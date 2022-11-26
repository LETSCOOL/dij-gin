// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package spec

// Info The object provides metadata about the API. The metadata MAY be used by the clients if needed,
// and MAY be presented in editing or documentation generation tools for convenience.
//
//	{
//	 "title": "Sample Pet Store App",
//	 "description": "This is a sample server for a pet store.",
//	 "termsOfService": "http://example.com/terms/",
//	 "contact": {
//	   "name": "API Support",
//	   "url": "http://www.example.com/support",
//	   "email": "support@example.com"
//	 },
//	 "license": {
//	   "name": "Apache 2.0",
//	   "url": "https://www.apache.org/licenses/LICENSE-2.0.html"
//	 },
//	 "version": "1.0.1"
//	}
type Info struct {
	// REQUIRED. The title of the API.
	Title string `json:"title"`

	// A short description of the API. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// A URL to the Terms of Service for the API. MUST be in the format of a URL.
	TermsOfService string `json:"termsOfService,omitempty"`

	// The contact information for the exposed API.
	Contact *Contact `json:"contact,omitempty"`

	// The license information for the exposed API.
	License *License `json:"license,omitempty"`

	// REQUIRED. The version of the OpenAPI document (which is distinct from the OpenAPI Specification version or the API implementation version).
	Version string `json:"version,omitempty"`
}
