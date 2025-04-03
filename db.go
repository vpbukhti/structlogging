package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"structlog/clog"
)

const maxDataLength = 40

func ImportantDatabaseCall(ctx context.Context, data string) (err error) {
	defer func() {
		val := recover()
		if val != nil {
			clog.Logger(ctx).Error("panic", slog.Any("panic", val), clog.WithStacktrace())
			err = fmt.Errorf("panic: %v", val)
		}
	}()

	if len(data) > maxDataLength {
		clog.Logger(ctx).Error("data is too long",
			slog.Int("max_length", maxDataLength),
			slog.Int("actual_length", len(data)),
		)
		return fmt.Errorf("data is too long: max length: %d, actual length: %d", maxDataLength, len(data))
	}

	if rand.Int32()%4 == 0 {
		panic("aaaaaaaaa!")
	}

	return nil
}
