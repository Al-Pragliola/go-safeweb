// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package defaults provides ready to use, safe, pre-configured instances of safehttp types.
package defaults

import (
	"errors"

	"github.com/google/go-safeweb/safehttp"
	"github.com/google/go-safeweb/safehttp/plugins/coop"
	"github.com/google/go-safeweb/safehttp/plugins/csp"
	"github.com/google/go-safeweb/safehttp/plugins/fetchmetadata"
	"github.com/google/go-safeweb/safehttp/plugins/framing"
	"github.com/google/go-safeweb/safehttp/plugins/hostcheck"
	"github.com/google/go-safeweb/safehttp/plugins/hsts"
	"github.com/google/go-safeweb/safehttp/plugins/staticheaders"
	"github.com/google/go-safeweb/safehttp/plugins/xsrf/xsrfhtml"
)

// ServeMuxConfig creates a safe and ready to use ServeMuxConfig with all necessary security interceptors installed.
// hosts should be all the hosts this mux will be served on and it can't be empty.
// xsrfKey is the secret application key to use for XSRF token generation and it can't be empty.
func ServeMuxConfig(hosts []string, xsrfKey string) (*safehttp.ServeMuxConfig, error) {

	if len(hosts) == 0 {
		return nil, errors.New("hosts slice cannot be empty")
	}

	if xsrfKey == "" {
		return nil, errors.New("xsrfKey cannot be empty")
	}

	c := safehttp.NewServeMuxConfig(nil)

	// Non-blocking:
	c.Intercept(staticheaders.Interceptor{})
	c.Intercept(hsts.Default())
	// TODO(empijei): add a report group once we support reporting.
	c.Intercept(coop.Default(""))
	// TODO(empijei): add a report-uri once we support reporting.
	for _, i := range csp.Default("") {
		c.Intercept(i)
	}

	// Blocking:
	c.Intercept(hostcheck.New(hosts...))
	c.Intercept(fetchmetadata.ResourceIsolationPolicy())
	c.Intercept(&xsrfhtml.Interceptor{SecretAppKey: xsrfKey})
	for _, i := range framing.Interceptors("") {
		c.Intercept(i)
	}

	return c, nil
}
