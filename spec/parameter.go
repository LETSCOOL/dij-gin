// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package spec

import (
	"reflect"
)

// Parameter Describes a single operation parameter.
//
// A unique parameter is defined by a combination of a name and location.
//
// Parameter Locations
// There are four possible parameter locations specified by the in field:
//
//   - path - Used together with Path Templating, where the parameter value is actually part of the operation's URL.
//     This does not include the host or base path of the API. For example, in /items/{itemId}, the path parameter is itemId.
//   - query - Parameters that are appended to the URL. For example, in /items?id=###, the query parameter is id.
//   - header - Custom headers that are expected as part of the request. Note that RFC7230 states header names are case-insensitive.
//   - cookie - Used to pass a specific cookie value to the API.
//
// A header parameter with an array of 64 bit integer numbers:
//
//	{
//	  "name": "token",
//	  "in": "header",
//	  "description": "token to be passed as a header",
//	  "required": true,
//	  "schema": {
//	    "type": "array",
//	    "items": {
//	      "type": "integer",
//	      "format": "int64"
//	    }
//	  },
//	  "style": "simple"
//	}
//
// A path parameter of a string value:
//
//	{
//	 "name": "username",
//	 "in": "path",
//	 "description": "username to fetch",
//	 "required": true,
//	 "schema": {
//	   "type": "string"
//	 }
//	}
//
// An optional query parameter of a string value, allowing multiple values by repeating the query parameter:
//
//	{
//	 "name": "id",
//	 "in": "query",
//	 "description": "ID of the object to fetch",
//	 "required": false,
//	 "schema": {
//	   "type": "array",
//	   "items": {
//	     "type": "string"
//	   }
//	 },
//	 "style": "form",
//	 "explode": true
//	}
//
// A free-form query parameter, allowing undefined parameters of a specific type:
//
//	{
//	 "in": "query",
//	 "name": "freeForm",
//	 "schema": {
//	   "type": "object",
//	   "additionalProperties": {
//	     "type": "integer"
//	   },
//	 },
//	 "style": "form"
//	}
//
// A complex parameter using content to define serialization:
//
//	{
//	 "in": "query",
//	 "name": "coordinates",
//	 "content": {
//	   "application/json": {
//	     "schema": {
//	       "type": "object",
//	       "required": [
//	         "lat",
//	         "long"
//	       ],
//	       "properties": {
//	         "lat": {
//	           "type": "number"
//	         },
//	         "long": {
//	           "type": "number"
//	         }
//	       }
//	     }
//	   }
//	 }
//	}
type Parameter struct {
	// REQUIRED. The name of the parameter. Parameter names are case-sensitive.
	//   - If in is "path", the name field MUST correspond to a template expression occurring within the path field in the Paths Object. See Path Templating for further information.
	//   - If in is "header" and the name field is "Accept", "Content-Type" or "Authorization", the parameter definition SHALL be ignored.
	//   - For all other cases, the name corresponds to the parameter name used by the in property.
	Name string `json:"name"` // ex. "body"

	// REQUIRED. The location of the parameter. Possible values are "query", "header", "path" or "cookie".
	In string `json:"in"` // ex. "body"

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

// WithSchemaAndStyle The rules for serialization of the parameter are specified in one of two ways.
// For simpler scenarios, a schema and style can describe the structure and syntax of the parameter.
type WithSchemaAndStyle struct {
	// Describes how the parameter value will be serialized depending on the type of the parameter value.
	// Default values (based on value of in): for query - form; for path - simple; for header - simple; for cookie - form.
	// In order to support common ways of serializing simple parameters, a set of style values are defined.
	//
	//  | style         |type                     |in          |Comments
	//  | ----------------------------------------------------------------------------------------------
	//  | matrix        |primitive, array, object |path        |Path-style parameters defined by RFC6570
	//  | label         |primitive, array, object |path        |Label style parameters defined by RFC6570
	//  | form          |primitive, array, object |query,cookie|Form style parameters defined by RFC6570. This option replaces collectionFormat with a csv (when explode is false) or multi (when explode is true) value from OpenAPI 2.0.
	//  | simple        |array                    |path,header |Simple style parameters defined by RFC6570. This option replaces collectionFormat with a csv value from OpenAPI 2.0.
	//  | spaceDelimited|array                    |query       |Space separated array values. This option replaces collectionFormat equal to ssv from OpenAPI 2.0.
	//  | pipeDelimited |array                    |query       |Pipe separated array values. This option replaces collectionFormat equal to pipes from OpenAPI 2.0.
	//  | deepObject    |object                   |query       |Provides a simple way of rendering nested objects using form parameters.
	Style string `json:"style,omitempty"`

	// When this is true, parameter values of type array or object generate separate parameters for each value of the array or key-value pair of the map.
	// For other types of parameters this property has no effect. When style is form, the default value is true. For all other styles, the default value is false.
	Explode bool `json:"explode,omitempty"`

	// Determines whether the parameter value SHOULD allow reserved characters, as defined by RFC3986 :/?#[]@!$&'()*+,;= to be included without percent-encoding.
	// This property only applies to parameters with an in value of query. The default value is false.
	AllowReserved bool `json:"allowReserved,omitempty"`

	// The schema defining the type used for the parameter.
	Schema *SchemaR `json:"schema,omitempty"`

	// Example of the parameter's potential value. The example SHOULD match the specified schema and encoding properties if present.
	// The example field is mutually exclusive of the examples field. Furthermore, if referencing a schema that contains an example,
	// the example value SHALL override the example provided by the schema. To represent examples of media types
	// that cannot naturally be represented in JSON or YAML, a string value can contain the example with escaping where necessary.
	Example any `json:"example,omitempty"`

	// Examples of the parameter's potential value. Each example SHOULD contain a value in the correct
	// format as specified in the parameter encoding. The examples field is mutually exclusive of the example field.
	// Furthermore, if referencing a schema that contains an example, the examples value SHALL override the example provided by the schema.
	Examples map[string]ExampleR `json:"examples,omitempty"`
}

// WithContent For more complex scenarios, the content property can define the media type and schema of the parameter.
// A parameter MUST contain either a schema property, or a content property, but not both.
// When example or examples are provided in conjunction with the schema object, the example MUST follow the prescribed serialization strategy for the parameter.
type WithContent struct {
	Content Content `json:"content,omitempty"`
}

func (p *Parameter) ApplyType(t reflect.Type) {
	p.Schema = &SchemaR{}
	p.Schema.ApplyType(t)
}

// ParameterR presents Parameter or Ref combination
type ParameterR struct {
	*Parameter `json:""`
	Ref        string `json:"$ref,omitempty"`
}

// Parameters presents a map for key-parameter pairs.
type Parameters map[string]ParameterR

// ParameterList presents a list/slice of parameters.
type ParameterList []ParameterR

func (l ParameterList) AppendRef(ref string) ParameterList {
	return append(l, ParameterR{Ref: ref})
}

func (l ParameterList) AppendParam(param *Parameter) ParameterList {
	return append(l, ParameterR{Parameter: param})
}
