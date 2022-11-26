// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package spec

// Tag adds metadata to a single tag that is used by the Operation Object.
// It is not mandatory to have a Tag Object per tag defined in the Operation Object instances.
//
//	{
//		"name": "pet",
//		"description": "Pets operations"
//	}
type Tag struct {
	// REQUIRED. The name of the tag.
	Name string `json:"name"`

	// A short description for the tag. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description"`

	// Additional external documentation for this tag.
	ExternalDocs *ExternalDoc `json:"externalDocs,omitempty"`
}

type Tags []Tag
