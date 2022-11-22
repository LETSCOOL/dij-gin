// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package spec

type WebSiteSpec struct {
	Swagger     string                        `json:"swagger"`
	Info        *Info                         `json:"info"`
	Host        string                        `json:"host"`
	BasePath    string                        `json:"basePath,omitempty"`
	Tags        []Tag                         `json:"tags,omitempty"`
	Schemes     []string                      `json:"schemes"`                       // ex: http, https
	Paths       Paths                         `json:"paths,omitempty"`               // ex: "/user":Path
	Security    map[string]SecurityDefinition `json:"securityDefinitions,omitempty"` // ex: {"api_key": SecurityDefinition}
	Definitions map[string]TypeDefinition     `json:"definitions,omitempty"`         // ex: {"type_name":TypeDefinition}
}

func (s *WebSiteSpec) AddMethod(path string, method string, def Method) {
	if s.Paths == nil {
		s.Paths = Paths{}
	}
	p, b := s.Paths[path]
	if !b {
		p = Path{}
		s.Paths[path] = p
	}
	p[method] = def
}

type Info struct {
	License        *License `json:"license,omitempty"`
	Contact        *Contact `json:"contact,omitempty"`
	Description    string   `json:"description"`
	TermsOfService string   `json:"termsOfService,omitempty"`
	Title          string   `json:"title"`
	Version        string   `json:"version,omitempty"`
}

type License struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
type Contact struct {
	Email string `json:"email"`
}

type Tag struct {
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	ExternalDocs *ExternalDoc `json:"externalDocs,omitempty"`
}

type ExternalDoc struct {
	Description string `json:"description"`
	Url         string `json:"url"`
}

type Paths = map[string]Path
type Path map[string]Method // ex: {"post":Method}

type Method struct {
	Summary     string              `json:"summary"`
	Security    []any               `json:"security,omitempty"`
	Consumes    []string            `json:"consumes,omitempty"` // ex: ["application/json", "multipart/form-data"]
	Produces    []string            `json:"produces,omitempty"` // ex: ["application/json"]
	Description string              `json:"description"`
	OperationID string              `json:"operationId,operationId"`
	Parameters  []Parameter         `json:"parameters,operationId"`
	Responses   map[string]Response `json:"responses,operationId"` // ex: {"200":response}, {"default":{"description":"successful operation"}}
	Tags        []string            `json:"tags,operationId"`
}

type Parameter struct {
	In          string      `json:"in"`   // ex. "body"
	Name        string      `json:"name"` // ex. "body"
	Format      string      `json:"format"`
	Description string      `json:"description"`
	Type        string      `json:"type,omitempty"`
	Schema      *TypeSchema `json:"schema,omitempty"`
	Required    bool        `json:"required"`
}

type Response struct {
	Schema      *TypeSchema `json:"schema,omitempty"`
	Description string      `json:"description"`
}

type TypeSchema struct {
	Ref                  string                `json:"$ref"`
	Type                 string                `json:"type,omitempty"` // ex: "array", "object", "string", etc.
	Items                *TypeSchema           `json:"items,omitempty"`
	AdditionalProperties *AdditionalProperties `json:"additionalProperties,omitempty"`
}

type AdditionalProperties struct {
	Type   string `json:"type"`
	Format string `json:"format"`
}

type SecurityDefinition struct {
	Type    string            `json:"type"` // ex: "apiKey", "oauth2"
	Name    string            `json:"name,omitempty"`
	In      string            `json:"in,omitempty"`
	AuthUrl string            `json:"authorizationUrl"`
	Flow    string            `json:"flow"`
	Scopes  map[string]string `json:"scopes"` // ex: {"read:pets": "read your pets","write:pets": "modify pets in your account"}
}

type TypeDefinition struct {
	Type       string              `json:"type"`               // ex: "object"
	Properties map[string]Property `json:"properties"`         // ex: {"property_name": Property}
	Required   []string            `json:"required,omitempty"` // names for required properties
}

type Property struct {
	Type        string   `json:"type,omitempty"`   // ex: "integer", "string", "boolean", "array"
	Format      string   `json:"format,omitempty"` // ex: "int32", "int64", "date-time"
	Description string   `json:"description,omitempty"`
	Enum        []string `json:"enum,omitempty"`
	Ref         string   `json:"$ref,omitempty"`

	//TODO: array definition
}
