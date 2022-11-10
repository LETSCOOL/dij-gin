package executor

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
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

// go test ./executor -v -run TestValidateWebType
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
	WebServer `http:"" description:""`

	ctrl1      *TestWebController1 `di:"^"`
	Middleware *TestWebMiddleware  `di:"^"`
}

func (s *TestWebServer) SetupRouter(router gin.IRouter) {
	router.GET("", s.GetRoot)
}

func (s *TestWebServer) GetRoot(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, "hello")
}

type TestWebController1 struct {
	WebController `http:"path=user"`

	Ctrl2 *TestWebController2 `di:"^"`
}

func (c *TestWebController1) GetUser(p struct {
	ctx *WebContext `http:"method=get, path=" description="取得使用者資訊"`
}) {
	ctx := p.ctx
	ctx.Get("")
}

type TestWebController2 struct {
	WebController
}

type TestWebMiddleware struct {
	WebMiddleware
}

// go test ./executor -v -run TestWebServerExec
func TestWebServerExec(t *testing.T) {
	t.Run("dij", func(t *testing.T) {
		wsTyp := reflect.TypeOf(TestWebServer{})
		//dij.EnableLog()
		if err := LaunchGin(wsTyp); err != nil {
			t.Error(err)
		}
	})
}
