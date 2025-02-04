package utils

import (
	"go.uber.org/zap"
)

func NewLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}

func String(key, val string) zap.Field {
	return zap.String(key, val)
}
