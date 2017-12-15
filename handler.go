package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/nullseed/devcoffee/services"
)

type CoffeeHandler struct {
	memberService services.MemberService
	statsService  services.StatsService
}

func NewCoffeeHandler(m services.MemberService, s services.StatsService) CoffeeHandler {
	return CoffeeHandler{
		memberService: m,
		statsService:  s,
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

	chosenMember, err := h.memberService.GetRandomMember()
	if err != nil {
		log.Printf("error getting random member: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if chosenMember == "" {
		writeResponse(w, "No one is around to make coffee ☹️")
	}

	err = h.statsService.Increment(chosenMember)
	if err != nil {
		log.Printf("error incrementing member stats: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeResponse(w, fmt.Sprintf("You're up <@%v>! ☕", chosenMember))
}

func (h CoffeeHandler) handleCoffeeStats(w http.ResponseWriter) {
	log.Println("Handling coffee stats")

	stats, err := h.statsService.Get()
	if err != nil {
		log.Printf("error getting stats: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(stats) == 0 {
		writeResponse(w, "No one has made coffee yet! ☕")
		return
	}

	var results []string
	for k, v := range stats {
		results = append(results, fmt.Sprintf("<@%v>: %v", k, v))
	}

	writeResponse(w, strings.Join(results, ", "))
}

func handleUnknownArgument(w http.ResponseWriter, arg string) {
	log.Printf("Unknown argument: %v\n", arg)
	writeResponse(w, "Unknown argument ☹️")
}

type response struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

func writeResponse(w http.ResponseWriter, text string) {
	res := response{
		ResponseType: "in_channel",
		Text:         text,
	}

	data, err := json.Marshal(res)
	if err != nil {
		log.Printf("error marshalling response: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}
