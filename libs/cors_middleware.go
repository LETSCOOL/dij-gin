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

func (m *CorsMiddleware) HandleCors(ctx struct {
	WebContext `http:""`
}) {
	// cors
	// ref: https://github.com/gin-contrib/cors
	// TODO: refine middleware architecture?
	//router.Use(cors.Default())
	if m.corsHandle == nil {
		m.corsHandle = cors.Default()
	}
	m.corsHandle(ctx.Context)
}
