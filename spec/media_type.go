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

type Content map[MediaTypeTitle]MediaType

func (c *Content) SetMediaType(coding MediaTypeTitle, mediaType MediaType) {
	(*c)[coding] = mediaType
}

type MediaTypeTitle string

const (
	UrlEncoded    MediaTypeTitle = "application/x-www-form-urlencoded"
	MultipartForm MediaTypeTitle = "multipart/form-data"
	PlainText     MediaTypeTitle = "text/plain"
	JsonObject    MediaTypeTitle = "application/json"
	XmlObject     MediaTypeTitle = "application/xml"
	HtmlPage      MediaTypeTitle = "text/html"
	OctetStream   MediaTypeTitle = "application/octet-stream"
	PngImage      MediaTypeTitle = "image/png"
	JpegImage     MediaTypeTitle = "image/jpeg"
)

type MediaTypeKind int

const (
	UnsupportedMediaType MediaTypeKind = iota
	PlainMediaType
	ObjectiveMediaType
	StreamMediaType
)

type MediaTypeSupport struct {
	Abbr  []string
	Title MediaTypeTitle
	Kind  MediaTypeKind
	Req   bool // request supports this Title
	Resp  bool // response supports this Title
}

var mediaTypeSupportList []MediaTypeSupport
var mediaTypeSupports map[string]MediaTypeSupport

func init() {
	mediaTypeSupportList = []MediaTypeSupport{
		{
			Abbr:  []string{"urlenc", "urlencoded"},
			Title: UrlEncoded,
			Kind:  PlainMediaType,
			Req:   true,
			Resp:  true,
		},
		{
			Abbr:  []string{"form", "multipart"},
			Title: MultipartForm,
			Kind:  PlainMediaType,
			Req:   true,
			Resp:  false,
		},
		{
			Abbr:  []string{"plain"},
			Title: PlainText,
			Kind:  PlainMediaType,
			Req:   false,
			Resp:  true,
		},
		{
			Abbr:  []string{"json"},
			Title: JsonObject,
			Kind:  ObjectiveMediaType,
			Req:   true,
			Resp:  true,
		},
		{
			Abbr:  []string{"xml"},
			Title: XmlObject,
			Kind:  ObjectiveMediaType,
			Req:   true,
			Resp:  true,
		},
		{
			Abbr:  []string{"html", "page"},
			Title: HtmlPage,
			Kind:  ObjectiveMediaType,
			Req:   false,
			Resp:  true,
		},
		{
			Abbr:  []string{"octet", "stream"},
			Title: OctetStream,
			Kind:  StreamMediaType,
			Req:   false,
			Resp:  true,
		},
		{
			Abbr:  []string{"png"},
			Title: PngImage,
			Kind:  StreamMediaType,
			Req:   false,
			Resp:  true,
		},
		{
			Abbr:  []string{"jpeg"},
			Title: JpegImage,
			Kind:  StreamMediaType,
			Req:   false,
			Resp:  true,
		},
	}

	mediaTypeSupports = map[string]MediaTypeSupport{}

	for _, s := range mediaTypeSupportList {
		for _, abbr := range s.Abbr {
			mediaTypeSupports[abbr] = s
		}
	}
}

// IsSupportedMediaType gets media type information.
func IsSupportedMediaType(abbr string) (kind MediaTypeKind, title MediaTypeTitle, supportReq bool, supportResp bool) {
	if support, ok := mediaTypeSupports[abbr]; ok {
		return support.Kind, support.Title, support.Req, support.Resp
	}
	kind = UnsupportedMediaType
	return
}

func GetSupportedMediaType(abbr string) (support MediaTypeSupport, ok bool) {
	support, ok = mediaTypeSupports[abbr]
	return
}
