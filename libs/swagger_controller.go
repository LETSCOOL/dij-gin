// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package libs

import (
	"embed"
	"encoding/json"
	. "github.com/letscool/dij-gin"
	"github.com/letscool/dij-gin/spec"
	"github.com/letscool/lc-go/dij"
	"github.com/letscool/lc-go/io"
	"io/fs"
	"net/http"
)

// content holds our static web server content.
//
//go:embed swagger-ui-dist/4.15.5/*
var content embed.FS

///go:embed swagger_example.json
//var swaggerJson string

// SwaggerController embeds a Swagger/OpenAPI entry.
// Swagger file validator: https://github.com/swagger-api/validator-badge
type SwaggerController struct {
	WebController `http:""`

	ref         *dij.DependencyReference `di:"webserver.dij.ref"`
	openapiSpec *spec.Openapi
}

func (s *SwaggerController) Open(name string) (fs.File, error) {
	if name == "swagger.json" || name == "./swagger.json" {
		// TODO: switch marshal compact or pretty format by debug or production mode
		data, err := json.MarshalIndent(s.openapiSpec, "", "  ")
		if err != nil {
			return nil, err
		}
		return io.NewRoMemFile("swagger.json", data), nil
	}

	fSys, err := fs.Sub(content, "swagger-ui-dist/4.15.5")
	if err != nil {
		return nil, err
	}

	return fSys.Open(name)
}

func (s *SwaggerController) SetupRouter(router WebRouter, _ ...any) {
	if rec, ok := (*(s.ref))[WebSpecRecord]; ok {
		config := (*(s.ref))[WebConfigKey].(*WebConfig)
		s.openapiSpec = rec.(*spec.Openapi)
		router.StaticFS(config.OpenApi.DocPath, http.FS(s))
	}
}
