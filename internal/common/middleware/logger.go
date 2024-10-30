package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StructuredLog(l *logrus.Entry) gin.HandlerFunc {
	return func(c *gin.Context) {
		//t := time.Now()
		c.Next()
		//elapsed := time.Since(t)
		//l.WithFields(logrus.Fields{
		//	"time_elapsed_ms": elapsed.Milliseconds(),
		//	"request_uri":     c.Request.RequestURI,
		//	"remote_addr":     c.RemoteIP(),
		//	"client_ip":       c.ClientIP(),
		//	"full_path":       c.FullPath(),
		//}).Info("request_out")
	}
}
