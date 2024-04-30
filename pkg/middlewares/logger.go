package middlewares

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		start := time.Now()
		ctx.Next()
		end := time.Now()

		latency := end.Sub(start)
		status := ctx.Writer.Status()
		clientIP := ctx.ClientIP()
		method := ctx.Request.Method

		log.Printf("\n[GIN] %v | %3d | %12v | %s | %-7s %s\n",
			end.Format("2006/01/02 - 15:04:05"),
			status,
			latency,
			clientIP,
			method,
			ctx.Request.URL.Path,
		)

	}
}
