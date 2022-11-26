// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package spec

// Ref A simple object to allow referencing other components in the specification, internally and externally.
//
// The Reference Object is defined by JSON Reference and follows the same structure, behavior and rules.
//
// For this specification, reference resolution is accomplished as defined by the JSON Reference specification and not by the JSON Schema specification.
//
// Reference object example
//
//	{
//		"$ref": "#/components/schemas/Pet"
//	}
//
// Relative schema document example
//
//	{
//	 "$ref": "Pet.json"
//	}
//
// Relative documents with embedded schema example
//
//	{
//	 "$ref": "definitions.json#/Pet"
//	}
//
// This may not be really used.
type Ref struct {
	// REQUIRED. The reference string.
	Ref string `json:"$ref"`
}
