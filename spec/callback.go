// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package spec

// Callback A map of possible out-of band callbacks related to the parent operation.
// Each value in the map is a Path Item Object that describes a set of requests that may be initiated by
// the API provider and the expected responses. The key value used to identify the path item object is an expression,
// evaluated at runtime, that identifies a URL to use for the callback operation.
//
// Not implement yet, refer to https://swagger.io/specification/
type Callback struct {
}

type Callbacks map[string]CallbackR

// CallbackR presents Callback or Ref combination
type CallbackR struct {
	*Callback `json:""`
	Ref       string `json:"$ref,omitempty"`
}
