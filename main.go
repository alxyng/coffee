package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/nlopes/slack"
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
	http.HandleFunc("/need-coffee-please", handleCoffeeRequest)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func handleCoffeeRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling coffee request")

	api := slack.New(os.Getenv("SLACK_TOKEN"))

	channelMembers, err := getChannelMembers(api)
	if err != nil {
		log.Printf("error getting channel members: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("Channel members: %v\n", channelMembers)

	activeMembers, err := getActiveMembers(api, channelMembers)
	if err != nil {
		log.Printf("error getting active members: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("Active members: %v\n", activeMembers)

	chosenMember := activeMembers[rand.Intn(len(activeMembers))]
	log.Printf("Chosen member: %v\n", chosenMember)

	res := response{
		ResponseType: "in_channel",
		Text:         fmt.Sprintf("You're up <@%v>!", chosenMember),
	}

	data, err := json.Marshal(res)
	if err != nil {
		log.Printf("error marshalling response: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func getChannelMembers(api *slack.Client) ([]string, error) {
	group, err := api.GetGroupInfo(os.Getenv("SLACK_CHANNEL"))
	if err != nil {
		return nil, err
	}

	return group.Members, nil
}

func getActiveMembers(api *slack.Client, channelMembers []string) ([]string, error) {
	var activeMembers []string

	ch := make(chan memberStatus)

	for _, member := range channelMembers {
		go getPresence(api, member, ch)
	}

	for range channelMembers {
		status := <-ch

		if status.Error != nil {
			return nil, status.Error
		}

		if status.Active {
			activeMembers = append(activeMembers, status.Member)
		}
	}

	return activeMembers, nil
}

func getPresence(api *slack.Client, member string, ch chan<- memberStatus) {
	presence, err := api.GetUserPresence(member)
	ch <- memberStatus{
		Member: member,
		Error:  err,
		Active: err == nil && presence.Presence == "active",
	}
}
