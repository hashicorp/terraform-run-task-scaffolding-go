// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package runtask

import (
	"log"
	"os"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-run-task-scaffolding-go/internal/sdk/api"

	"github.com/hashicorp/terraform-run-task-scaffolding-go/internal/sdk/handler"
)

const (
	DefaultBind = ":22180"
	DefaultPath = "/runtask"
	HMACKey     = "secret123"
)

// ScaffoldingRunTask defines the run task implementation.
type ScaffoldingRunTask struct {
	config handler.Configuration
	logger *log.Logger
}

// Configure defines the configuration for the server and run task.
// This method is called before the server is initialized.
func (r *ScaffoldingRunTask) Configure() {
	r.config = handler.Configuration{
		Addr:    DefaultBind,
		Path:    DefaultPath,
		HmacKey: HMACKey,
	}
}

// VerifyRequest defines custom run task integration logic.
// This method is called after the run task receives and validates the run task request from TFC.
func (r *ScaffoldingRunTask) VerifyRequest(request api.Request) (*handler.CallbackBuilder, error) {

	// Run custom verification logic
	//if request.OrganizationName != "TFC-ORG" {
	//	return handler.NewCallbackBuilder(api.TaskFailed).WithMessage("Unexpected Org Name")
	//}

	r.logger.Println("Successfully verified request")
	return handler.NewCallbackBuilder(api.TaskPassed).WithMessage("Custom Passed Message"), nil
}

// VerifyPlan defines custom integration logic for verifying the run's plan from TFC.
// This method is only called if the run task is running in the post-plan or pre-apply stages
// and if VerifyRequest returns a nil response with no error.
func (r *ScaffoldingRunTask) VerifyPlan(request api.Request, plan tfjson.Plan) (*handler.CallbackBuilder, error) {
	return nil, nil
}

// NewRunTask instantiates a new ScaffoldingRunTask with a new Logger.
func NewRunTask() *ScaffoldingRunTask {
	return &ScaffoldingRunTask{
		logger: log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
	}
}
