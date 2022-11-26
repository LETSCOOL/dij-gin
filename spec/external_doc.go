// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package spec

// ExternalDoc allows referencing an external resource for extended documentation.
type ExternalDoc struct {
	// A short description of the target documentation. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description"`

	// REQUIRED. The URL for the target documentation. Value MUST be in the format of a URL.
	Url string `json:"url"`
}
