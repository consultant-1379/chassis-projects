package main

import (
	"net/http"
	"time"
)

func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func runProbe(port string) {
	server := &http.Server{Addr: ":" + port, ReadTimeout: 1 * time.Second, WriteTimeout: 1 * time.Second}
	http.HandleFunc("/healthz", HealthzHandler)
	logger.Info("readiness probe started on port: " + port)
	logger.Critical(server.ListenAndServe())
}
