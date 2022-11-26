// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package spec

// Example In all cases, the example value is expected to be compatible with the type schema of its associated value.
// Tooling implementations MAY choose to validate compatibility automatically, and reject the example value(s) if incompatible.
type Example struct {
	// Short description for the example.
	Summary string `json:"summary,omitempty"`

	// Long description for the example. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// Embedded literal example. The value field and externalValue field are mutually exclusive.
	// To represent examples of media types that cannot naturally be represented in JSON or YAML,
	// use a string value to contain the example, escaping where necessary.
	Value any `json:"value,omitempty"`

	// A URL that points to the literal example. This provides the capability to reference examples that cannot easily be included in JSON or YAML documents. The value field and externalValue field are mutually exclusive.
	ExternalValue string `json:"externalValue,omitempty"`
}

type Examples map[string]ExampleR

// ExampleR presents Example or Ref combination
type ExampleR struct {
	*Example `json:""`
	Ref      string `json:"$ref,omitempty"`
}
