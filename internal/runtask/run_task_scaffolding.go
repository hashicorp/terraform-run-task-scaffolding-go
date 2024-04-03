// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package runtask

import (
	"fmt"
	"log"
	"os"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-run-task-scaffolding-go/internal/sdk/api"

	"github.com/hashicorp/terraform-run-task-scaffolding-go/internal/sdk/handler"
)

// ScaffoldingRunTask defines the run task implementation.
type ScaffoldingRunTask struct {
	config handler.Configuration
	logger *log.Logger
}

// Configure defines the configuration for the server and run task.
// This method is called before the server is initialized.
func (r *ScaffoldingRunTask) Configure(addr string, path string, hmacKey string) {
	r.config = handler.Configuration{
		Addr:    fmt.Sprintf(":%s", addr),
		Path:    path,
		HmacKey: hmacKey,
	}
}

// VerifyRequest defines custom run task integration logic.
// This method is called after the run task receives and validates the run task request from TFC.
func (r *ScaffoldingRunTask) VerifyRequest(request api.Request) (*handler.CallbackBuilder, error) {

	// Run custom verification logic
	r.logger.Println("Successfully verified request")
	return handler.NewCallbackBuilder(api.TaskPassed).WithMessage("Custom Passed Message"), nil
}

// VerifyPlan defines custom integration logic for verifying the run's plan from TFC.
// This method is only called if the run task is running in the post-plan or pre-apply stages
// and if VerifyRequest returns a nil response with no error.
func (r *ScaffoldingRunTask) VerifyPlan(request api.Request, plan tfjson.Plan) (*handler.CallbackBuilder, error) {
	r.logger.Println("Successfully verified plan")
	return handler.NewCallbackBuilder(api.TaskPassed).WithMessage("Custom Passed Message"), nil
}

// NewRunTask instantiates a new ScaffoldingRunTask with a new Logger.
func NewRunTask() *ScaffoldingRunTask {
	return &ScaffoldingRunTask{
		logger: log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
	}
}
