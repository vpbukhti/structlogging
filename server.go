package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"structlog/clog"

	"github.com/go-chi/chi/v5"
)

const APIPath = "/api"
const Port = 8081

func SetupServer(ctx context.Context) *http.Server {
	router := chi.NewRouter()

	router.Use(RequestIDMiddleware)
	router.Use(TraceIDServerMiddleware)
	router.Use(AuthenticationMiddleware)

	router.Post(APIPath, ImportantHandler)

	server := &http.Server{
		Addr:        fmt.Sprintf(":%d", Port),
		Handler:     router,
		BaseContext: func(l net.Listener) context.Context { return ctx },
	}

	return server
}

const UserIDHeader = "X-USER-ID"

func AuthenticationMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get(UserIDHeader)

		ctx := r.Context()
		logger := clog.Logger(ctx)
		logger = logger.With("user_id", userID)
		ctx = clog.WithLogger(ctx, logger)
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	})
}

const TraceIDHeader = "X-TRACE-ID"

func TraceIDServerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get(TraceIDHeader)

		ctx := r.Context()
		logger := clog.Logger(ctx)
		logger = logger.With("trace_id", traceID)
		ctx = clog.WithLogger(ctx, logger)
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	})
}

var requestID int64

func RequestIDMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := clog.Logger(ctx)
		logger = logger.With("request_id", fmt.Sprintf("request_%d", requestID))
		requestID++
		ctx = clog.WithLogger(ctx, logger)
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	})
}

type ImportantRequest struct {
	ID   int64
	Data string
}

func ImportantHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := ImportantRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		clog.Logger(ctx).Error("")
		sendError(w, err)
		return
	}

	ctx = clog.WithLogger(ctx, clog.Logger(ctx).With("entry_id", req.ID))

	err = ImportantDatabaseCall(ctx, req.Data)
	if err != nil {
		clog.Logger(ctx).Error("unable to make an important db call", slog.String("error", err.Error()))
		sendError(w, err)
		return
	}

	clog.Logger(ctx).Info("success", slog.Any("req", req))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func sendError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}
