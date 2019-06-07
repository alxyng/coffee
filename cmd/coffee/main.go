package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/nullseed/coffee/internal/config"
	"github.com/nullseed/coffee/internal/handlers"
	"github.com/nullseed/coffee/services"
	"github.com/nullseed/coffee/services/member"
	"github.com/nullseed/coffee/services/stats"

	"github.com/apex/gateway"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/nlopes/slack"
)

const (
	bucket = "coffee-storage.myunidays.com"
	key    = "results.json"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	http.Handle("/", createCoffeeHandler())
	log.Fatal(gateway.ListenAndServe("", nil))
}

func createCoffeeHandler() http.Handler {
	memberService := createMemberService()
	statsService := createStatsService()
	return handlers.NewCoffeeHandler(memberService, statsService)
}

func createStatsService() services.StatsService {
	awsSession := config.CreateAWSSession()
	db := dynamodb.New(awsSession)
	return stats.NewDynamoDBStatsService(db)
}

func createMemberService() services.MemberService {
	api := slack.New(os.Getenv("SLACK_TOKEN"))
	channel := os.Getenv("SLACK_CHANNEL")
	return member.NewSlackMemberService(api, channel)
}
