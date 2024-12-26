package logging

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	Method   = "method"
	Args     = "args"
	Cost     = "cost_ms"
	Response = "response"
	Error    = "err"
)

func WhenMySQL(ctx context.Context, method string, args ...any) (logrus.Fields, func(any, *error)) {
	fields := logrus.Fields{
		Method: method,
		Args:   formatMySQLArgs(args),
	}
	start := time.Now()
	return fields, func(resp any, err *error) {
		level, msg := logrus.InfoLevel, "mysql_success"
		fields[Cost] = time.Since(start).Milliseconds()
		fields[Response] = resp

		if err != nil && (*err != nil) {
			level, msg = logrus.ErrorLevel, "mysql_error"
			fields[Error] = (*err).Error()
		}

		logrus.WithContext(ctx).WithFields(fields).Logf(level, "%s", msg)
	}
}

func formatMySQLArgs(args []any) string {
	var item []string
	for _, arg := range args {
		item = append(item, formatMySQLArg(arg))
	}
	return strings.Join(item, "||")
}

func formatMySQLArg(arg any) string {
	switch v := arg.(type) {
	default:
		bytes, err := json.Marshal(v)
		if err != nil {
			return "unsupported type in formatMySQLArg||err=" + err.Error()
		}
		return string(bytes)
	}
}
