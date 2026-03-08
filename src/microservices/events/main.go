package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

var (
	kafkaBrokers string
	kafkaTopic   string
	writer       *kafka.Writer
)

type Event struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Timestamp string      `json:"timestamp"`
	Payload   interface{} `json:"payload"`
}

type EventResponse struct {
	Status    string `json:"status"`
	Partition int    `json:"partition"`
	Offset    int64  `json:"offset"`
	Event     Event  `json:"event"`
}

type MovieEvent struct {
	MovieID     int      `json:"movie_id"`
	Title       string   `json:"title"`
	Action      string   `json:"action"`
	UserID      int      `json:"user_id,omitempty"`
	Rating      float64  `json:"rating,omitempty"`
	Genres      []string `json:"genres,omitempty"`
	Description string   `json:"description,omitempty"`
}

type UserEvent struct {
	UserID    int    `json:"user_id"`
	Username  string `json:"username,omitempty"`
	Email     string `json:"email,omitempty"`
	Action    string `json:"action"`
	Timestamp string `json:"timestamp"`
}

type PaymentEvent struct {
	PaymentID  int     `json:"payment_id"`
	UserID     int     `json:"user_id"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
	Timestamp  string  `json:"timestamp"`
	MethodType string  `json:"method_type,omitempty"`
}

func main() {
	kafkaBrokers = getEnv("KAFKA_BROKERS", "kafka:9092")
	kafkaTopic = getEnv("KAFKA_TOPIC", "cinema-events")
	port := getEnv("PORT", "8082")

	writer = &kafka.Writer{
		Addr:         kafka.TCP(kafkaBrokers),
		Topic:        kafkaTopic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
	}
	defer writer.Close()

	go startConsumer()

	http.HandleFunc("/api/events/health", handleHealth)
	http.HandleFunc("/api/events/movie", handleMovieEvent)
	http.HandleFunc("/api/events/user", handleUserEvent)
	http.HandleFunc("/api/events/payment", handlePaymentEvent)

	log.Printf("Events service starting on port %s, Kafka: %s, topic: %s", port, kafkaBrokers, kafkaTopic)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"status": true})
}

func handleMovieEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload MovieEvent
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	eventID := fmt.Sprintf("movie-%d-%s", payload.MovieID, payload.Action)
	event := Event{
		ID:        eventID,
		Type:      "movie",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Payload:   payload,
	}

	partition, offset, err := publishEvent(event)
	if err != nil {
		log.Printf("[PRODUCER] Error publishing movie event: %v", err)
		writeError(w, "Failed to publish event", http.StatusInternalServerError)
		return
	}

	log.Printf("[PRODUCER] Movie event published: id=%s, partition=%d, offset=%d", eventID, partition, offset)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(EventResponse{
		Status:    "success",
		Partition: partition,
		Offset:    offset,
		Event:     event,
	})
}

func handleUserEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload UserEvent
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	eventID := fmt.Sprintf("user-%d-%s", payload.UserID, payload.Action)
	event := Event{
		ID:        eventID,
		Type:      "user",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Payload:   payload,
	}

	partition, offset, err := publishEvent(event)
	if err != nil {
		log.Printf("[PRODUCER] Error publishing user event: %v", err)
		writeError(w, "Failed to publish event", http.StatusInternalServerError)
		return
	}

	log.Printf("[PRODUCER] User event published: id=%s, partition=%d, offset=%d", eventID, partition, offset)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(EventResponse{
		Status:    "success",
		Partition: partition,
		Offset:    offset,
		Event:     event,
	})
}

func handlePaymentEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload PaymentEvent
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	eventID := fmt.Sprintf("payment-%d-%s", payload.PaymentID, payload.Status)
	event := Event{
		ID:        eventID,
		Type:      "payment",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Payload:   payload,
	}

	partition, offset, err := publishEvent(event)
	if err != nil {
		log.Printf("[PRODUCER] Error publishing payment event: %v", err)
		writeError(w, "Failed to publish event", http.StatusInternalServerError)
		return
	}

	log.Printf("[PRODUCER] Payment event published: id=%s, partition=%d, offset=%d", eventID, partition, offset)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(EventResponse{
		Status:    "success",
		Partition: partition,
		Offset:    offset,
		Event:     event,
	})
}

func publishEvent(event Event) (int, int64, error) {
	data, err := json.Marshal(event)
	if err != nil {
		return 0, 0, err
	}

	msg := kafka.Message{
		Key:   []byte(event.ID),
		Value: data,
	}

	err = writer.WriteMessages(context.Background(), msg)
	if err != nil {
		return 0, 0, err
	}

	return 0, 0, nil
}

func startConsumer() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{kafkaBrokers},
		Topic:       kafkaTopic,
		GroupID:     "events-consumer-group",
		StartOffset: kafka.LastOffset,
		MaxWait:     1 * time.Second,
	})
	defer reader.Close()

	log.Printf("[CONSUMER] Started listening on topic: %s", kafkaTopic)

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("[CONSUMER] Error reading message: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		var event Event
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("[CONSUMER] Error unmarshaling message: %v", err)
			continue
		}

		log.Printf("[CONSUMER] Received event: id=%s, type=%s, partition=%d, offset=%d, payload=%v",
			event.ID, event.Type, msg.Partition, msg.Offset, event.Payload)
	}
}

func writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
