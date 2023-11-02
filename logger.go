package resttest

import (
	"github.com/blendle/zapdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() (*zap.Logger, error) {

	c := zapdriver.NewProductionConfig()
	c.DisableStacktrace = true
	c.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)

	return c.Build()
}
