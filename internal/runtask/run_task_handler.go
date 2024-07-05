// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package runtask

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-run-task-scaffolding-go/internal/sdk/api"
	"github.com/hashicorp/terraform-run-task-scaffolding-go/internal/sdk/handler"
)

func HandleRequests(task *ScaffoldingRunTask) {
	r := mux.NewRouter()

	task.logger.Println("Registering " + task.config.Path + " route")
	r.HandleFunc(task.config.Path, handleTFCRequestWrapper(task, sendTFCCallbackResponse())).Methods(http.MethodPost)

	task.logger.Println("Registering /healthcheck route")
	r.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		task.logger.Println("/healthcheck called")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(map[string]string{"status": "available"})
		if err != nil {
			return
		}
	}).Methods(http.MethodGet)

	task.logger.Printf("Starting server on port %s", task.config.Addr)
	err := http.ListenAndServe(task.config.Addr, r)
	if err != nil {
		return
	}
}

func handleTFCRequestWrapper(task *ScaffoldingRunTask, original func(http.ResponseWriter, *http.Request, api.Request, *ScaffoldingRunTask, *handler.CallbackBuilder)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		task.logger.Println(task.config.Path + " called")

		// Parse request
		var runTaskReq api.Request
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			task.logger.Println("Error occurred while parsing the request")
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(reqBody, &runTaskReq)
		if err != nil {
			task.logger.Println("Error occurred while parsing the request")
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		requestSha := r.Header.Get(handler.HeaderTaskSignature)

		if requestSha != "" && task.config.HmacKey == "" {
			task.logger.Printf("Received a request for %s with a signature but this server cannot validate signed requests\n", r.URL)
			http.Error(w, "Unexpected x-tfc-task-signature header", http.StatusBadRequest)
			return
		}

		if requestSha == "" && task.config.HmacKey != "" {
			task.logger.Printf("Received an unsigned request for %s\n", r.URL)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if requestSha != "" {
			// Calculate expected HMAC
			verified, err := handler.VerifyHMAC(reqBody, []byte(r.Header.Get(handler.HeaderTaskSignature)), []byte(task.config.HmacKey))

			if err != nil {
				task.logger.Println("Unable to verify given HMAC key")
				http.Error(w, "Error verifying signed request", http.StatusInternalServerError)
				return
			}

			if !verified {
				task.logger.Println("Received an unauthorized request")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

		if runTaskReq.IsEndpointValidation() {
			task.logger.Println("Successfully validated TFC request")
			w.WriteHeader(http.StatusOK)
			return
		}

		callbackResp, err := task.VerifyRequest(runTaskReq)
		if err != nil {
			task.logger.Println("Error occurred during run task request verification")
			http.Error(w, "Error during run task request verification:"+err.Error(), http.StatusInternalServerError)
			return
		}

		// Get TFC Plan if the task is running in the post-plan or pre-apply stages
		if runTaskReq.Stage == api.PostPlan || runTaskReq.Stage == api.PreApply {
			plan, err := retrieveTFCPlan(runTaskReq)

			if err != nil {
				task.logger.Println("Error occurred while retrieving plan from TFC")
				http.Error(w, "Bad Request: "+err.Error(), http.StatusNotFound)
				return
			}
			task.logger.Println("Successfully retrieved plan from TFC")

			callbackResp, err = task.VerifyPlan(runTaskReq, plan)
			if err != nil {
				task.logger.Println("Error occurred while verifying plan")
				http.Error(w, "Error verifying plan:"+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		original(w, r, runTaskReq, task, callbackResp)
	}
}

func sendTFCCallbackResponse() func(w http.ResponseWriter, r *http.Request, reqBody api.Request, task *ScaffoldingRunTask, cbBuilder *handler.CallbackBuilder) {

	return func(w http.ResponseWriter, r *http.Request, reqBody api.Request, task *ScaffoldingRunTask, cbBuilder *handler.CallbackBuilder) {

		respBody, err := cbBuilder.MarshallJSON()
		if err != nil {
			task.logger.Println("Unable to marshall callback response to TFC")
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Send PATCH callback response to TFC
		request, err := sendTFCRequest(reqBody.TaskResultCallbackURL, http.MethodPatch, reqBody.AccessToken, respBody)
		if request != nil {
			_ = r.Body.Close()
		}
		if err != nil {
			task.logger.Println("Error occurred while sending the callback response to TFC")
			http.Error(w, "Bad Request:"+err.Error(), http.StatusNotFound)
			return
		}

		task.logger.Println("Sent run task response to TFC")
	}

}

func retrieveTFCPlan(req api.Request) (tfjson.Plan, error) {

	// Call TFC to get plan
	resp, err := sendTFCRequest(req.PlanJSONAPIURL, "GET", req.AccessToken, nil)
	if err != nil {
		return tfjson.Plan{}, err
	}

	var tfPlan tfjson.Plan

	if resp == nil {
		return tfPlan, fmt.Errorf("expected Terraform plan from TFC but received none")
	}

	respBody, err := io.ReadAll(resp.Body)

	_ = resp.Body.Close()

	if err != nil {
		return tfPlan, err
	}

	err = json.Unmarshal(respBody, &tfPlan)

	if err != nil {
		return tfPlan, err
	}

	return tfPlan, nil
}

func sendTFCRequest(url string, method string, accessToken string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	// Required headers to send to TFC
	req.Header.Set("Content-Type", api.JsonApiMediaTypeHeader)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	return http.DefaultClient.Do(req)
}
