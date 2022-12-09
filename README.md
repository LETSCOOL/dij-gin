# dij-gin

A dij-style gin library. [gin](https://github.com/gin-gonic/gin) is one of 
most popular web frameworks for golang. [dij](https://github.com/LETSCOOL/lc-go)
stands for dependency injection. This library provides a dij-style gin wrapper.

## Contents
__________
- [Gin Style](#gin-style)
  - [Get method](#get-method)
  - [Hierarchy](#hierarchycontroller)
- [dij-gin Style](#dij-gin-style)
  - [Query](#query)
  - [Where variable data came from?](#where-variable-data-came-from)
  - Customize path name and http method
  - Body
    - Form
    - Json
  - [Validator](#validator)
  - [Response](#response)
  - [Middlewares](#middlewares)
    - [Log](#log) 
    - [Basic Auth](#basic-auth)
    - [CORS](#cors)
  - [OpenAPI generation](#openapi-generation)
    - tag/group
  - Runtime environment
- [Http tag](#http-tag)
- [TODO List](#todo-list)


## Gin Style
___
### Get method
The *WebContext* embeds *gin.Context*, any gin helper functions can be used directly.
```go
package main

import (
	. "github.com/letscool/dij-gin"
	"log"
	"net/http"
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
	//dij.EnableLog()
	if err := LaunchGin(&TWebServer{}); err != nil {
		log.Fatalln(err)
	}
}
```
The function name *GetHello* should combine a method and a path name.
Method should be one of valid http methods, for examples: **Get**, **Post**, **Delete**, etc.

### Hierarchy/Controller
You can group a few http functions in one controller, 
which likes what *gin.router.Group* function does.

```go
package main

import (
	. "github.com/letscool/dij-gin"
	"log"
	"net/http"
)

type TWebServer struct {
	WebServer

	userCtl *TUserController `di:""` // inject dependency by class/default name
}

type TUserController struct {
	WebController `http:"user"` //group by 'user' path
}

// Get a http request with "get" method.
// Url should like this in local: http://localhost:8000/user
func (u *TUserController) Get(ctx WebContext) {
	ctx.IndentedJSON(http.StatusOK, "/user")
}

// GetMe a http request with "get" method.
// Url should like this in local: http://localhost:8000/user/me
func (u *TUserController) GetMe(ctx WebContext) {
	ctx.IndentedJSON(http.StatusOK, "/user/me")
}

func main() {
	//dij.EnableLog()
	if err := LaunchGin(&TWebServer{}); err != nil {
		log.Fatalln(err)
	}
}
```

## dij-gin Style
_____
dij-gin style includes many features:
- Easy way to retrieve path parameters, query parameters, etc.
- Easy way to response data with different http code.
- Support OpenAPI(v3.0) generation.  
- Http functions support to enable or disable by runtime environment, aka: prod, dev, and test.


### Query
```go
package main

import (
  "fmt"
  . "github.com/letscool/dij-gin"
  "log"
  "net/http"
)

type TWebServer struct {
  WebServer
}

// GetHello a http request with "get" method.
// Url should like this in local: http://localhost:8000/hello?name=wayne&age=123.
// The result will be:
//
//	"/hello wayne, 123 years old"
func (s *TWebServer) GetHello(ctx struct {
  WebContext
  Name string `http:"name"`
  Age  int    `http:"age"`
}) {
  ctx.IndentedJSON(http.StatusOK, fmt.Sprintf("/hello %s, %d years old", ctx.Name, ctx.Age))
}

func main() {
  if err := LaunchGin(&TWebServer{}); err != nil {
    log.Fatalln(err)
  }
}
```

### Where variable data came from?

Add an attribute "in=xxx" in http tag. About http tag setting, see
the [reference](#http-tag)

```go
package main

import (
	"fmt"
	. "github.com/letscool/dij-gin"
	"log"
	"net/http"
)

type TWebServer struct {
	WebServer

	userCtl *TUserController `di:""` // inject dependency by class/default name
}

type TUserController struct {
	WebController `http:"user"` //group by 'user' path
}

// PutUserById a http request with "get" method.
// Curl this url should like this in local:
//
//	curl -d "age=34&name=wayne" -X PUT http://localhost:8000/user/2345/profile
//
// The result should be:
//
//	"update user(#2345)'s name wayne and age 34"
func (u *TUserController) PutUserById(ctx struct {
	WebContext `http:":id/profile"`
	Id         int    `http:"id,in=path"`
	Name       string `http:"name,in=body"`
	Age        int    `http:"age,in=body"`
}) {
	ctx.IndentedJSON(http.StatusOK, fmt.Sprintf("update user(#%d)'s name %s and age %d", ctx.Id, ctx.Name, ctx.Age))
}

func main() {
	if err := LaunchGin(&TWebServer{}); err != nil {
		log.Fatalln(err)
	}
}
```

### Validator

dij-gin uses [go-playground/validator/v10](https://github.com/go-playground/validator) for validation.
gin uses 'binding' as validation tag key instead of 'validate', 'validate' tag key is chose from go-playground/validator/v10 official.
dij-gin still uses 'validate' tag key.

```go
package main

import (
	"fmt"
	. "github.com/letscool/dij-gin"
	"log"
	"net/http"
)

type TWebServer struct {
	WebServer

	userCtl *TUserController `di:""` // inject dependency by class/default name
}

type TUserController struct {
	WebController `http:"user"` //group by 'user' path
}

// GetUserById a http request with "get" method.
// Url should like this in local: http://localhost:8000/user/2345/profile.
// The result will be:
//
//	{"message":"Key: 'Id' Error:Field validation for 'Id' failed on the 'lte' tag","code":"400"}
func (u *TUserController) GetUserById(ctx struct {
	WebContext `http:":id/profile"`
	Id         int `http:"id,in=path" validate:"gte=100,lte=999"`
}) {
	ctx.IndentedJSON(http.StatusOK, fmt.Sprintf("get user(#%d)'s profile", ctx.Id))
}

func main() {
	if err := LaunchGin(&TWebServer{}); err != nil {
		log.Fatalln(err)
	}
}
```

### Response

```go
package main

import (
  "errors"
  . "github.com/letscool/dij-gin"
  "log"
)

type TWebServer struct {
  WebServer
}

// GetResp a http request with "get" method.
// Url should like this in local: http://localhost:8000/resp?select=1 .
// Use *curl -v* command to see response code.
func (s *TWebServer) GetResp(ctx struct {
  WebContext
  Select int `http:"select"`
}) (result struct {
  Ok200       *string // the range of last three characters is between 2xx and 5xx, so the response code = 200
  Ok          *string `http:"201"` // force response code to 201
  Redirect302 *string // redirect data should be string type, because it is redirect location.
  Error       error   // default response code for error is 400
}) {
  switch ctx.Select {
  case 1:
    data := "ok"
    result.Ok200 = &data
  case 2:
    data := "ok"
    result.Ok = &data
  case 3:
    url := "https://github.com/letscool"
    result.Redirect302 = &url
  default:
    result.Error = errors.New("an error")
  }
  return
}

func main() {
  if err := LaunchGin(&TWebServer{}); err != nil {
    log.Fatalln(err)
  }
}
```

### Middlewares

#### Log
- Log all http methods for a controller and all it's sub-controllers
```go
package main

import (
  "fmt"
  "github.com/gin-gonic/gin"
  . "github.com/letscool/dij-gin"
  "github.com/letscool/dij-gin/libs"
  "log"
  "net/http"
)

type TWebServer struct {
  WebServer `http:",middleware=log"`

  _ *libs.LogMiddleware `di:""`
}

// GetHello a http request with "get" method.
// Url should like this in local: http://localhost:8000/hello
func (t *TWebServer) GetHello(ctx WebContext) {
  ctx.IndentedJSON(http.StatusOK, "/hello")
}

func main() {
  //f, _ := os.Create("gin.log") // log to file
  config := NewWebConfig().
          //SetDefaultWriter(io.MultiWriter(f)).
          SetDependentRef(libs.RefKeyForLogFormatter, (gin.LogFormatter)(func(params gin.LogFormatterParams) string {
            // your custom format
            return fmt.Sprintf("[%s-%s] \"%s %s\"\n",
              params.ClientIP,
              params.TimeStamp.Format("15:04:05.000"),
              params.Method,
              params.Path,
            )
          }))
  if err := LaunchGin(&TWebServer{}, config); err != nil {
    log.Fatalln(err)
  }
}
```

- Log functions only which set *log* middleware
```go
package main

import (
  "fmt"
  "github.com/gin-gonic/gin"
  . "github.com/letscool/dij-gin"
  "github.com/letscool/dij-gin/libs"
  "log"
  "net/http"
)

type TWebServer struct {
  WebServer `http:""`

  _ *libs.LogMiddleware `di:""`
}

// GetHelloWithLog a http request with "get" method.
// Url should like this in local: http://localhost:8000/hello_with_log
func (t *TWebServer) GetHelloWithLog(ctx struct {
  WebContext `http:"hello_with_log,middleware=log"`
}) {
  ctx.IndentedJSON(http.StatusOK, "hello with log")
}

// GetHelloWithoutLog a http request with "get" method.
// Url should like this in local: http://localhost:8000/hello_without_log
func (t *TWebServer) GetHelloWithoutLog(ctx struct {
  WebContext `http:"hello_without_log"`
}) {
  ctx.IndentedJSON(http.StatusOK, "hello without log")
}

func main() {
  //f, _ := os.Create("gin.log") // log to file
  config := NewWebConfig().
          //SetDefaultWriter(io.MultiWriter(f)).
          SetDependentRef(libs.RefKeyForLogFormatter, (gin.LogFormatter)(func(params gin.LogFormatterParams) string {
            // your custom format
            return fmt.Sprintf("[%s-%s] \"%s %s\"\n",
              params.ClientIP,
              params.TimeStamp.Format("15:04:05.000"),
              params.Method,
              params.Path,
            )
          }))
  if err := LaunchGin(&TWebServer{}, config); err != nil {
    log.Fatalln(err)
  }
}
```

#### Basic Auth

```go
package main

import (
  "crypto/subtle"
  "encoding/base64"
  . "github.com/letscool/dij-gin"
  "github.com/letscool/dij-gin/libs"
  "log"
)

type TWebServer struct {
  WebServer

  _ *TUserController `di:""`
}

type TUserController struct {
  WebController `http:"user"`

  _ *libs.BasicAuthMiddleware `di:""`
}

// GetMe a http request with "get" method.
// Url should like this in local: http://localhost:8000/user/me
func (u *TUserController) GetMe(ctx struct {
  WebContext `http:",middleware=basic_auth"`
}) (result struct {
  Account *Account `http:"200,json"`
}) {
  result.Account = ctx.MustGet(libs.BasicAuthUserKey).(*Account)
  return
}

func main() {
  ac := &FakeAccountDb{}
  ac.initFakeDb()
  config := NewWebConfig().
    SetDependentRef(libs.RefKeyForBasicAuthAccountCenter, ac)
  if err := LaunchGin(&TWebServer{}, config); err != nil {
    log.Fatalln(err)
  }
}
```

Account Db information
```go
type Account struct {
	User  string `json:"user"`
	Email string `json:"email"`
	pass  string
	realm string
}

// FakeAccountDb should implement libs.AccountForBasicAuth interface
type FakeAccountDb struct {
	accounts []Account
	creds    map[string]*Account
}

func (a *FakeAccountDb) initFakeDb() {
	a.accounts = []Account{
		{"john", "john@fake.com", "abc", ""},
		{"wayne", "wayne@fake.com", "abc", ""},
	}

	a.creds = map[string]*Account{}
	for i := range a.accounts {
		account := &a.accounts[i]
		base := account.User + ":" + account.pass
		cred := "Basic " + base64.StdEncoding.EncodeToString([]byte(base))
		a.creds[cred] = account
	}
}

func (a *FakeAccountDb) GetRealm() string {
	return "Authorization Required"
}

func (a *FakeAccountDb) SearchCredential(credential string) (account any, found bool) {
	for key, value := range a.creds {
		if subtle.ConstantTimeCompare([]byte(key), []byte(credential)) == 1 {
			return value, true
		}
	}
	return nil, false
}
```


#### CORS
(on-going)


### OpenAPI generation
When you use dij-gin style to setup server, dij-gin server will automatically
generate OpenAPI document if you need.

```go
// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	. "github.com/letscool/dij-gin"
	"github.com/letscool/dij-gin/libs"
	"log"
)

type TWebServer struct {
	WebServer

	openapi *libs.SwaggerController `di:""` // Bind OpenApi controller in root.
}

// GetResp a http request with "get" method.
// Url should like this in local: http://localhost:8000/resp?select=1 .
// Use *curl -v* command to see response code.
func (s *TWebServer) GetResp(ctx struct {
	WebContext
	Select int `http:"select"`
}) (result struct {
	Ok200 *string // the range of last three characters is between 2xx and 5xx, so the response code = 200
	Ok    *string `http:"201"` // force response code to 201
	Error error   // default response code for error is 400
}) {
	switch ctx.Select {
	case 1:
		data := "ok"
		result.Ok200 = &data
	case 2:
		data := "ok"
		result.Ok = &data
	default:
		result.Error = errors.New("an error")
	}
	return
}

// The OpenAPI page will be enabled in location: http://localhost:8000/doc.
func main() {
	config := NewWebConfig().
		SetOpenApi(func(o *OpenApiConfig) {
			o.SetEnabled(true).UseHttpOnly().SetDocPath("doc")
		})
	if err := LaunchGin(&TWebServer{}, config); err != nil {
		log.Fatalln(err)
	}
}
```


## Http Tag
______

(on-going)

#### Attributes

- path
- name
- method
- env
- tag
- middleware

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


## TODO List
_____

Still many function should be implemented, such as:
- Redirect response
- Response urlencoded data
- More middlewares
- Fix bugs
- More examples for http tag settings
- Add unit tests
- Dynamic path for controller
- NoRoute