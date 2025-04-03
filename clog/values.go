package clog

import (
	"context"
	"log/slog"
	"runtime"
	"strconv"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func WithTracing(ctx context.Context) (traceID, spanID slog.Attr) {
	span, _ := tracer.SpanFromContext(ctx) //nolint // if empty it'll be a noop span; no npe
	return slog.Uint64("dd.trace_id", span.Context().TraceID()),
		slog.Uint64("dd.span_id", span.Context().SpanID())
}

// WithStacktrace usage is discouraged
// Use error wrapping + sensible error messages to log on the top of the stack
// Use Logger(ctx) + logger.With(...) + WithLogger(ctx) to accumulate context to log at the botttom of the stack
func WithStacktrace() slog.Attr {
	stack := make([]uintptr, 256)
	length := runtime.Callers(2, stack)

	type stackFrame struct {
		Func   string `json:"func"`
		Caller string `json:"caller"`
	}

	stackTrace := make([]stackFrame, 0, length)

	frames := runtime.CallersFrames(stack[:length])
	for {
		frame, more := frames.Next()
		if !more {
			break
		}

		stackTrace = append(stackTrace,
			stackFrame{
				Func:   frame.Func.Name(),
				Caller: frame.File + ":" + strconv.Itoa(frame.Line),
			},
		)
	}

	return slog.Any("stacktrace", stackTrace)
}
