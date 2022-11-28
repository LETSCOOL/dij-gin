// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package spec

import "strings"

type Openapi struct {
	// REQUIRED. This string MUST be the semantic version number of the OpenAPI Specification version that the OpenAPI document uses.
	// The openapi field SHOULD be used by tooling specifications and clients to interpret the OpenAPI document. This is not related to the API info.version string.
	Openapi string `json:"openapi"` // ex: openapi: "3.0.0"

	// REQUIRED. Provides metadata about the API. The metadata MAY be used by tooling as required.
	Info *Info `json:"info"`

	// An array of Server Objects, which provide connectivity information to a target server.
	// If the servers property is not provided, or is an empty array, the default value would be a Server Object with a url value of /.
	Servers []Server `json:"servers"`

	// REQUIRED. The available paths and operations for the API.
	Paths Paths `json:"paths,omitempty"` // ex: "/user":Path

	// An element to hold various schemas for the specification.
	Components *Components `json:"components,omitempty"`

	// A declaration of which security mechanisms can be used across the API.
	// The list of values includes alternative security requirement objects that can be used.
	// Only one of the security requirement objects need to be satisfied to authorize a request.
	// Individual operations can override this definition. To make security optional, an empty security requirement ({}) can be included in the array.
	Security []SecurityRequirement `json:"securityDefinitions,omitempty"` // ex: {"api_key": SecurityDefinition}

	// A list of tags used by the specification with additional metadata. The order of the tags can be used to reflect on their order by the parsing tools.
	// Not all tags that are used by the Operation Object must be declared. The tags that are not declared MAY be organized randomly or based on the tools' logic.
	// Each tag name in the list MUST be unique.
	Tags Tags `json:"tags,omitempty"`

	// Additional external documentation.
	ExternalDocs *ExternalDoc `json:"externalDocs,omitempty"`
}

func (s *Openapi) AddPathOperation(path string, method string, def Operation) {
	if s.Paths == nil {
		s.Paths = Paths{}
	}
	p, b := s.Paths[path]
	if !b {
		p = Path{}
	}
	switch strings.ToLower(method) {
	case "get":
		p.Get = &def
	case "put":
		p.Put = &def
	case "post":
		p.Post = &def
	case "delete":
		p.Delete = &def
	case "options":
		p.Options = &def
	case "head":
		p.Head = &def
	case "patch":
		p.Patch = &def
	case "trace":
		p.Trace = &def
	}
	s.Paths[path] = p
	// p[method] = def
}
