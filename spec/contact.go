// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package spec

// Contact information for the exposed API.
//
//	{
//	  "name": "API Support",
//	  "url": "http://www.example.com/support",
//	  "email": "support@example.com"
//	}
type Contact struct {
	// The identifying name of the contact person/organization.
	Name string `json:"name"`
	// The URL pointing to the contact information. MUST be in the format of a URL.
	Url string `json:"url"`
	// The email address of the contact person/organization. MUST be in the format of an email address.
	Email string `json:"email"`
}
