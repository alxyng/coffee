package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
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

	writeResponse(w, strings.Join(results, "\n"))
}

type memberStats struct {
	Name        string
	CoffeesMade int
	Error       error
}

type ByCoffeesMade []memberStats

func (a ByCoffeesMade) Len() int {
	return len(a)
}

func (a ByCoffeesMade) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByCoffeesMade) Less(i, j int) bool {
	return a[i].CoffeesMade < a[j].CoffeesMade
}

func (h CoffeeHandler) getMemberStatsWithNames(stats map[string]int) ([]string, error) {
	ch := make(chan memberStats)
	for k, v := range stats {
		go h.getMemberName(k, v, ch)
	}

	var results []memberStats
	for _ = range stats {
		s := <-ch

		if s.Error != nil {
			return nil, s.Error
		}

		results = append(results, s)
	}

	sort.Sort(sort.Reverse(ByCoffeesMade(results)))

	var output []string
	for _, r := range results {
		output = append(output, fmt.Sprintf("%v: %v", r.Name, r.CoffeesMade))
	}

	return output, nil
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
