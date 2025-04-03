package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"structlog/clog"
	"time"

	"github.com/google/uuid"
)

var traceID int64

type tracerRoundTripper struct {
	rt http.RoundTripper
}

func (tt *tracerRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set(TraceIDHeader, getTraceID(r.Context()))

	return tt.rt.RoundTrip(r)
}

func newTrace(ctx context.Context) context.Context {
	trace := fmt.Sprintf("trace_%d", traceID)
	traceID++

	return context.WithValue(ctx, "trace_id", trace)
}

func getTraceID(ctx context.Context) string {
	return ctx.Value("trace_id").(string)
}

func SendImportantRequests(ctx context.Context) error {
	client := &http.Client{
		Transport: &tracerRoundTripper{http.DefaultTransport},
	}

	t := time.NewTicker(time.Second * 2)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			ctx = newTrace(ctx)
			ctx = clog.WithLogger(ctx, clog.Logger(ctx).With("trace_id", getTraceID(ctx)))

			err := sendImportantRequest(ctx, client)
			if err != nil {
				return fmt.Errorf("unable to send request: %w", err)
			}
		}
	}
}

var entryID int64

func sendImportantRequest(ctx context.Context, client *http.Client) error {
	reqEntryID := entryID % 5
	entryID++

	ctx = clog.WithLogger(ctx, clog.Logger(ctx).With("entry_id", reqEntryID))

	reqData := uuid.NewString()
	if rand.Int()%2 == 0 {
		reqData += uuid.NewString() + uuid.NewString()
	}

	data, err := json.Marshal(ImportantRequest{
		ID:   reqEntryID,
		Data: reqData,
	})
	if err != nil {
		return fmt.Errorf("unable to marsal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d%s", Port, APIPath), bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("unable to create new request: %w", err)
	}

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("unable to make a request: %w", err)
	}
	defer resp.Body.Close()

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read response body: %w", err)
	}

	clog.Logger(ctx).Info("response", slog.String("body", string(data)))

	return nil
}
