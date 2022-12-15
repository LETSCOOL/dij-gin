// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package spec

// SecurityScheme describes security scheme for OpenAPI document.
// https://swagger.io/docs/specification/authentication/
type SecurityScheme struct {
	Type             string `json:"type"`
	Scheme           string `json:"scheme,omitempty"`
	BearerFormat     string `json:"bearerFormat,omitempty"`
	In               InWay  `json:"in,omitempty"`
	Name             string `json:"name,omitempty"`
	OpenIdConnectUrl string `json:"openIdConnectUrl,omitempty"`
}

type SecuritySchemes map[string]SecuritySchemeR

// SecuritySchemeR presents SecurityScheme or Ref combination
type SecuritySchemeR struct {
	*SecurityScheme `json:""`
	Ref             string `json:"$ref,omitempty"`
}

func (s *SecuritySchemes) AppendScheme(name string, scheme SecurityScheme) *SecuritySchemes {
	(*s)[name] = SecuritySchemeR{
		SecurityScheme: &scheme,
	}
	return s
}

func (s *SecuritySchemes) AppendBasicAuth(name string) *SecuritySchemes {
	return s.AppendScheme(name, SecurityScheme{
		Type:   "http",
		Scheme: "basic",
	})
}

func (s *SecuritySchemes) AppendBearerAuth(name string) *SecuritySchemes {
	return s.AppendScheme(name, SecurityScheme{
		Type:         "http",
		Scheme:       "bearer",
		BearerFormat: "JWT",
	})
}

// AppendApiKeyAuth appends an api-key scheme with name, this scheme get api-key by paramName from paramIn way.
// The way paramIn can only be "header", "query" or "cookie" (aka. InHeaderWay, InQueryWay, InCookieWay).
func (s *SecuritySchemes) AppendApiKeyAuth(name string, paramIn InWay, paramName string) *SecuritySchemes {
	return s.AppendScheme(name, SecurityScheme{
		Type: "apiKey",
		In:   paramIn,
		Name: paramName,
	})
}

func (s *SecuritySchemes) AppendOpenId(name string, openIdConnectUrl string) *SecuritySchemes {
	return s.AppendScheme(name, SecurityScheme{
		Type:             "openIdConnect",
		OpenIdConnectUrl: openIdConnectUrl,
	})
}
