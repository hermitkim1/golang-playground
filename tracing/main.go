package main

import (
	"context"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
)

// Define context keys
type contextKey string

const (
	traceIDKey contextKey = "trace_id"
	spanIDKey  contextKey = "span_id"
)

// Define TracingHook struct
type TracingHook struct{}

// Run method: Executed when a log event occurs
func (h TracingHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	ctx := e.GetCtx()
	traceID := ctx.Value(traceIDKey)
	spanID := ctx.Value(spanIDKey)

	if traceID != nil {
		e.Str("trace_id", traceID.(string))
	}
	if spanID != nil {
		e.Str("span_id", spanID.(string))
	}
}

// Define Tracing middleware
func TracingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Start a new tracing context
		tracer := otel.Tracer("example-tracer")
		ctx, span := tracer.Start(c.Request().Context(), c.Path())
		defer span.End()

		// Store trace and span IDs in the context
		traceID := span.SpanContext().TraceID().String()
		spanID := span.SpanContext().SpanID().String()

		ctx = context.WithValue(ctx, traceIDKey, traceID)
		ctx = context.WithValue(ctx, spanIDKey, spanID)

		// Create a logger with trace_id and span_id and store it in the context
		logger := log.With().Str("trace_id", traceID).Str("span_id", spanID).Logger()
		ctx = logger.WithContext(ctx)

		// Set the context in the request
		c.SetRequest(c.Request().WithContext(ctx))

		// Call the next handler
		return next(c)
	}
}

func init() {
	// Set up the base logger (used globally)
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()

	// Add TracingHook to the logger
	log.Logger = log.Logger.Hook(TracingHook{})
}

func main() {

	// Create Echo instance
	e := echo.New()

	// Register Tracing middleware to apply to all handlers
	e.Use(TracingMiddleware)

	// Register example handler
	e.GET("/", func(c echo.Context) error {
		// Log with context
		ctx := c.Request().Context()

		log.Ctx(ctx).Info().Msg("Request received")

		// Call an internal function within the handler
		someInternalFunction(ctx)

		// Simple response
		return c.String(http.StatusOK, "Hello, World!")
	})

	// Start the server
	e.Start(":8080")
}

// Example function within a handler
func someInternalFunction(ctx context.Context) {
	// Log with context in this function as well
	log.Ctx(ctx).Info().Msg("Inside someInternalFunction")
}
