package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/apex/gateway"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/nlopes/slack"
	"github.com/nullseed/devcoffee/services"
	"github.com/nullseed/devcoffee/services/member"
	"github.com/nullseed/devcoffee/services/stats"
)

const (
	bucket = "coffee-storage.myunidays.com"
	key    = "results.json"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	statsService := createStatsService()
	memberService := createMemberService()

	http.Handle("/", NewCoffeeHandler(memberService, statsService))

	log.Fatal(gateway.ListenAndServe("", nil))
}

func createStatsService() services.StatsService {
	var awsSession *session.Session

	if os.Getenv("AWS_SAM_LOCAL") == "true" {
		awsSession = session.Must(session.NewSession())
	} else {
		awsSession = session.Must(session.NewSession())
	}

	s3Client := s3.New(awsSession)
	return stats.NewS3StatsService(stats.S3StatsOptions{
		Bucket:     bucket,
		Key:        key,
		Downloader: stats.NewS3Downloader(bucket, key, s3Client),
		Uploader:   stats.NewS3Uploader(bucket, key, s3Client),
	})
}

func createMemberService() services.MemberService {
	api := slack.New(os.Getenv("SLACK_TOKEN"))
	channel := os.Getenv("SLACK_CHANNEL")
	return member.NewSlackMemberService(api, channel)
}
