package logging

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

func WhenCommandExecute(ctx context.Context, commandName string, cmd any, err error) {
	fields := logrus.Fields{
		"cmd": cmd,
	}
	if err == nil {
		logf(ctx, logrus.InfoLevel, fields, "%s_command_success", commandName)
	} else {
		logf(ctx, logrus.ErrorLevel, fields, "%s_command_failed", commandName)
	}
}

func WhenRequest(ctx context.Context, method string, args ...any) (logrus.Fields, func(any, *error)) {
	fields := logrus.Fields{
		Method: method,
		Args:   formatArgs(args),
	}
	start := time.Now()
	return fields, func(resp any, err *error) {
		level, msg := logrus.InfoLevel, "_request_success"
		fields[Cost] = time.Since(start).Milliseconds()
		fields[Response] = resp

		if err != nil && (*err != nil) {
			level, msg = logrus.ErrorLevel, "_request_failed"
			fields[Error] = (*err).Error()
		}

		logf(ctx, level, fields, "%s", msg)
	}
}

func WhenEventPublish(ctx context.Context, args ...any) (logrus.Fields, func(any, *error)) {
	fields := logrus.Fields{
		Args: formatArgs(args),
	}
	start := time.Now()
	return fields, func(resp any, err *error) {
		level, msg := logrus.InfoLevel, "_mq_publish_success"
		fields[Cost] = time.Since(start).Milliseconds()
		fields[Response] = resp

		if err != nil && (*err != nil) {
			level, msg = logrus.ErrorLevel, "_mq_publish_failed"
			fields[Error] = (*err).Error()
		}

		logf(ctx, level, fields, "%s", msg)
	}
}
