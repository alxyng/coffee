package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/apex/gateway"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/nlopes/slack"
	"github.com/nullseed/devcoffee/services"
)

const (
	bucket = "coffee-storage"
	key    = "results.json"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	awsSession := session.Must(session.NewSession())

	statsService := services.NewS3StatsService(services.S3StatsOptions{
		Bucket:     bucket,
		Key:        key,
		Downloader: s3manager.NewDownloader(awsSession),
		Uploader:   s3manager.NewUploader(awsSession),
	})

	api := slack.New(os.Getenv("SLACK_TOKEN"))
	channel := os.Getenv("SLACK_CHANNEL")

	memberService := services.NewSlackMemberService(api, channel)

	handler := NewCoffeeHandler(memberService, statsService)
	http.Handle("/", handler)

	log.Println("Ready")

	log.Fatal(gateway.ListenAndServe("", nil))
}
