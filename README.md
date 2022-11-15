## dij-gin

A dij-style gin library. [gin](https://github.com/gin-gonic/gin) is one of 
most popular web frameworks for golang. [dij](https://github.com/LETSCOOL/lc-go)
stands for dependency injection. This library provides a dij-style gin wrapper.

### Examples

A simple example:
```go
package main

import (
	. "github.com/letscool/dij-gin/executor"
	"log"
	"net/http"
	"reflect"
)

type TWebServer struct {
	WebServer
}

// GetHello a http request with "get" method.
// Url should like this in local: http://localhost:8000/hello
func (s *TWebServer) GetHello(ctx WebContext) {
	ctx.IndentedJSON(http.StatusOK, "/hello")
}

func main() {
	wsTyp := reflect.TypeOf(TWebServer{})
	//dij.EnableLog()
	if err := LaunchGin(wsTyp); err != nil {
		log.Fatalln(err)
	}
}
```


More examples: [go-examples](https://github.com/LETSCOOL/go-examples)