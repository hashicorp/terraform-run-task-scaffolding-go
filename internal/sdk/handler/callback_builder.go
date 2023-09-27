// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"encoding/json"

	"terraform-run-task-scaffolding-go/internal/sdk/api"
)

// TypeTaskResults is the data type used in run task responses.
const TypeTaskResults = "task-results"

type callbackResponse struct {
	Data callbackData `json:"data"`
}

type callbackData struct {
	Type       string       `json:"type"`
	Attributes api.Response `json:"attributes"`
}

type CallbackBuilder struct {
	resp callbackResponse
}

func NewCallbackBuilder(status api.TaskStatus) *CallbackBuilder {
	return &CallbackBuilder{
		resp: callbackResponse{
			Data: callbackData{
				Type: TypeTaskResults,
				Attributes: api.Response{
					Status: status,
				},
			},
		},
	}
}

func (cb *CallbackBuilder) WithMessage(message string) *CallbackBuilder {
	cb.resp.Data.Attributes.Message = message
	return cb
}

func (cb *CallbackBuilder) MarshallJSON() ([]byte, error) {
	return json.Marshal(cb.resp)
}
