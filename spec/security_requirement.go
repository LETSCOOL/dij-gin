// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package spec

// SecurityRequirement Lists the required security schemes to execute this operation.
// The name used for each property MUST correspond to a security scheme declared in the Security Schemes under the Components Object.
//
// Security Requirement Objects that contain multiple schemes require that all schemes MUST be satisfied for a request to be authorized.
// This enables support for scenarios where multiple query parameters or HTTP headers are required to convey security information.
//
// When a list of Security Requirement Objects is defined on the OpenAPI Object or Operation Object,
// only one of the Security Requirement Objects in the list needs to be satisfied to authorize the request.
//
// Each name MUST correspond to a security scheme which is declared in the Security Schemes under the Components Object.
// If the security scheme is of type "oauth2" or "openIdConnect", then the value is a list of scope names required for the execution,
// and the list MAY be empty if authorization does not require a specified scope. For other security scheme types, the array MUST be empty.
//
// Non-OAuth2 Security Requirement
//
//	{
//	 "api_key": []
//	}
//
// OAuth2 Security Requirement
//
//	{
//	 "petstore_auth": [
//	   "write:pets",
//	   "read:pets"
//	 ]
//	}
//
// Optional OAuth2 security as would be defined in an OpenAPI Object or an Operation Object:
//
//	{
//	 "security": [
//	   {},
//	   {
//	     "petstore_auth": [
//	       "write:pets",
//	       "read:pets"
//	     ]
//	   }
//	 ]
//	}
type SecurityRequirement map[string][]string
