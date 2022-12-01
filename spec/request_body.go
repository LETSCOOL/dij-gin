// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package spec

// RequestBody describes a single request body.
//
// A request body with a referenced model definition.
//
//	{
//	 "description": "user to add to the system",
//	 "content": {
//	   "application/json": {
//	     "schema": {
//	       "$ref": "#/components/schemas/User"
//	     },
//	     "examples": {
//	         "user" : {
//	           "summary": "User Example",
//	           "externalValue": "http://foo.bar/examples/user-example.json"
//	         }
//	       }
//	   },
//	   "application/xml": {
//	     "schema": {
//	       "$ref": "#/components/schemas/User"
//	     },
//	     "examples": {
//	         "user" : {
//	           "summary": "User example in XML",
//	           "externalValue": "http://foo.bar/examples/user-example.xml"
//	         }
//	       }
//	   },
//	   "text/plain": {
//	     "examples": {
//	       "user" : {
//	           "summary": "User example in Plain text",
//	           "externalValue": "http://foo.bar/examples/user-example.txt"
//	       }
//	     }
//	   },
//	   "*/*": {
//	     "examples": {
//	       "user" : {
//	           "summary": "User example in other format",
//	           "externalValue": "http://foo.bar/examples/user-example.whatever"
//	       }
//	     }
//	   }
//	 }
//	}
//
// A body parameter that is an array of string values:
//
//	{
//	 "description": "user to add to the system",
//	 "content": {
//	   "text/plain": {
//	     "schema": {
//	       "type": "array",
//	       "items": {
//	         "type": "string"
//	       }
//	     }
//	   }
//	 }
//	}
//
// In contrast with the 2.0 specification, file input/output content in OpenAPI is described with the same semantics as any other schema type.
// Refer to https://swagger.io/specification/.
type RequestBody struct {
	// A brief description of the request body. This could contain examples of use. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// REQUIRED. The content of the request body. The key is a media type or media type range and the value describes it.
	// For requests that match multiple keys, only the most specific key is applicable. e.g. text/plain overrides text/*
	Content Content `json:"content"`

	// Determines if the request body is required in the request. Defaults to false.
	Required bool `json:"required,omitempty"`
}

type RequestBodies map[string]RequestBodyR

// RequestBodyR presents RequestBody or Ref combination
type RequestBodyR struct {
	*RequestBody `json:""`
	Ref          string `json:"$ref,omitempty"`
}

func (r *RequestBodyR) SetMediaType(coding MediaTypeTitle, mediaType MediaType) {
	if r.RequestBody == nil {
		r.RequestBody = &RequestBody{Content: Content{}}
	}
	if r.RequestBody.Content == nil {
		r.RequestBody.Content = Content{}
	}
	r.RequestBody.Content.SetMediaType(coding, mediaType)
}
