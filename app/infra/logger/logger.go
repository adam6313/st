package logger

import (
	"strings"

	"storage/app/infra/config"

	"github.com/axolotlteam/thunder/logger"
	"github.com/sirupsen/logrus"
)

// NewLogger -
func NewLogger(c config.Config) *logrus.Entry {
	logger.NewLogrus()

	var l logger.Level

	switch strings.ToLower(c.Log.Level) {
	case "trace":
		l = logger.TraceLevel
	case "debug":
		l = logger.DebugLevel
	case "info":
		l = logger.InfoLevel
	case "error":
		l = logger.ErrorLevel
	case "warn":
		l = logger.WarnLevel
	default:
		l = logger.DebugLevel
	}
	logger.SetLevel(l)
	return logger.WithField("server", c.Info.Name)
}
