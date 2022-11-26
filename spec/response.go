// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package spec

// Responses A container for the expected responses of an operation. The container maps a HTTP response code to the expected response.
//
// The documentation is not necessarily expected to cover all possible HTTP response codes because they may not be known in advance. However, documentation is expected to cover a successful operation response and any known errors.
//
// The default MAY be used as a default response object for all HTTP codes that are not covered individually by the specification.
//
// The Responses Object MUST contain at least one response code, and it SHOULD be the response for a successful operation call.
//
//	{
//	 "200": {
//	   "description": "a pet to be returned",
//	   "content": {
//	     "application/json": {
//	       "schema": {
//	         "$ref": "#/components/schemas/Pet"
//	       }
//	     }
//	   }
//	 },
//	 "default": {
//	   "description": "Unexpected error",
//	   "content": {
//	     "application/json": {
//	       "schema": {
//	         "$ref": "#/components/schemas/ErrorModel"
//	       }
//	     }
//	   }
//	 }
//	}
type Responses map[string]ResponseR

// Response describes a single response from an API Operation, including design-time, static links to operations based on the response.
//
// Response of an array of a complex type:
//
//	{
//	 "description": "A complex object array response",
//	 "content": {
//	   "application/json": {
//	     "schema": {
//	       "type": "array",
//	       "items": {
//	         "$ref": "#/components/schemas/VeryComplexType"
//	       }
//	     }
//	   }
//	 }
//	}
//
// Response with a string type:
//
//	{
//	 "description": "A simple string response",
//	 "content": {
//	   "text/plain": {
//	     "schema": {
//	       "type": "string"
//	     }
//	   }
//	 }
//	}
//
// Plain text response with headers:
//
//	{
//	 "description": "A simple string response",
//	 "content": {
//	   "text/plain": {
//	     "schema": {
//	       "type": "string",
//	       "example": "whoa!"
//	     }
//	   }
//	 },
//	 "headers": {
//	   "X-Rate-Limit-Limit": {
//	     "description": "The number of allowed requests in the current period",
//	     "schema": {
//	       "type": "integer"
//	     }
//	   },
//	   "X-Rate-Limit-Remaining": {
//	     "description": "The number of remaining requests in the current period",
//	     "schema": {
//	       "type": "integer"
//	     }
//	   },
//	   "X-Rate-Limit-Reset": {
//	     "description": "The number of seconds left in the current period",
//	     "schema": {
//	       "type": "integer"
//	     }
//	   }
//	 }
//	}
//
// Response with no return value:
//
//	{
//	 "description": "object created"
//	}
type Response struct {
	// REQUIRED. A short description of the response. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description"`

	// Maps a header name to its definition. RFC7230 states header names are case-insensitive.
	// If a response header is defined with the name "Content-Type", it SHALL be ignored.
	Headers map[string]HeaderR `json:"headers,omitempty"`

	// A map containing descriptions of potential response payloads. The key is a media type or media type range and the value describes it.
	// For responses that match multiple keys, only the most specific key is applicable. e.g. text/plain overrides text/*
	Content map[string]MediaType `json:"content,omitempty"`

	// A map of operations links that can be followed from the response. The key of the map is a short name for the link,
	// following the naming constraints of the names for Component Objects.
	Links map[string]LinkR `json:"links,omitempty"`
}

// ResponseR presents Response or Ref combination
type ResponseR struct {
	*Response `json:""`
	Ref       string `json:"$ref,omitempty"`
}
