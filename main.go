package main

import (
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
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

func main() {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/error", errorHandler)
	mux.HandleFunc("/ready", readyHandler)
	mux.HandleFunc("/healthz", healthzHandler)

	slog.Info("listening", slog.Int("port", 8080))
	http.ListenAndServe(":8080", mux) //nolint:errcheck,gosec
}
