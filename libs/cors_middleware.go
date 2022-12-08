package libs

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	. "github.com/letscool/dij-gin"
)

type CorsMiddleware struct {
	WebMiddleware

	corsHandle gin.HandlerFunc
}

func (m *CorsMiddleware) DidDependencyInitialization() {
	// cors
	// ref: https://github.com/gin-contrib/cors
	m.corsHandle = cors.Default()
}

func (m *CorsMiddleware) HandleCors(ctx struct {
	WebContext `http:""`
}) {
	m.corsHandle(ctx.Context)
}
