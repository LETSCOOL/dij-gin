package libs

import (
	"github.com/gin-gonic/gin"
	. "github.com/letscool/dij-gin"
)

const RefKeyForLogFormatter = "_.mdl.log.formatter"

type LogMiddleware struct {
	WebMiddleware

	f          gin.LogFormatter `di:"_.mdl.log.formatter"`
	logHandler gin.HandlerFunc
}

func (l *LogMiddleware) DidDependencyInitialization() {
	//l.logHandler = gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
	//	// your custom format
	//	return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
	//		param.ClientIP,
	//		param.TimeStamp.Format(time.RFC1123),
	//		param.Method,
	//		param.Path,
	//		param.Request.Proto,
	//		param.StatusCode,
	//		param.Latency,
	//		param.Request.UserAgent(),
	//		param.ErrorMessage,
	//	)
	//})
	l.logHandler = gin.LoggerWithFormatter(l.f)
}

func (l *LogMiddleware) HandleLog(ctx struct {
	WebContext `http:""`
}) {
	l.logHandler(ctx.Context)
}
