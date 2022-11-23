// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package libs

import (
	"embed"
	"encoding/json"
	. "github.com/letscool/dij-gin"
	"github.com/letscool/lc-go/io"
	"io/fs"
	"net/http"
)

// content holds our static web server content.
//
//go:embed swagger-ui-dist/4.15.5/*
var content embed.FS

//go:embed swagger_example.json
var swaggerJson string

// SwaggerController embeds a Swagger/OpenAPI entry.
// Swagger file validator: https://github.com/swagger-api/validator-badge
type SwaggerController struct {
	WebController `http:""`
}

func (s *SwaggerController) Open(name string) (fs.File, error) {
	if name == "swagger.json" || name == "./swagger.json" {
		data, err := json.Marshal(s.Spec)
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
	//http.FileSystem()
	//fSys, err := fs.Sub(content, "swagger-ui-dist/4.15.5")
	//if err != nil {
	//	log.Fatal(err)
	//}
	router.StaticFS("docs", http.FS(s))
	//router.GET(".docs/swagger.json", func(ctx *gin.Context) {
	//	ctx.String(200, swaggerJson)
	//})

}
