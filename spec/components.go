// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package spec

// Components Holds a set of reusable objects for different aspects of the OAS.
// All objects defined within the components object will have no effect on the API unless they are explicitly referenced from properties outside the components object.
//
// This object MAY be extended with Specification Extensions.
//
// All the fixed fields declared below are objects that MUST use keys that match the regular expression: ^[a-zA-Z0-9\.\-_]+$.
//
//	"components": {
//	 "schemas": {
//	   "GeneralError": {
//	     "type": "object",
//	     "properties": {
//	       "code": {
//	         "type": "integer",
//	         "format": "int32"
//	       },
//	       "message": {
//	         "type": "string"
//	       }
//	     }
//	   },
//	   "Category": {
//	     "type": "object",
//	     "properties": {
//	       "id": {
//	         "type": "integer",
//	         "format": "int64"
//	       },
//	       "name": {
//	         "type": "string"
//	       }
//	     }
//	   },
//	   "Tag": {
//	     "type": "object",
//	     "properties": {
//	       "id": {
//	         "type": "integer",
//	         "format": "int64"
//	       },
//	       "name": {
//	         "type": "string"
//	       }
//	     }
//	   }
//	 },
//	 "parameters": {
//	   "skipParam": {
//	     "name": "skip",
//	     "in": "query",
//	     "description": "number of items to skip",
//	     "required": true,
//	     "schema": {
//	       "type": "integer",
//	       "format": "int32"
//	     }
//	   },
//	   "limitParam": {
//	     "name": "limit",
//	     "in": "query",
//	     "description": "max records to return",
//	     "required": true,
//	     "schema" : {
//	       "type": "integer",
//	       "format": "int32"
//	     }
//	   }
//	 },
//	 "responses": {
//	   "NotFound": {
//	     "description": "Entity not found."
//	   },
//	   "IllegalInput": {
//	     "description": "Illegal input for operation."
//	   },
//	   "GeneralError": {
//	     "description": "General Error",
//	     "content": {
//	       "application/json": {
//	         "schema": {
//	           "$ref": "#/components/schemas/GeneralError"
//	         }
//	       }
//	     }
//	   }
//	 },
//	 "securitySchemes": {
//	   "api_key": {
//	     "type": "apiKey",
//	     "name": "api_key",
//	     "in": "header"
//	   },
//	   "petstore_auth": {
//	     "type": "oauth2",
//	     "flows": {
//	       "implicit": {
//	         "authorizationUrl": "http://example.org/api/oauth/dialog",
//	         "scopes": {
//	           "write:pets": "modify pets in your account",
//	           "read:pets": "read your pets"
//	         }
//	       }
//	     }
//	   }
//	 }
//	}
type Components struct {
	// An object to hold reusable Schema Objects.
	Schemas Schemas `json:"schemas,omitempty"`

	// An object to hold reusable Response Objects.
	Responses Responses `json:"responses,omitempty"`

	// An object to hold reusable Parameter Objects.
	Parameters Parameters `json:"parameters,omitempty"`

	// An object to hold reusable Example Objects.
	Examples Examples `json:"examples,omitempty"`

	// An object to hold reusable Request Body Objects.
	RequestBodies RequestBodies `json:"requestBodies,omitempty"`

	// An object to hold reusable Header Objects.
	Headers Headers `json:"headers,omitempty"`

	// An object to hold reusable Security Scheme Objects.
	SecuritySchemes SecuritySchemes `json:"securitySchemes,omitempty"`

	// An object to hold reusable Link Objects.
	Links Links `json:"links,omitempty"`

	// An object to hold reusable Callback Objects.
	Callbacks Callbacks `json:"callbacks,omitempty"`
}
