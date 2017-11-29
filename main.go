package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
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

	db, err := bolt.Open("devcoffee.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	service, err := NewDiskStatsService(db)
	if err != nil {
		log.Fatal(err)
	}

	handler := NewCoffeeHandler(service)
	http.Handle("/need-coffee-please", handler)

	log.Fatal(http.ListenAndServe(":3000", nil))
}
