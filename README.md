# Terraform Run Task Scaffolding (Golang)

This repository is a *template* for a [Terraform Cloud](https://developer.hashicorp.com/terraform/cloud-docs/workspaces/settings/run-tasks) and/or [Terraform Enterprise](https://developer.hashicorp.com/terraform/enterprise/workspaces/settings/run-tasks) run task. It is intended as a starting point for creating Terraform run tasks, containing:

- A service handler for handling run task request/response (`internal/runtask/run_task_handler`),
- A scaffolding template file for configuring the service or business verification logic of the run task (`internal/runtask/run_task_scaffolding`),
- Miscellaneous meta files.

These files contain boilerplate code that you will need to edit to create your own Terraform run task. Detailed documentation for run task integration can be found on the  [HashiCorp Developer](https://developer.hashicorp.com/terraform/cloud-docs/integrations/run-tasks) platform.

Please see the [GitHub template repository documentation](https://help.github.com/en/github/creating-cloning-and-archiving-repositories/creating-a-repository-from-a-template) for how to create a new repository from this template on GitHub.

## Requirements

- A Terraform Cloud account or Terraform Enterprise >= v202206-1
  - To create a run task, you must have a user account with the [Manage Run Tasks](https://developer.hashicorp.com/terraform/cloud-docs/users-teams-organizations/permissions#manage-run-tasks) permission. To associate run tasks with a workspace, you need the [Manage Workspace Run Tasks](https://developer.hashicorp.com/terraform/cloud-docs/users-teams-organizations/permissions#general-workspace-permissions) permission on that particular workspace.
- [Go](https://golang.org/doc/install) >= 1.20

## Using The Run Task in TFC/E
_We highly recommend using a service like [ngrok](https://ngrok.com/) to quickly test your run task without a full cloud deployment_

1. Run the run task using the Go `run` command:

```shell
go run main.go
```
2. Take note of the `path`, and `hmac` values configured in `internal/runtask/run_task_scaffolding.Configure()`
   - Default values: `path` = `/runtask`, and `hmac` = `secret123`
3. Follow the steps on the Hashicorp Developer platform for [Creating a Run Task](https://developer.hashicorp.com/terraform/cloud-docs/workspaces/settings/run-tasks#creating-a-run-task) and [Associating Run Tasks with a Workspace](https://developer.hashicorp.com/terraform/cloud-docs/workspaces/settings/run-tasks#associating-run-tasks-with-a-workspace)
   - For the **Endpoint URL** field, append the `path` value to the end of your service's URL, ex: `http://myservice.io/runtask`
   - For the **HMAC key** field, use the configured `hmac` value 



## Adding Dependencies

This run task uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up-to-date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform run task:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.