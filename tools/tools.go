// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build generate

package tools

import (
	//NOTE: This is for HashiCorp specific licensing automation and can be deleted after creating a new repo with this template.
	_ "github.com/hashicorp/copywrite"
)

//NOTE: This is for HashiCorp specific licensing automation and can be deleted after creating a new repo with this template.
//go:generate go run github.com/hashicorp/copywrite headers -d .. --config ../.copywrite.hcl
