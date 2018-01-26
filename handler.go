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
	member := r.FormValue("user_id")
	arg := r.FormValue("text")

	if arg == "" {
		h.handleCoffeeDraw(w)
		return
	}

	if arg == "ready" {
		h.handleCoffeeReady(member, w)
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

	writeResponse(w, fmt.Sprintf("You're up <@%v>! ☕", chosenMember))
}

func (h CoffeeHandler) handleCoffeeReady(member string, w http.ResponseWriter) {
	log.Println("Handling coffee ready")

	err := h.statsService.Increment(member)
	if err != nil {
		log.Printf("error incrementing member stats: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeResponse(w, "<!here> Coffee's ready! ☕")
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

	results, err := h.getMemberStatsWithNames(stats)
	if err != nil {
		log.Printf("error getting member stats with names: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeResponse(w, strings.Join(results, ", "))
}

type memberStats struct {
	Name        string
	CoffeesMade int
	Error       error
}

func (h CoffeeHandler) getMemberStatsWithNames(stats map[string]int) ([]string, error) {
	ch := make(chan memberStats)
	for k, v := range stats {
		go h.getMemberName(k, v, ch)
	}

	var results []string
	for _ = range stats {
		stats := <-ch

		if stats.Error != nil {
			return nil, stats.Error
		}

		results = append(results, fmt.Sprintf("%v: %v", stats.Name, stats.CoffeesMade))
	}

	return results, nil
}

func (h CoffeeHandler) getMemberName(member string, coffeesMade int, ch chan<- memberStats) {
	name, err := h.memberService.GetMemberName(member)
	ch <- memberStats{
		Name:        name,
		CoffeesMade: coffeesMade,
		Error:       err,
	}
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
