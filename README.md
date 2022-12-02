## dij-gin

A dij-style gin library. [gin](https://github.com/gin-gonic/gin) is one of 
most popular web frameworks for golang. [dij](https://github.com/LETSCOOL/lc-go)
stands for dependency injection. This library provides a dij-style gin wrapper.

### Examples

A simple example:
```go
package main

import (
	. "github.com/letscool/dij-gin"
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

A query example:
```go
package main

import (
	"fmt"
	. "github.com/letscool/dij-gin"
	"log"
	"net/http"
	"reflect"
)

type TWebServer struct {
	WebServer
}

// GetHello a http request with "get" method.
// Url should like this in local: http://localhost:8000/hello?name=wayne&age=123
func (s *TWebServer) GetHello(ctx struct {
	WebContext
	name string
	age  int
}) {
	//fmt.Printf("%s", ctx.Query("name"))
	ctx.IndentedJSON(http.StatusOK, fmt.Sprintf("/hello %s, %d years old", ctx.name, ctx.age))
}

func main() {
	wsTyp := reflect.TypeOf(TWebServer{})
	//dij.EnableLog()
	if err := LaunchGin(wsTyp); err != nil {
		log.Fatalln(err)
	}
}
```

### Http Tag

#### Attributes
##### Coding/Media Type for Request Input
The http tag includes an attribute "[AttrKey]" for request and response body.
and "mime=[MIME_TYPE]" for response body only.

|      AttrKey       | Req/Resp | MIME Type                         |
|:------------------:|:--------:|:----------------------------------|
|  form, multipart   |   Req    | multipart/form-data               |
| urlenc, urlencoded |   BOTH   | application/x-www-form-urlencoded |
|        json        |   Both   | application/json                  |
|        xml         |   Both   | application/xml                   |
|       plain        |   Resp   | text/plain                        |
|     page, html     |   Resp   | text/html                         |
|       octet        |   Resp   | application/octet-stream          |
|     jpeg, png      |   Resp   | image/jpeg,png                    |


#### Data way for Request Input Variables
The http tag includes an attribute "in=[AttrKey]"

| AttrKey | Default situation                    | Meaning |
|:-------:|:-------------------------------------|:--------|
| header  |                                      |         |
| cookie  |                                      |         |
|  path   | If variable name is included in path |         |
|  query  |                                      |         |
|  body   |                                      |         |

More examples: [go-examples](https://github.com/LETSCOOL/go-examples)