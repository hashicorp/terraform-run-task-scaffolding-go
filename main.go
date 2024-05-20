// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"flag"

	"github.com/hashicorp/terraform-run-task-scaffolding-go/internal/runtask"
)

func main() {
	// Define command-line parameters
	var addr = flag.String("addr", "22180", "the port the run task HTTP server will run on")
	var path = flag.String("path", "/runtask", "the URL path for the run task to receive HTTP request from TFC or TFE")
	var hmacKey = flag.String("hmacKey", "", "the customizable secret which TFC or TFE will use to sign requests to the run task")
	flag.Parse()

	task := runtask.NewRunTask()
	task.Configure(*addr, *path, *hmacKey)
	runtask.HandleRequests(task)
}
