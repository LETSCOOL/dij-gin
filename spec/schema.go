// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package spec

import (
	"fmt"
	"reflect"
	"strconv"
)

// The Schema Object allows the definition of input and output data types. These types can be objects, but also primitives and arrays.
// This object is an extended subset of the JSON Schema Specification Wright Draft 00.
// Refer to: https://swagger.io/specification/.
//
// For more information about the properties, see JSON Schema Core and JSON Schema Validation. Unless stated otherwise,
// the property definitions follow the JSON Schema.
// Refer to: https://datatracker.ietf.org/doc/html/draft-wright-json-schema-00 and https://datatracker.ietf.org/doc/html/draft-wright-json-schema-validation-00.
//
// Properties
// The following properties are taken directly from the JSON Schema definition and follow the same specifications:
//
//   - title
//   - multipleOf
//   - maximum
//   - exclusiveMaximum
//   - minimum
//   - exclusiveMinimum
//   - maxLength
//   - minLength
//   - pattern (This string SHOULD be a valid regular expression, according to the Ecma-262 Edition 5.1 regular expression dialect)
//   - maxItems
//   - minItems
//   - uniqueItems
//   - maxProperties
//   - minProperties
//   - required
//   - enum
//
// The following properties are taken from the JSON Schema definition but their definitions were adjusted to the OpenAPI Specification.
//
//   - type - Value MUST be a string. Multiple types via an array are not supported.
//   - allOf - Inline or referenced schema MUST be of a Schema Object and not a standard JSON Schema.
//   - oneOf - Inline or referenced schema MUST be of a Schema Object and not a standard JSON Schema.
//   - anyOf - Inline or referenced schema MUST be of a Schema Object and not a standard JSON Schema.
//   - not - Inline or referenced schema MUST be of a Schema Object and not a standard JSON Schema.
//   - items - Value MUST be an object and not an array. Inline or referenced schema MUST be of a Schema Object and not a standard JSON Schema. items MUST be present if the type is array.
//   - properties - Property definitions MUST be a Schema Object and not a standard JSON Schema (inline or referenced).
//   - additionalProperties - Value can be boolean or object. Inline or referenced schema MUST be of a Schema Object and not a standard JSON Schema. Consistent with JSON Schema, additionalProperties defaults to true.
//   - description - CommonMark syntax MAY be used for rich text representation.
//   - format - See Data Type Formats for further details. While relying on JSON Schema's defined formats, the OAS offers a few additional predefined formats.
//   - default - The default value represents what would be assumed by the consumer of the input as the value of the schema if one is not provided. Unlike JSON Schema, the value MUST conform to the defined type for the Schema Object defined at the same level. For example, if type is string, then default can be "foo" but cannot be 1.
//
// Alternatively, any time a Schema Object can be used, a Reference Object can be used in its place. This allows referencing definitions instead of defining them inline.
//
// Additional properties defined by the JSON Schema specification that are not mentioned here are strictly unsupported.
//
// Other than the JSON Schema subset fields, the following fields MAY be used for further schema documentation.
//
// Primitive sample
//
//	{
//	 "type": "string",
//	 "format": "email"
//	}
//
// Simple model
//
//	{
//	 "type": "object",
//	 "required": [
//	   "name"
//	 ],
//	 "properties": {
//	   "name": {
//	     "type": "string"
//	   },
//	   "address": {
//	     "$ref": "#/components/schemas/Address"
//	   },
//	   "age": {
//	     "type": "integer",
//	     "format": "int32",
//	     "minimum": 0
//	   }
//	 }
//	}
//
// Model with map/dictionary properties - for a simple string to string mapping:
//
//	{
//	 "type": "object",
//	 "additionalProperties": {
//	   "type": "string"
//	 }
//	}
//
// Model with map/dictionary properties - for a string to model mapping:
//
//	{
//	 "type": "object",
//	 "additionalProperties": {
//	   "$ref": "#/components/schemas/ComplexModel"
//	 }
//	}
//
// Model with example
//
//	{
//	 "type": "object",
//	 "properties": {
//	   "id": {
//	     "type": "integer",
//	     "format": "int64"
//	   },
//	   "name": {
//	     "type": "string"
//	   }
//	 },
//	 "required": [
//	   "name"
//	 ],
//	 "example": {
//	   "name": "Puma",
//	   "id": 1
//	 }
//	}
//
// Models with composition
//
//	{
//	 "components": {
//	   "schemas": {
//	     "ErrorModel": {
//	       "type": "object",
//	       "required": [
//	         "message",
//	         "code"
//	       ],
//	       "properties": {
//	         "message": {
//	           "type": "string"
//	         },
//	         "code": {
//	           "type": "integer",
//	           "minimum": 100,
//	           "maximum": 600
//	         }
//	       }
//	     },
//	     "ExtendedErrorModel": {
//	       "allOf": [
//	         {
//	           "$ref": "#/components/schemas/ErrorModel"
//	         },
//	         {
//	           "type": "object",
//	           "required": [
//	             "rootCause"
//	           ],
//	           "properties": {
//	             "rootCause": {
//	               "type": "string"
//	             }
//	           }
//	         }
//	       ]
//	     }
//	   }
//	 }
//	}
//
// Models with polymorphism support
//
//	{
//	 "components": {
//	   "schemas": {
//	     "Pet": {
//	       "type": "object",
//	       "discriminator": {
//	         "propertyName": "petType"
//	       },
//	       "properties": {
//	         "name": {
//	           "type": "string"
//	         },
//	         "petType": {
//	           "type": "string"
//	         }
//	       },
//	       "required": [
//	         "name",
//	         "petType"
//	       ]
//	     },
//	     "Cat": {
//	       "description": "A representation of a cat. Note that `Cat` will be used as the discriminator value.",
//	       "allOf": [
//	         {
//	           "$ref": "#/components/schemas/Pet"
//	         },
//	         {
//	           "type": "object",
//	           "properties": {
//	             "huntingSkill": {
//	               "type": "string",
//	               "description": "The measured skill for hunting",
//	               "default": "lazy",
//	               "enum": [
//	                 "clueless",
//	                 "lazy",
//	                 "adventurous",
//	                 "aggressive"
//	               ]
//	             }
//	           },
//	           "required": [
//	             "huntingSkill"
//	           ]
//	         }
//	       ]
//	     },
//	     "Dog": {
//	       "description": "A representation of a dog. Note that `Dog` will be used as the discriminator value.",
//	       "allOf": [
//	         {
//	           "$ref": "#/components/schemas/Pet"
//	         },
//	         {
//	           "type": "object",
//	           "properties": {
//	             "packSize": {
//	               "type": "integer",
//	               "format": "int32",
//	               "description": "the size of the pack the dog is from",
//	               "default": 0,
//	               "minimum": 0
//	             }
//	           },
//	           "required": [
//	             "packSize"
//	           ]
//	         }
//	       ]
//	     }
//	   }
//	 }
//	}
type Schema struct {
	// ================== JSON Schema ======================
	// JsonSchema presents partial json schema
	// Refer to spec: https://datatracker.ietf.org/doc/html/draft-wright-json-schema-validation-00

	// **JSON**  The value of both of these keywords MUST be a string.
	//   Both of these keywords can be used to decorate a user interface with
	//   information about the data produced by this user interface.  A title
	//   will preferrably be short, whereas a description will provide
	//   explanation about the purpose of the instance described by this
	//   schema.
	Title string `json:"title,omitempty"`

	// **JSON** The value of both of these keywords MUST be a string.
	//   Both of these keywords can be used to decorate a user interface with
	//   information about the data produced by this user interface.  A title
	//   will preferrably be short, whereas a description will provide
	//   explanation about the purpose of the instance described by this
	//   schema.
	//
	// CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// **JSON** The value of this keyword MUST be either a string or an array.  If it
	//   is an array, elements of the array MUST be strings and MUST be unique.
	//   String values MUST be one of the seven primitive types defined by the
	//   core specification.
	//   An instance matches successfully if its primitive type is one of the
	//   types defined by keyword.  Recall: "number" includes "integer".
	//
	// Value MUST be a string. Multiple types via an array are not supported.
	Type string `json:"type,omitempty"`

	// **JSON** Implementations MAY support the "format" keyword.  Should they choose
	//   to do so:
	//      they SHOULD implement validation for attributes defined below;
	//      they SHOULD offer an option to disable validation for this keyword,
	//   Implementations MAY add custom format attributes.  Save for agreement
	//   between parties, schema authors SHALL NOT expect a peer
	//   implementation to support this keyword and/or custom format
	//   attributes.
	//
	// See Data Type Formats for further details. While relying on JSON Schema's defined formats, the OAS offers a few additional predefined formats.
	Format string `json:"format,omitempty"`

	// **JSON** The value of this keyword MUST be an array.  This array SHOULD have
	//   at least one element.  Elements in the array SHOULD be unique.
	//   Elements in the array MAY be of any type, including null.
	//   An instance validates successfully against this keyword if its value
	//   is equal to one of the elements in this keyword's array value.
	Enum string `json:"enum,omitempty"`

	// **JSON** The value of "properties" MUST be an object.  Each value of this
	//   object MUST be an object, and each object MUST be a valid JSON Schema.
	//   If absent, it can be considered the same as an empty object.
	// Property definitions MUST be a Schema Object and not a standard JSON Schema (inline or referenced).
	Properties map[string]SchemaR `json:"properties,omitempty"`

	// **JSON** The value of "additionalProperties" MUST be a boolean or a schema.
	//   If "additionalProperties" is absent, it may be considered present with an empty schema as a value.
	//   If "additionalProperties" is true, validation always succeeds.
	//   If "additionalProperties" is false, validation succeeds only if the instance is an object and
	//   all properties on the instance were covered by "properties" and/or "patternProperties".
	//   If "additionalProperties" is an object, validate the value as a schema to all the properties
	//   that weren't validated by "properties" nor "patternProperties".
	// Value can be boolean or object. Inline or referenced schema MUST be of a Schema Object and not a standard JSON Schema.
	// Consistent with JSON Schema, additionalProperties defaults to true.
	//
	// 所有不列在properties與patternProperties的屬性都要符合這additionProperties規範。
	// 不存在表示任意物件、ture表示皆可、false表示不允許額外的屬性(都需要列在properties與patternProperties裏面)、Schema表示需要符合該定義。
	AdditionalProperties *SchemaR `json:"additionalProperties,omitempty"`

	// **JSON** The value of this keyword MUST be an array. This array MUST have at
	//   least one element. Elements of this array MUST be strings, and MUST be unique.
	//   An object instance is valid against this keyword if its property set
	//   contains all elements in this keyword's array value.
	Required []string `json:"required,omitempty"`

	// **JSON** The value of this keyword MUST be an integer. This integer MUST be greater than, or equal to, 0.
	//   An object instance is valid against "maxProperties" if its number of
	//   properties is less than, or equal to, the value of this keyword.
	MaxProperties any `json:"maxProperties,omitempty"`

	// **JSON** The value of this keyword MUST be an integer. This integer MUST be
	//   greater than, or equal to, 0.
	//   An object instance is valid against "minProperties" if its number of
	//   properties is greater than, or equal to, the value of this keyword.
	//   If this keyword is not present, it may be considered present with a
	//   value of 0.
	MinProperties any `json:"minProperties,omitempty"`

	// **JSON** The value of "additionalItems" MUST be either a boolean or an object.
	//   If it is an object, this object MUST be a valid JSON Schema.
	//   The value of "items" MUST be either a schema or array of schemas.
	//   Successful validation of an array instance in regard to these two
	//   keywords is determined as follows:
	//		if "items" is not present, or its value is an object, validation
	//		of the instance always succeeds, regardless of the value of "additionalItems";
	//		if the value of "additionalItems" is boolean value true or an
	//		object, validation of the instance always succeeds;
	//		if the value of "additionalItems" is boolean value false and the
	//		value of "items" is an array, the instance is valid if its size is
	//		less than, or equal to, the size of "items".
	//		If either keyword is absent, it may be considered present with an
	//		empty schema.
	// Value MUST be an object and not an array. Inline or referenced schema MUST be of a Schema Object and not
	// a standard JSON Schema. items MUST be present if the type is array.
	Items any `json:"items,omitempty"`

	// **JSON** The value of this keyword MUST be an integer. This integer MUST be greater than, or equal to, 0.
	//   An array instance is valid against "maxItems" if its size is less
	//   than, or equal to, the value of this keyword.
	MaxItems any `json:"maxItems,omitempty"`

	// **JSON**  The value of this keyword MUST be an integer. This integer MUST be greater than, or equal to, 0.
	//   An array instance is valid against "minItems" if its size is greater
	//   than, or equal to, the value of this keyword.
	//   If this keyword is not present, it may be considered present with a
	//   value of 0.
	MinItems any `json:"minItems,omitempty"`

	// **JSON** The value of this keyword MUST be a boolean.
	//   If this keyword has boolean value false, the instance validates
	//   successfully.  If it has boolean value true, the instance validates
	//   successfully if all of its elements are unique.
	//   If not present, this keyword may be considered present with boolean
	//   value false.
	UniqueItems bool `json:"uniqueItems,omitempty"`

	// **JSON** There are no restrictions placed on the value of this keyword.
	//   This keyword can be used to supply a default JSON value associated
	//   with a particular schema.  It is RECOMMENDED that a default value be
	//   valid against the associated schema.
	//
	// The default value represents what would be assumed by the consumer of the input as the value of the schema if one is not provided.
	// Unlike JSON Schema, the value MUST conform to the defined type for the Schema Object defined at the same level.
	// For example, if type is string, then default can be "foo" but cannot be 1.
	Default any `json:"default,omitempty"`

	// **JSON** This keyword's value MUST be an array. This array MUST have at least one element.
	//   Elements of the array MUST be objects.  Each object MUST be a valid JSON Schema.
	//   An instance validates successfully against this keyword if it validates successfully against
	//   all schemas defined by this keyword's value.
	// Inline or referenced schema MUST be of a Schema Object and not a standard JSON Schema.
	AllOf []SchemaR `json:"allOf,omitempty"`

	// **JSON** This keyword's value MUST be an array. This array MUST have at least one element.
	//   Elements of the array MUST be objects.  Each object MUST be a valid JSON Schema.
	//   An instance validates successfully against this keyword if it validates successfully against
	//   at least one schema defined by this keyword's value.
	// Inline or referenced schema MUST be of a Schema Object and not a standard JSON Schema.
	AnyOf []SchemaR `json:"anyOf,omitempty"`

	// **JSON** This keyword's value MUST be an array. This array MUST have at least one element.
	//   Elements of the array MUST be objects.  Each object MUST be a valid JSON Schema.
	//   An instance validates successfully against this keyword if it validates successfully against
	//   exactly one schema defined by this keyword's value.
	// Inline or referenced schema MUST be of a Schema Object and not a standard JSON Schema.
	OneOf []SchemaR `json:"oneOf,omitempty"`

	// **JSON** This keyword's value MUST be an object. This object MUST be a valid JSON Schema.
	//   An instance is valid against this keyword if it fails to validate successfully against
	//   the schema defined by this keyword.
	// Inline or referenced schema MUST be of a Schema Object and not a standard JSON Schema.
	Not *SchemaR `json:"not,omitempty"`

	// **JSON** The value of "multipleOf" MUST be a number, strictly greater than 0.
	//   A numeric instance is only valid if division by this keyword's value
	//   results in an integer.
	MultipleOf any `json:"multipleOf,omitempty"`

	// **JSON** The value of "maximum" MUST be a number, representing an upper limit for a numeric instance.
	//   If the instance is a number, then this keyword validates if
	//   "exclusiveMaximum" is true and instance is less than the provided
	//   value, or else if the instance is less than or exactly equal to the
	//   provided value.
	Maximum any `json:"maximum,omitempty"`

	// **JSON** The value of "exclusiveMaximum" MUST be a boolean, representing
	//   whether the limit in "maximum" is exclusive or not.  An undefined
	//   value is the same as false.
	//   If "exclusiveMaximum" is true, then a numeric instance SHOULD NOT be
	//   equal to the value specified in "maximum".  If "exclusiveMaximum" is
	//   false (or not specified), then a numeric instance MAY be equal to the
	//   value of "maximum".
	ExclusiveMaximum bool `json:"exclusiveMaximum,omitempty"`

	// **JSON** The value of "minimum" MUST be a number, representing a lower limit for a numeric instance.
	//   If the instance is a number, then this keyword validates if
	//   "exclusiveMinimum" is true and instance is greater than the provided
	//   value, or else if the instance is greater than or exactly equal to
	//   the provided value.
	Minimum any `json:"minimum,omitempty"`

	// **JSON** The value of "exclusiveMinimum" MUST be a boolean, representing whether
	//   the limit in "minimum" is exclusive or not. An undefined value is the same as false.
	//   If "exclusiveMinimum" is true, then a numeric instance SHOULD NOT be
	//   equal to the value specified in "minimum".  If "exclusiveMinimum" is
	//   false (or not specified), then a numeric instance MAY be equal to the
	//   value of "minimum".
	ExclusiveMinimum bool `json:"exclusiveMinimum,omitempty"`

	// **JSON** The value of this keyword MUST be a non-negative integer.
	//   The value of this keyword MUST be an integer. This integer MUST be greater than, or equal to, 0.
	//   A string instance is valid against this keyword if its length is less
	//   than, or equal to, the value of this keyword.
	//   The length of a string instance is defined as the number of its
	//   characters as defined by RFC 7159 [RFC7159].
	MaxLength any `json:"maxLength,omitempty"`

	// **JSON** A string instance is valid against this keyword if its length is
	//   greater than, or equal to, the value of this keyword.
	//   The length of a string instance is defined as the number of its
	//   characters as defined by RFC 7159 [RFC7159].
	//   The value of this keyword MUST be an integer.  This integer MUST be
	//   greater than, or equal to, 0.
	//   "minLength", if absent, may be considered as being present with
	//   integer value 0.
	MinLength any `json:"minLength,omitempty"`

	// **JSON** The value of this keyword MUST be a string.  This string SHOULD be a
	//   valid regular expression, according to the ECMA 262 regular expression dialect.
	//
	//   A string instance is considered valid if the regular expression matches the instance successfully.
	//   Recall: regular expressions are not implicitly anchored.
	// (This string SHOULD be a valid regular expression, according to the Ecma-262 Edition 5.1 regular expression dialect)
	Pattern string `json:"pattern,omitempty"`

	// ==================== OpenAPI Only ==================

	// A true value adds "null" to the allowed type specified by the type keyword, only if type is explicitly defined within the same Schema Object.
	// Other Schema Object constraints retain their defined behavior, and therefore may disallow the use of null as a value.
	// A false value leaves the specified or default type unmodified. The default value is false.
	Nullable bool `json:"nullable,omitempty"`

	// Adds support for polymorphism. The discriminator is an object name that is used to differentiate between other schemas which may satisfy the payload description.
	// See Composition and Inheritance for more details.
	Discriminator *Discriminator `json:"discriminator,omitempty"`

	// Relevant only for Schema "properties" definitions. Declares the property as "read only".
	// This means that it MAY be sent as part of a response but SHOULD NOT be sent as part of the request.
	// If the property is marked as readOnly being true and is in the required list, the required will take effect on the response only. A property MUST NOT be marked as both readOnly and writeOnly being true. Default value is false.
	ReadOnly bool `json:"readOnly,omitempty"`

	// Relevant only for Schema "properties" definitions. Declares the property as "write only". Therefore, it MAY be sent as part of a request but SHOULD NOT be sent as part of the response. If the property is marked as writeOnly being true and is in the required list, the required will take effect on the request only. A property MUST NOT be marked as both readOnly and writeOnly being true. Default value is false.
	WriteOnly bool `json:"writeOnly,omitempty"`

	// This MAY be used only on properties schemas. It has no effect on root schemas. Adds additional metadata to describe the XML representation of this property.
	Xml *Xml `json:"xml,omitempty"`

	// Additional external documentation for this schema.
	ExternalDocs *ExternalDoc `json:"externalDocs,omitempty"`

	// A free-form property to include an example of an instance for this schema.
	// To represent examples that cannot be naturally represented in JSON or YAML,
	// a string value can be used to contain the example with escaping where necessary.
	Example any `json:"example,omitempty"`

	// Specifies that a schema is deprecated and SHOULD be transitioned out of usage. Default value is false.
	Deprecated bool `json:"deprecated,omitempty"`
}

