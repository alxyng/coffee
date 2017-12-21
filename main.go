package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
	"github.com/nlopes/slack"
	"github.com/nullseed/devcoffee/services"
)

const (
	dataDirectory = "data"
	databaseName  = "devcoffee.db"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	_, err := os.Stat(dataDirectory)
	if os.IsNotExist(err) {
		log.Println("Creating data directory")
		err = os.Mkdir(dataDirectory, os.ModePerm)
		if err != nil {
			log.Fatalf("error creating data directory: %v", err)
		}
	}

	db, err := bolt.Open(filepath.Join(dataDirectory, databaseName), 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	statsService, err := services.NewDiskStatsService(db)
	if err != nil {
		log.Fatal(err)
	}

	api := slack.New(os.Getenv("SLACK_TOKEN"))
	channel := os.Getenv("SLACK_CHANNEL")

	memberService := services.NewSlackMemberService(api, channel)

	handler := NewCoffeeHandler(memberService, statsService)
	http.Handle("/need-coffee-please", handler)

	log.Println("Ready")

	log.Fatal(http.ListenAndServe(":3000", nil))
}
