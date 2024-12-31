package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RequestLog(l *logrus.Entry) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestIn(c, l)
		defer requestOut(c, l)
		c.Next()
	}
}

func requestOut(c *gin.Context, l *logrus.Entry) {
	response, _ := c.Get("response")
	start, _ := c.Get("request_start")
	startTime := start.(time.Time)
	l.WithContext(c.Request.Context()).WithFields(logrus.Fields{
		"proc_time_ms": time.Since(startTime).Milliseconds(),
		"response":     response,
	}).Info("__request_out")
}

func requestIn(c *gin.Context, l *logrus.Entry) {
	c.Set("request_start", time.Now())
	body := c.Request.Body
	bodyBytes, _ := io.ReadAll(body)
	c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	var compactJson bytes.Buffer
	_ = json.Compact(&compactJson, bodyBytes)
	l.WithContext(c.Request.Context()).WithFields(logrus.Fields{
		"start": time.Now().Unix(),
		"args":  compactJson.String(),
		"from":  c.RemoteIP(),
		"uri":   c.Request.RequestURI,
	}).Info("__request_in")
}
