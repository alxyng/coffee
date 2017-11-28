package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/nlopes/slack"
)

type CoffeeHandler struct {
	statsService StatsService
}

func NewCoffeeHandler(statsService StatsService) CoffeeHandler {
	return CoffeeHandler{
		statsService: statsService,
	}
}

func (h CoffeeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	arg := r.FormValue("text")

	if arg == "" {
		h.handleCoffeeDraw(w)
		return
	}

	if arg == "stats" {
		h.handleCoffeeStats(w)
		return
	}

	handleUnknownArgument(w, arg)
}

func (h CoffeeHandler) handleCoffeeDraw(w http.ResponseWriter) {
	log.Println("Handling coffee draw")

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

	if len(activeMembers) == 0 {
		log.Println("No active members")
		writeResponse(w, "No one is around to make coffee ☹️")
		return
	}

	chosenMember := activeMembers[rand.Intn(len(activeMembers))]
	log.Printf("Chosen member: %v\n", chosenMember)
	h.statsService.Increment(chosenMember)
	writeResponse(w, fmt.Sprintf("You're up <@%v>! ☕", chosenMember))
}

func (h CoffeeHandler) handleCoffeeStats(w http.ResponseWriter) {
	log.Println("Handling coffee stats")

	stats := h.statsService.Get()
	log.Printf("Stats: %v\n", stats)

	if len(stats) == 0 {
		writeResponse(w, "No one has made coffee yet! ☕")
		return
	}

	var results []string
	for k, v := range stats {
		results = append(results, fmt.Sprintf("<%v>: %v", k, v))
	}

	text := strings.Join(results, ", ")
	writeResponse(w, text)
}

func handleUnknownArgument(w http.ResponseWriter, arg string) {
	log.Printf("Unknown argument: %v\n", arg)
	writeResponse(w, "Unknown argument ☹️")
}

func writeResponse(w http.ResponseWriter, text string) {
	res := response{
		ResponseType: "in_channel",
		Text:         text,
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
