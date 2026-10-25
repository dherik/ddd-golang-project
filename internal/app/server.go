package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"slices"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/dherik/ddd-golang-project/internal/app/api"
	"github.com/dherik/ddd-golang-project/internal/app/messaging"
	"github.com/dherik/ddd-golang-project/internal/app/task"
	"github.com/dherik/ddd-golang-project/internal/app/user"
	"github.com/dherik/ddd-golang-project/internal/infrastructure/clock"
	"github.com/dherik/ddd-golang-project/internal/infrastructure/messaging/rabbitmq"
	"github.com/dherik/ddd-golang-project/internal/infrastructure/persistence"
	"github.com/dherik/ddd-golang-project/internal/infrastructure/security"
)

type Server struct {
	Datasource         persistence.Datasource
	RabbitMQDataSource rabbitmq.RabbitMQDataSource
	HTTPPort           int
	JWTSecret          string

	// Listener, when set, is served instead of binding HTTPPort. Callers that
	// need to know the effective port before the server is serving — the
	// integration tests, which must not collide with a locally running app or
	// the docker-compose stack — create the listener themselves and read the
	// port back from it. When nil, HTTPPort is bound as usual.
	Listener net.Listener
}

const requestIDKey = "request_id"

var requestIDCounter atomic.Uint64

func (s *Server) Start() {
	pgsql, err := persistence.NewPostgreRepository(s.Datasource)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	taskRepository := persistence.NewTaskRepository(pgsql)
	userRepository := persistence.NewUserRepository(pgsql)
	hasher := security.NewBcryptHasher()
	timeSource := clock.NewSystemClock()

	setupSlog()

	rabbitMQ, err := rabbitmq.NewRabbitMQ(s.RabbitMQDataSource.ConnectionString())
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := messaging.Declare(rabbitMQ); err != nil {
		log.Fatalf("Failed to declare the broker topology: %v", err)
	}

	taskPublisher, err := messaging.NewPublisher(s.RabbitMQDataSource.ConnectionString())
	if err != nil {
		log.Fatalf("Failed to create the event publisher: %v", err)
	}
	defer taskPublisher.Close()

	taskService := task.NewService(taskRepository, timeSource, taskPublisher)
	userService := user.NewService(userRepository, hasher, timeSource)
	taskHandler := api.NewTaskHandler(taskService)
	loginHandler := api.NewLoginHandler(userService, s.JWTSecret)
	userHandler := api.NewUserHandler(userService)
	routes := api.NewRouter(taskHandler, loginHandler, userHandler, s.JWTSecret)

	consumer, err := messaging.NewConsumer(s.RabbitMQDataSource.ConnectionString())
	if err != nil {
		log.Fatalf("Failed to create the event consumer: %v", err)
	}

	calendarHandler := &messaging.CalendarHandler{}
	consumer.Add(messaging.QueueCalendarBirthday, calendarHandler.HandleMessage)

	invoiceDueHandler := messaging.NewInvoiceDueHandler(taskService)
	consumer.Add(messaging.QueueBillingInvoiceDue, invoiceDueHandler.HandleMessage)

	go func() {
		if err := consumer.Run(ctx); err != nil {
			slog.Error("event consumer stopped", "err", err)
			return
		}
		slog.Error("event consumer stopped unexpectedly: no handler is running")
	}()

	mux := http.NewServeMux()
	routes.SetupRoutes(mux)

	handler := chain(mux,
		recoverMiddleware,
		requestIDMiddleware,
		requestLoggerMiddleware,
	)

	httpServer := &http.Server{
		Handler: handler,
	}

	if s.Listener != nil {
		// The caller owns the listener, so it already knows the bound address.
		slog.Info("Initializing the HTTP server...", "addr", s.Listener.Addr().String())
		if err := httpServer.Serve(s.Listener); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
		return
	}

	httpServer.Addr = fmt.Sprintf(":%d", s.HTTPPort)
	slog.Info("Initializing the HTTP server...", "port", s.HTTPPort)
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}

func setupSlog() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

func chain(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	for _, mw := range slices.Backward(mws) {
		h = mw(h)
	}
	return h
}

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("panic recovered", "err", rec)
				api.CustomHTTPErrorHandler(api.NewHTTPError(http.StatusInternalServerError, "Internal Server Error"), w, r)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = strconv.FormatUint(requestIDCounter.Add(1), 36)
		}
		w.Header().Set("X-Request-ID", id)
		ctx := context.WithValue(r.Context(), requestIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func requestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		latency := time.Since(start)

		level := slog.LevelInfo
		if rec.status >= http.StatusInternalServerError {
			level = slog.LevelError
		} else if rec.status >= http.StatusBadRequest {
			level = slog.LevelWarn
		}
		slog.LogAttrs(r.Context(), level, "request",
			slog.String("method", r.Method),
			slog.String("uri", r.URL.RequestURI()),
			slog.Int("status", rec.status),
			slog.String("latency", latency.String()),
			slog.String("request_id", requestIDFromContext(r.Context())),
		)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status  int
	written bool
}

func (s *statusRecorder) WriteHeader(code int) {
	if !s.written {
		s.status = code
		s.written = true
	}
	s.ResponseWriter.WriteHeader(code)
}

func (s *statusRecorder) Write(b []byte) (int, error) {
	if !s.written {
		s.status = http.StatusOK
		s.written = true
	}
	return s.ResponseWriter.Write(b)
}

func requestIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(requestIDKey).(string); ok {
		return v
	}
	return ""
}
