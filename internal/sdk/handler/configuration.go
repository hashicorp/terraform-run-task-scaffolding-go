// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package handler

type Configuration struct {
	// Addr specifies the TCP address for the server to listen on.
	Addr string
	// Path defines a matcher for the route URL path. It accepts a template with zero or more URL variables enclosed by {}.
	// The template must start with a "/".
	Path string
	// HmacKey defines the HMAC Key used for verifying the TFC request.
	HmacKey string
}
