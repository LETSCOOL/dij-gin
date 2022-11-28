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
The http tag includes an attribute "[AttrKey]"

|      AttrKey       | Meaning                           |
|:------------------:|:----------------------------------|
|  form, multipart   | multipart/form-data               |
| urlenc, urlencoded | application/x-www-form-urlencoded |
|        json        | application/json                  |
|        xml         | application/xml                   |


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