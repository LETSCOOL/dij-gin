// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dij_gin_test

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	. "github.com/letscool/dij-gin"
	"github.com/letscool/dij-gin/libs"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

type TestByWebServerValue struct {
	WebServer
}

type TestByWebServerPtr struct {
	*WebServer
}

type TestByWebServerValueExt struct {
	WebServer
}

func (w TestByWebServerValueExt) iAmAWebServer() {

}

type TestByWebMiddlewareValue struct {
	WebMiddleware
}

// go test ./ -v -run TestValidateWebType
func TestValidateWebType(t *testing.T) {
	//dij.EnableLog()
	t.Run("TestByWebServerValue", func(t *testing.T) {
		ws := TestByWebServerValue{}
		if typ := reflect.TypeOf(ws); !IsTypeOfWebServer(typ) {
			t.Errorf("%v is not a web server type\n", typ)
			if !IsTypeOfWebController(typ) {
				t.Errorf("%v is not a web controller type\n", typ)
			}
		}
		if typ := reflect.TypeOf(&ws); !IsTypeOfWebServer(typ) {
			t.Errorf("%v is not a web server type\n", typ)
			if !IsTypeOfWebController(typ) {
				t.Errorf("%v is not a web controller type\n", typ)
			}
		}
	})

	t.Run("TestByWebServerPtr", func(t *testing.T) {
		ws := TestByWebServerPtr{}
		if typ := reflect.TypeOf(ws); !IsTypeOfWebServer(typ) {
			t.Errorf("%v is not a web server type\n", typ)
			if !IsTypeOfWebController(typ) {
				t.Errorf("%v is not a web controller type\n", typ)
			}
		}
		if typ := reflect.TypeOf(&ws); !IsTypeOfWebServer(typ) {
			t.Errorf("%v is not a web server type\n", typ)
			if !IsTypeOfWebController(typ) {
				t.Errorf("%v is not a web controller type\n", typ)
			}
		}
	})

	t.Run("TestByWebServerValueExt", func(t *testing.T) {
		ws := TestByWebServerValueExt{}
		if typ := reflect.TypeOf(ws); !IsTypeOfWebServer(typ) {
			t.Errorf("%v is not a web server type\n", typ)
			if !IsTypeOfWebController(typ) {
				t.Errorf("%v is not a web controller type\n", typ)
			}
		}
		if typ := reflect.TypeOf(&ws); !IsTypeOfWebServer(typ) {
			t.Errorf("%v is not a web server type\n", typ)
			if !IsTypeOfWebController(typ) {
				t.Errorf("%v is not a web controller type\n", typ)
			}
		}
	})

	t.Run("TestByWebMiddlewareValue", func(t *testing.T) {
		ws := TestByWebMiddlewareValue{}
		if typ := reflect.TypeOf(ws); !IsTypeOfWebMiddleware(typ) {
			t.Errorf("%v is not a web middleware type\n", typ)
		}
		if typ := reflect.TypeOf(&ws); !IsTypeOfWebMiddleware(typ) {
			t.Errorf("%v is not a web middleware type\n", typ)
		}
	})
}

type TestWebServer struct {
	WebServer `http:"middleware=abc&cors" description:""`

	ctrl1   *TestWebController1     `di:"^"`
	mdl1    *TestWebMiddleware      `di:"^"`
	swagger *libs.SwaggerController `di:""`
	_       *libs.CorsMiddleware    `di:""`
}

func (s *TestWebServer) Get(ctx WebContext) {
	ctx.IndentedJSON(http.StatusOK, "/")
}

func (s *TestWebServer) GetRoot(ctx WebContext) {
	ctx.IndentedJSON(http.StatusOK, "/root")
}

func (s *TestWebServer) GetHello(ctx struct {
	WebContext `http:"middleware=efg" description:"説嗨"`
	a          float32
}) {
	ctx.IndentedJSON(http.StatusOK, fmt.Sprintf("/hello %f", ctx.a))
}

// PostJson shows post request with json style
// curl: curl -X POST http://localhost:8000/json -H 'Content-Type: application/json' -d '{"a":123,"b":"data+b"}'
func (s *TestWebServer) PostJson(ctx struct {
	WebContext `http:"" description:"貼一個Json"`
	json       struct {
		A int    `form:"a" json:"a" binding:"required" http:""`
		B string `form:"b" json:"b" binding:"required"`
	}
}) {
	fmt.Printf("json: %v\n", ctx.json)
	//fmt.Printf("a: %d, b: %s\n", ctx.a, ctx.b)
	ctx.IndentedJSON(http.StatusOK, ctx.json)
}

type TestWebController1 struct {
	WebController `http:"path=user,tag=User"`

	Ctrl2 *TestWebController2 `di:"^"`
}