type Schemas map[string]SchemaR

// SchemaR presents Schema or Ref combination
type SchemaR struct {
	*Schema `json:""`
	Ref     string `json:"$ref,omitempty"`
}

// ApplyType coverts type in golang to json/swagger type
// ref: https://swagger.io/docs/specification/data-models/data-types/
func (s *Schema) ApplyType(t reflect.Type) {
	switch t.Kind() {
	case reflect.Bool:
		s.Type = "boolean"
	case reflect.Int, reflect.Uint:
		s.Type = "integer"
		s.Format = fmt.Sprintf("int%d", strconv.IntSize)
	case reflect.Int8, reflect.Uint8, reflect.Int16, reflect.Uint16, reflect.Int32, reflect.Uint32:
		s.Type = "integer"
		s.Format = "int32"
	case reflect.Int64, reflect.Uint64:
		s.Type = "integer"
		s.Format = "int64"
	case reflect.Float64:
		s.Type = "number"
		s.Format = "double"
	case reflect.Float32:
		s.Type = "number"
		s.Format = "float"
	case reflect.String:
		s.Type = "number"
		s.Format = ""
		// TODO: support datetime
	case reflect.Struct, reflect.Map:
		s.Type = "object"
		s.Format = ""
	case reflect.Array, reflect.Slice:
		s.Type = "array"
		s.Format = ""
		// TODO: set item type ?
	default:
		//
	}
}

func (s *SchemaR) ApplyOneOf(schemas ...SchemaR) {
	s.OneOf = append(s.OneOf, schemas...)
}

func (s *SchemaR) ApplyAllOf(schemas ...SchemaR) {
	s.AllOf = append(s.AllOf, schemas...)
}

func (s *SchemaR) ApplyAnyOf(schemas ...SchemaR) {
	s.AnyOf = append(s.AnyOf, schemas...)
}

func (s *SchemaR) ApplyNot(schema SchemaR) {
	s.Not = &schema
}

func (s *SchemaR) ApplyAdditionalProperties(schema SchemaR) {
	s.AdditionalProperties = &schema
}

func (s *SchemaR) ApplyType(t reflect.Type) {
	s.Schema = &Schema{}
	s.Schema.ApplyType(t)
}
