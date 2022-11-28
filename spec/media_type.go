// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package spec

// MediaType Each Object provides schema and examples for the media type identified by its key.
//
//	{
//	 "application/json": {
//	   "schema": {
//	        "$ref": "#/components/schemas/Pet"
//	   },
//	   "examples": {
//	     "cat" : {
//	       "summary": "An example of a cat",
//	       "value":
//	         {
//	           "name": "Fluffy",
//	           "petType": "Cat",
//	           "color": "White",
//	           "gender": "male",
//	           "breed": "Persian"
//	         }
//	     },
//	     "dog": {
//	       "summary": "An example of a dog with a cat's name",
//	       "value" :  {
//	         "name": "Puma",
//	         "petType": "Dog",
//	         "color": "Black",
//	         "gender": "Female",
//	         "breed": "Mixed"
//	       },
//	     "frog": {
//	         "$ref": "#/components/examples/frog-example"
//	       }
//	     }
//	   }
//	 }
//	}
type MediaType struct {
	// The schema defining the content of the request, response, or parameter.
	Schema *SchemaR `json:"schema,omitempty"`

	// Example of the media type. The example object SHOULD be in the correct format as specified by the media type.
	// The example field is mutually exclusive of the examples field. Furthermore, if referencing a schema which contains an example,
	// the example value SHALL override the example provided by the schema.
	Example any `json:"example,omitempty"`

	// Examples of the media type. Each example object SHOULD match the media type and specified schema if present.
	// The examples field is mutually exclusive of the example field. Furthermore, if referencing a schema which contains an example,
	// the examples value SHALL override the example provided by the schema.
	Examples map[string]ExampleR `json:"examples,omitempty"`

	// A map between a property name and its encoding information. The key, being the property name,
	// MUST exist in the schema as a property. The encoding object SHALL only apply to requestBody objects
	// when the media type is multipart or application/x-www-form-urlencoded.
	Encoding map[string]Encoding `json:"encoding,omitempty"`
}

type Content map[MediaTypeCoding]MediaType

func (c *Content) SetMediaType(coding MediaTypeCoding, mediaType MediaType) {
	(*c)[coding] = mediaType
}

type MediaTypeCoding string

const (
	UrlEncoded    MediaTypeCoding = "application/x-www-form-urlencoded"
	MultipartForm MediaTypeCoding = "multipart/form-data"
	JsonObject    MediaTypeCoding = "application/json"
	XmlObject     MediaTypeCoding = "application/xml"
)