func (c *TestWebController1) GetUserMe(ctx struct {
	WebContext `http:"me, method=get, middleware=" description:"取得使用者資訊"`
}) (result struct {
	Default *string `http:"200,"`
	Error   error
}) {
	ctx.IndentedJSON(http.StatusOK, "/user")

	return
}

type User struct {
	Uid  int
	Name string
}

func (c *TestWebController1) GetUserById(ctx struct {
	WebContext `http:":id/profile, method=get" description:"取得使用者資訊"`
	Id         int `http:"id" validate:"gte=100,lte=999"`
}) (result struct {
	Data *User `http:"200," description:"使用者資訊"`
}) {
	//ctx.IndentedJSON(http.StatusOK, fmt.Sprintf("/user/%d", ctx.id))
	//fmt.Printf("Id: %d\n", ctx.id)
	result.Data = &User{
		Uid:  ctx.Id,
		Name: "234",
	}
	return
}

type TestWebController2 struct {
	WebController
}

type TestWebMiddleware struct {
	WebMiddleware
}

func (m *TestWebMiddleware) HandleAbc(ctx struct {
	WebContext `http:""`
}) {
	fmt.Printf("*** Hi i am Abc Middleware ***\n")
	//ctx.Next()
}

func (m *TestWebMiddleware) HandleEfg(ctx struct {
	WebContext `http:""`
}) {
	fmt.Printf("*** Hi i am Efg Middleware ***\n")
	//ctx.Next()
}

// go test ./ -v -run TestWebServerExec
func TestWebServerExec(t *testing.T) {
	t.Run("dij", func(t *testing.T) {
		config := NewWebConfig().
			SetAddress("localhost").
			SetRtMode(RtDev).
			SetOpenApi(func(o *OpenApiConfig) {
				o.SetEnabled(true).UseHttpOnly().SetDocPath("doc")
			})
		t.Log(config)
		//dij.EnableLog()
		//wsTyp := reflect.TypeOf(TestWebServer{})
		//if err := LaunchGin(wsTyp, config); err != nil {
		//	t.Error(err)
		//}
		inst := &TestWebServer{}
		if err := LaunchGin(inst, config); err != nil {
			t.Error(err)
		}
	})
}

// go test ./ -v -run TestRegex
func TestRegex(t *testing.T) {
	t.Run("request name and method from function name", func(t *testing.T) {
		re := regexp.MustCompile(`^(get|post|put|patch|delete|head|connect|options|trace)`)
		if v := re.Find([]byte(strings.ToLower("DeleteHello"))); len(v) == 0 {
			t.Error("error regex: ", string(v))
		} else {
			t.Log(string(v))
		}
	})
	t.Run("response code from field name", func(t *testing.T) {
		re := regexp.MustCompile(`^((\w*[\D+|^][2-5]\d{2})|default)$`)
		data := []struct{ data, result string }{
			{"a2345", ""}, {"a345", "a345"}, {"adefault9", ""},
			{"default300", "default300"}, {"DeFault", "default"},
		}
		for i, d := range data {
			v := re.Find([]byte(strings.ToLower(d.data)))
			if d.result != string(v) {
				if len(v) == 0 {
					t.Errorf("%d. %s ==> No match\n", i, d.data)
				} else {
					t.Errorf("%d. %s ==> Found %s\n", i, d.data, string(v))
				}
			}
		}
	})
}

// go test ./ -v -run TestValidator
func TestValidator(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		v := validator.New()
		v.SetTagName("binding")

		a := struct {
			J int    `binding:"gte=0,lte=130"`
			S string `binding:"required"`
		}{}
		if err := v.Struct(&a); err == nil {
			//validationErrors := err.(validator.ValidationErrors)
			t.Errorf("should be error")
		}
		a.S = "1234"
		a.J = 99
		if err := v.Struct(&a); err != nil {
			//validationErrors := err.(validator.ValidationErrors)
			t.Error(err)
		}
	})
	t.Run("Deep", func(t *testing.T) {
		v := validator.New()
		v.SetTagName("binding")

		a := struct {
			J int    `binding:"gte=0,lte=130"`
			S string `binding:"required"`
			O struct {
				SS string `binding:"required"`
			}
		}{}
		if err := v.Struct(&a); err == nil {
			//validationErrors := err.(validator.ValidationErrors)
			t.Errorf("should be error")
		}
		a.S = "1234"
		a.J = 99
		if err := v.Struct(&a); err == nil {
			//validationErrors := err.(validator.ValidationErrors)
			t.Errorf("should be error")
		}
		a.O.SS = "bb"
		if err := v.Struct(&a); err != nil {
			//validationErrors := err.(validator.ValidationErrors)
			t.Error(err)
		}
	})
}
