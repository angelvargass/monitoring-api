package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/angelvargass/monitoring-api/internal/tracing"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func errorHandler(w http.ResponseWriter, r *http.Request) { //nolint:unparam,revive
	w.WriteHeader(500)
	w.Write([]byte("error endpoint")) //nolint:errcheck
}

func readyHandler(w http.ResponseWriter, r *http.Request) { //nolint:unparam,revive
	w.Write([]byte("ready")) //nolint:errcheck
}

func healthzHandler(w http.ResponseWriter, r *http.Request) { //nolint:unparam,revive
	w.Write([]byte("ok")) //nolint:errcheck
}

func otelHandler(w http.ResponseWriter, r *http.Request) { //nolint:unparam,revive
	w.Write([]byte("Hello from OTEL traced API!\n")) //nolint:errcheck
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	shutdown := tracing.InitTracer(ctx, "monitoring-api")
	defer func() {
		if err := shutdown(ctx); err != nil {
			slog.Error("error shutting down tracer", slog.String("error", err.Error()))
		}
	}()

	mux := http.NewServeMux()
	otelHandler := otelhttp.NewHandler(http.HandlerFunc(otelHandler), "otelHandler")

	mux.Handle("/otel", otelHandler)
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/error", errorHandler)
	mux.HandleFunc("/ready", readyHandler)
	mux.HandleFunc("/healthz", healthzHandler)

	slog.Info("listening", slog.Int("port", 8080))
	err := http.ListenAndServe(":8080", mux) //nolint:errcheck,gosec
	if err != nil {
		slog.Error("server error", slog.String("error", err.Error()))
	}
}
