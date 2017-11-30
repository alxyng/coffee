package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/nlopes/slack"
	"github.com/nullseed/devcoffee/services"
)

func main() {
	rand.Seed(time.Now().Unix())

	db, err := bolt.Open("devcoffee.db", 0600, nil)
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

	log.Fatal(http.ListenAndServe(":3000", nil))
}
