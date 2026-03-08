package main

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	monolithURL        string
	moviesServiceURL   string
	eventsServiceURL   string
	gradualMigration   bool
	migrationPercent   int
)

func main() {
	monolithURL = getEnv("MONOLITH_URL", "http://monolith:8080")
	moviesServiceURL = getEnv("MOVIES_SERVICE_URL", "http://movies-service:8081")
	eventsServiceURL = getEnv("EVENTS_SERVICE_URL", "http://events-service:8082")

	gm := getEnv("GRADUAL_MIGRATION", "true")
	gradualMigration = gm == "true"

	mp := getEnv("MOVIES_MIGRATION_PERCENT", "50")
	var err error
	migrationPercent, err = strconv.Atoi(mp)
	if err != nil {
		migrationPercent = 50
	}

	port := getEnv("PORT", "8000")

	http.HandleFunc("/", handleProxy)

	log.Printf("Strangler Fig Proxy starting on port %s", port)
	log.Printf("Gradual migration: %v, movies migration percent: %d%%", gradualMigration, migrationPercent)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleProxy(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if path == "/health" {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Strangler Fig Proxy is healthy"))
		return
	}

	if strings.HasPrefix(path, "/api/events") {
		proxyTo(w, r, eventsServiceURL)
		return
	}

	if strings.HasPrefix(path, "/api/movies") {
		if gradualMigration {
			roll := rand.Intn(100)
			if roll < migrationPercent {
				log.Printf("[PROXY] %s %s → movies-service (roll=%d, threshold=%d%%)", r.Method, path, roll, migrationPercent)
				proxyTo(w, r, moviesServiceURL)
			} else {
				log.Printf("[PROXY] %s %s → monolith (roll=%d, threshold=%d%%)", r.Method, path, roll, migrationPercent)
				proxyTo(w, r, monolithURL)
			}
		} else {
			log.Printf("[PROXY] %s %s → monolith (gradual migration disabled)", r.Method, path)
			proxyTo(w, r, monolithURL)
		}
		return
	}

	log.Printf("[PROXY] %s %s → monolith", r.Method, path)
	proxyTo(w, r, monolithURL)
}

func proxyTo(w http.ResponseWriter, r *http.Request, targetBase string) {
	targetURL := targetBase + r.URL.RequestURI()

	proxyReq, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, "Failed to create proxy request", http.StatusInternalServerError)
		return
	}

	for key, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(key, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		log.Printf("[PROXY] Error proxying to %s: %v", targetURL, err)
		http.Error(w, "Service unavailable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
