// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package spec

// The Header Object follows the structure of the Parameter Object with the following changes:
//
//  1. name MUST NOT be specified, it is given in the corresponding headers map.
//  2. in MUST NOT be specified, it is implicitly in header.
//  3. All traits that are affected by the location MUST be applicable to a location of header (for example, style).
//
// A simple header of type integer:
//
//	{
//	 "description": "The number of allowed requests in the current period",
//	 "schema": {
//	   "type": "integer"
//	 }
//	}
type Header struct {
	// A brief description of the parameter. This could contain examples of use. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// Determines whether this parameter is mandatory. If the parameter location is "path", this property is REQUIRED and its value MUST be true.
	// Otherwise, the property MAY be included and its default value is false.
	Required bool `json:"required,omitempty"`

	// Specifies that a parameter is deprecated and SHOULD be transitioned out of usage. Default value is false.
	Deprecated bool `json:"deprecated,omitempty"`

	// Sets the ability to pass empty-valued parameters. This is valid only for query parameters and allows sending a parameter with an empty value.
	// Default value is false. If style is used, and if behavior is n/a (cannot be serialized), the value of allowEmptyValue SHALL be ignored. Use of this property is NOT RECOMMENDED, as it is likely to be removed in a later revision.
	AllowEmptyValue bool `json:"allowEmptyValue,omitempty"`

	// WithSchemaAndStyle The rules for serialization of the parameter are specified in one of two ways.
	// For simpler scenarios, a schema and style can describe the structure and syntax of the parameter.
	WithSchemaAndStyle

	// For more complex scenarios, the content property can define the media type and schema of the parameter.
	// A parameter MUST contain either a schema property, or a content property, but not both.
	// When example or examples are provided in conjunction with the schema object, the example MUST follow the prescribed serialization strategy for the parameter.
	WithContent
}

type Headers map[string]HeaderR

// HeaderR presents Header or Ref combination
type HeaderR struct {
	*Header `json:""`
	Ref     string `json:"$ref,omitempty"`
}
