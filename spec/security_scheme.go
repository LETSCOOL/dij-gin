// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package spec

type SecurityScheme struct {
}

type SecuritySchemes map[string]SecuritySchemeR

// SecuritySchemeR presents SecurityScheme or Ref combination
type SecuritySchemeR struct {
	*SecurityScheme `json:""`
	Ref             string `json:"$ref,omitempty"`
}
