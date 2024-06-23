package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	slog.Info("Starting server")
	healthzFail := false
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		b, _ := json.Marshal(map[string]string{"message": "ok"})
		w.Write(b)
	})

	http.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if healthzFail {
			b, _ := json.Marshal(map[string]string{"status": "fail"})
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(b)
			return
		}
		b, _ := json.Marshal(map[string]string{"status": "ok"})
		w.Write(b)
	})

	http.HandleFunc("GET /healthz/toggle", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		healthzFail = !healthzFail
		b, _ := json.Marshal(map[string]bool{"healthz status": healthzFail})
		w.Write(b)
	})

	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT env var is required")
	}
	slog.Info(fmt.Sprintf("Listening on port %s", port))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		panic(err)
	}
}
