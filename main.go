package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	slog.Info("Starting server")
	start := time.Now()
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		az, err := myAz()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			b, _ := json.Marshal(map[string]string{"message": err.Error()})
			w.Write(b)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		b, _ := json.Marshal(map[string]string{
			"message": "ok",
			"az":      az,
		})
		w.Write(b)
	})

	http.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if time.Since(start).Seconds() > 60 {
			b, _ := json.Marshal(map[string]string{"status": "fail"})
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(b)
			return
		}
		b, _ := json.Marshal(map[string]string{"status": "ok"})
		w.Write(b)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	slog.Info(fmt.Sprintf("Listening on port %s", port))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		panic(err)
	}
}

func myAz() (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodPut, "http://169.254.169.254/latest/api/token", nil)
	if err != nil {
		return "", fmt.Errorf("erro ao criar solicitação: %w", err)
	}
	req.Header.Set("X-aws-ec2-metadata-token-ttl-seconds", "21600")
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro ao obter metadados: %w", err)
	}
	defer resp.Body.Close()
	token, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erro ao ler resposta: %w", err)
	}
	req, err = http.NewRequest(http.MethodGet, "http://169.254.169.254/latest/meta-data/placement/availability-zone", nil)
	if err != nil {
		return "", fmt.Errorf("erro ao criar solicitação: %w", err)
	}
	req.Header.Set("X-aws-ec2-metadata-token", string(token))
	resp, err = client.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro ao obter metadados: %w", err)
	}
	defer resp.Body.Close()

	az, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erro ao ler resposta: %w", err)
	}
	return string(az), nil
}
