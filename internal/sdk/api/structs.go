// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"time"
)

// VerificationToken is a nonsense Terraform Cloud API token that should NEVER be valid.
const VerificationToken = "test-token"
const JsonApiMediaTypeHeader = "application/vnd.api+json"

type TaskStatus string

const (
	TaskFailed  TaskStatus = "failed"
	TaskPassed  TaskStatus = "passed"
	TaskRunning TaskStatus = "running"
)

const (
	PrePlan   string = "pre_plan"
	PostPlan  string = "post_plan"
	PreApply  string = "pre_apply"
	PostApply string = "post_apply"
)

type Request struct {
	AccessToken                     string    `json:"access_token"`
	ConfigurationVersionDownloadURL string    `json:"configuration_version_download_url,omitempty"`
	ConfigurationVersionID          string    `json:"configuration_version_id,omitempty"`
	IsSpeculative                   bool      `json:"is_speculative"`
	OrganizationName                string    `json:"organization_name"`
	PayloadVersion                  int       `json:"payload_version"`
	RunAppURL                       string    `json:"run_app_url"`
	RunCreatedAt                    time.Time `json:"run_created_at"`
	RunCreatedBy                    string    `json:"run_created_by"`
	RunID                           string    `json:"run_id"`
	RunMessage                      string    `json:"run_message"`
	Stage                           string    `json:"stage"`
	TaskResultCallbackURL           string    `json:"task_result_callback_url"`
	TaskResultEnforcementLevel      string    `json:"task_result_enforcement_level"`
	TaskResultID                    string    `json:"task_result_id"`
	VcsBranch                       string    `json:"vcs_branch,omitempty"`
	VcsCommitURL                    string    `json:"vcs_commit_url,omitempty"`
	VcsPullRequestURL               string    `json:"vcs_pull_request_url,omitempty"`
	VcsRepoURL                      string    `json:"vcs_repo_url,omitempty"`
	WorkspaceAppURL                 string    `json:"workspace_app_url"`
	WorkspaceID                     string    `json:"workspace_id"`
	WorkspaceName                   string    `json:"workspace_name"`
	WorkspaceWorkingDirectory       string    `json:"workspace_working_directory,omitempty"`
	PlanJSONAPIURL                  string    `json:"plan_json_api_url,omitempty"` // Specific to post-plan and pre-apply stage
}

// IsEndpointValidation returns true if the Request is from the
// run task service to validate this API endpoint. Callers should
// immediately return an HTTP 200 status code for these requests.
func (r Request) IsEndpointValidation() bool {
	return r.AccessToken == VerificationToken
}

type CallbackResponse struct {
	Data CallbackData `json:"data"`
}

type CallbackData struct {
	Type       string   `json:"type"`
	Attributes Response `json:"attributes"`
}

type Response struct {
	// A short message describing the status of the task.
	Message string `json:"message,omitempty"`
	// Must be one of TaskFailed, TaskPassed or TaskRunning
	Status TaskStatus `json:"status"`
	// URL that the user can use to get more information from the external service
	URL string `json:"url,omitempty"`
}
