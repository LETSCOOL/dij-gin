package dij_gin

import "github.com/gin-gonic/gin"

type WebRouter interface {
	gin.IRouter
	BasePath() string
}

type WebRoutes interface {
	gin.IRoutes
	BasePath() string
}
