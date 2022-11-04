package main

import (
	"github.com/gin-gonic/gin"
	"github.com/letscool/lc-go/dij"
	"github.com/yuchi518/dij-gin/executor"
	"reflect"
)

type TWebServer struct {
	_ *TAuthRoute `di:"" router:"auth,"`
	_ *TUserRoute `di:"" router:"user,"`
}

func (s *TWebServer) GetRoot(p struct {
	ctx gin.Context `http:"abc"`
}) {

}

type TAuthRoute struct {
}

func (a *TAuthRoute) GetToken() {

}

type TUserRoute struct {
}

func (u *TUserRoute) GetUser() {

}

func main() {
	dij.EnableLog()
	executor.LaunchWebServer(reflect.TypeOf(TWebServer{}))

}
