// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"terraform-run-task-scaffolding-go/internal/runtask"
)

func main() {
	task := runtask.NewRunTask()
	task.Configure()
	runtask.HandleRequests(task)
}
