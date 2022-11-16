package dij_gin

import (
	"fmt"
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
	WebServer `http:"middleware=abc" description:""`

	ctrl1 *TestWebController1 `di:"^"`
	mdl1  *TestWebMiddleware  `di:"^"`
}

func (s *TestWebServer) Get(ctx WebContext) {
	ctx.IndentedJSON(http.StatusOK, "/")
}

func (s *TestWebServer) GetRoot(ctx WebContext) {
	ctx.IndentedJSON(http.StatusOK, "/root")
}

func (s *TestWebServer) GetHello(ctx struct {
	WebContext `http:"" description:""`
	a          float32
}) {
	ctx.IndentedJSON(http.StatusOK, fmt.Sprintf("/hello %f", ctx.a))
}

// PostJson shows post request with json style
// curl: curl -X POST http://localhost:8000/json -H 'Content-Type: application/json' -d '{"a":123,"b":"data+b"}'
func (s *TestWebServer) PostJson(ctx struct {
	WebContext `http:"" description:""`
	json       struct {
		A int    `form:"a" json:"a" binding:"required"`
		B string `form:"b" json:"b" binding:"required"`
	}
}) {
	fmt.Printf("json: %v\n", ctx.json)
	//fmt.Printf("a: %d, b: %s\n", ctx.a, ctx.b)
	ctx.IndentedJSON(http.StatusOK, ctx.json)
}

type TestWebController1 struct {
	WebController `http:"path=user"`

	Ctrl2 *TestWebController2 `di:"^"`
}

func (c *TestWebController1) GetUserMe(ctx struct {
	WebContext `http:"me, method=get, middleware=" description="取得使用者資訊"`
}) {
	ctx.IndentedJSON(http.StatusOK, "/user")
}

func (s *TestWebController1) GetUserById(ctx struct {
	WebContext `http:":id/profile, method=get" description="取得使用者資訊"`
	id         int
}) (result struct {
	data *string `resp:"200,"`
}) {
	ctx.IndentedJSON(http.StatusOK, fmt.Sprintf("/user/%d", ctx.id))
	fmt.Printf("Id: %d\n", ctx.id)
	a := "234"
	result.data = &a
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

// go test ./ -v -run TestWebServerExec
func TestWebServerExec(t *testing.T) {
	t.Run("dij", func(t *testing.T) {
		wsTyp := reflect.TypeOf(TestWebServer{})
		//dij.EnableLog()
		if err := LaunchGin(wsTyp); err != nil {
			t.Error(err)
		}
	})
}

// go test ./ -v -run TestRegex
func TestRegex(t *testing.T) {
	t.Run("request name and method", func(t *testing.T) {
		re := regexp.MustCompile(`^(get|post|put|patch|delete|head|connect|options|trace)`)
		if v := re.Find([]byte(strings.ToLower("DeleteHello"))); len(v) == 0 {
			t.Error("error regex: ", string(v))
		} else {
			t.Log(string(v))
		}
	})
}