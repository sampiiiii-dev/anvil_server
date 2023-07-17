package logs

import "go.uber.org/zap"

func HireScribe() *zap.Logger {
	scribe, _ := zap.NewProduction()
	defer scribe.Sync()

	return scribe
}
