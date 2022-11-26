// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package spec

// License information for the exposed API.
//
//	{
//	 "name": "Apache 2.0",
//	 "url": "https://www.apache.org/licenses/LICENSE-2.0.html"
//	}
type License struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
