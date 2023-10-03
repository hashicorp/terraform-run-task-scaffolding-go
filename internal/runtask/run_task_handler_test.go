// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package runtask

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"

	"github.com/hashicorp/terraform-run-task-scaffolding-go/internal/sdk/api"
)

func Test_runTask_handler(t *testing.T) {
	testToken := fmt.Sprintf("\"access_token\": \"%s\"", api.VerificationToken)
	tests := []struct {
		name         string
		method       string
		endPoint     string
		bodyJSON     string
		wantHTTPCode int
	}{
		{
			name:         "with validation token",
			method:       http.MethodPost,
			endPoint:     "/test",
			bodyJSON:     "{" + testToken + "}",
			wantHTTPCode: http.StatusOK,
		},
		{
			name:         "no body",
			method:       http.MethodPost,
			endPoint:     "/test",
			bodyJSON:     "{" + "}",
			wantHTTPCode: http.StatusNotFound,
		},
		{
			name:         "with malformed body",
			method:       http.MethodPost,
			endPoint:     "/test",
			bodyJSON:     "{",
			wantHTTPCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := NewRunTask()
			task.config.HmacKey = ""
			router := CreateMockService(task)

			req, _ := http.NewRequest(tt.method, tt.endPoint, strings.NewReader(tt.bodyJSON))
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if rr.Code != tt.wantHTTPCode {
				t.Errorf("Expected response code %d, got %d", tt.wantHTTPCode, rr.Code)
			}
		})
	}
}

func CreateMockService(task *ScaffoldingRunTask) *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/test", handleTFCRequestWrapper(task, sendTFCCallbackResponse())).Methods(http.MethodPost)

	return router
}
