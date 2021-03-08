package slog_test

import (
	"io/ioutil"
	"testing"

	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

// go test -v -cpu=4 -run=none -bench=. -benchtime=10s -benchmem bench_test.go
// code refer:
// 	https://github.com/phuslu/log
var msg = "The quick brown fox jumps over the lazy dog"

func BenchmarkGookitSlogNegative(b *testing.B) {
	logger := slog.NewWithHandlers(
		handler.NewIOWriter(ioutil.Discard, []slog.Level{slog.ErrorLevel}),
	)
	for i := 0; i < b.N; i++ {
		logger.Info("rate", "15", "low", 16, "high", 123.2, msg)
	}
}

func BenchmarkGookitSlogPositive(b *testing.B) {
	logger := slog.NewWithHandlers(
		handler.NewIOWriter(ioutil.Discard, slog.NormalLevels),
	)
	for i := 0; i < b.N; i++ {
		logger.Info("rate", "15", "low", 16, "high", 123.2, msg)
	}
}
