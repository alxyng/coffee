package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"
)

type memberStatus struct {
	Member string
	Error  error
	Active bool
}

type response struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

func main() {
	rand.Seed(time.Now().Unix())

	statsService := NewMemoryStatsService()
	handler := NewCoffeeHandler(statsService)
	http.Handle("/need-coffee-please", handler)

	log.Fatal(http.ListenAndServe(":3000", nil))
}
