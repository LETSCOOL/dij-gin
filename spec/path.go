// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package spec

// Paths Holds the relative paths to the individual endpoints and their operations.
// The path is appended to the URL from the Server Object in order to construct the full URL. The Paths MAY be empty, due to ACL constraints.
//
// A relative path to an individual endpoint. The field name MUST begin with a forward slash (/).
// The path is appended (no relative URL resolution) to the expanded URL from the Server Object's url field in order to construct the full URL.
// Path templating is allowed. When matching URLs, concrete (non-templated) paths would be matched before their templated counterparts.
// Templated paths with the same hierarchy but different templated names MUST NOT exist as they are identical.
// In case of ambiguous matching, it's up to the tooling to decide which one to use.
//
//	{
//	 "/pets": {
//	   "get": {
//	     "description": "Returns all pets from the system that the user has access to",
//	     "responses": {
//	       "200": {
//	         "description": "A list of pets.",
//	         "content": {
//	           "application/json": {
//	             "schema": {
//	               "type": "array",
//	               "items": {
//	                 "$ref": "#/components/schemas/pet"
//	               }
//	             }
//	           }
//	         }
//	       }
//	     }
//	   }
//	 }
//	}
type Paths = map[string]Path

// Path describes the operations available on a single path. A Path Item MAY be empty, due to ACL constraints.
// The path itself is still exposed to the documentation viewer but they will not know which operations and parameters are available.
//
//	{
//	 "get": {
//	   "description": "Returns pets based on ID",
//	   "summary": "Find pets by ID",
//	   "operationId": "getPetsById",
//	   "responses": {
//	     "200": {
//	       "description": "pet response",
//	       "content": {
//	         "*/*": {
//	           "schema": {
//	             "type": "array",
//	             "items": {
//	               "$ref": "#/components/schemas/Pet"
//	             }
//	           }
//	         }
//	       }
//	     },
//	     "default": {
//	       "description": "error payload",
//	       "content": {
//	         "text/html": {
//	           "schema": {
//	             "$ref": "#/components/schemas/ErrorModel"
//	           }
//	         }
//	       }
//	     }
//	   }
//	 },
//	 "parameters": [
//	   {
//	     "name": "id",
//	     "in": "path",
//	     "description": "ID of pet to use",
//	     "required": true,
//	     "schema": {
//	       "type": "array",
//	       "items": {
//	         "type": "string"
//	       }
//	     },
//	     "style": "simple"
//	   }
//	 ]
//	}
type Path struct {
	// Allows for an external definition of this path item. The referenced structure MUST be in the format of a Path Item Object.
	// In case a Path Item Object field appears both in the defined object and the referenced object, the behavior is undefined.
	Ref string `json:"$ref,omitempty"`

	// An optional, string summary, intended to apply to all operations in this path.
	Summary string `json:"summary,omitempty"`

	// An optional, string description, intended to apply to all operations in this path. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// A definition of a GET operation on this path.
	Get *Operation `json:"get,omitempty"`

	// A definition of a PUT operation on this path.
	Put *Operation `json:"put,omitempty"`

	// A definition of a POST operation on this path.
	Post *Operation `json:"post,omitempty"`

	// A definition of a DELETE operation on this path.
	Delete *Operation `json:"delete,omitempty"`

	// A definition of a OPTIONS operation on this path.
	Options *Operation `json:"options,omitempty"`

	// A definition of a HEAD operation on this path.
	Head *Operation `json:"head,omitempty"`

	// A definition of a PATCH operation on this path.
	Patch *Operation `json:"patch,omitempty"`

	// A definition of a TRACE operation on this path.
	Trace *Operation `json:"trace,omitempty"`

	// An alternative server array to service all operations in this path.
	Servers []Server `json:"servers,omitempty"`

	// A list of parameters that are applicable for all the operations described under this path.
	// These parameters can be overridden at the operation level, but cannot be removed there.
	// The list MUST NOT include duplicated parameters. A unique parameter is defined by a combination of a name and location.
	// The list can use the Reference Object to link to parameters that are defined at the OpenAPI Object's components/parameters.
	Parameters ParameterList `json:"parameters,omitempty"`
}
