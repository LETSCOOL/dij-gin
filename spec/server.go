// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package spec

// Server is an object representing a Server.
//
//	{
//	 "url": "https://development.gigantic-server.com/v1",
//	 "description": "Development server"
//	}
//
// The following shows how multiple servers can be described, for example, at the OpenAPI Object's servers:
//
//	{
//	 "servers": [
//	   {
//	     "url": "https://development.gigantic-server.com/v1",
//	     "description": "Development server"
//	   },
//	   {
//	     "url": "https://staging.gigantic-server.com/v1",
//	     "description": "Staging server"
//	   },
//	   {
//	     "url": "https://api.gigantic-server.com/v1",
//	     "description": "Production server"
//	   }
//	 ]
//	}
type Server struct {
	// REQUIRED. A URL to the target host. This URL supports Server Variables and MAY be relative, to indicate
	// that the host location is relative to the location where the OpenAPI document is being served.
	// Variable substitutions will be made when a variable is named in {brackets}.
	Url string `json:"url"` // ex. https://development.gigantic-server.com/v1

	// An optional string describing the host designated by the URL. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// A map between a variable name and its value. The value is used for substitution in the server's URL template.
	Variables map[string]ServerVariable `json:"variables,omitempty"`
}

// ServerVariable is an object representing a Server Variable for server URL template substitution.
//
//	{
//	  "url": "https://{username}.gigantic-server.com:{port}/{basePath}",
//	  "description": "The production API server",
//	  "variables": {
//	    "username": {
//	      "default": "demo",
//	      "description": "this value is assigned by the service provider, in this example `gigantic-server.com`"
//	    },
//	    "port": {
//	      "enum": [
//	        "8443",
//	        "443"
//	      ],
//	      "default": "8443"
//	    },
//	    "basePath": {
//	      "default": "v2"
//	    }
//	  }
//	}
type ServerVariable struct {
	// An enumeration of string values to be used if the substitution options are from a limited set. The array SHOULD NOT be empty.
	Enum []string `json:"enum,omitempty"`

	// REQUIRED. The default value to use for substitution, which SHALL be sent if an alternate value is not supplied.
	// Note this behavior is different from the Schema Object's treatment of default values, because in those cases parameter values are optional.
	// If the enum is defined, the value SHOULD exist in the enum's values.
	Default string `json:"default"`

	// An optional description for the server variable. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`
}
