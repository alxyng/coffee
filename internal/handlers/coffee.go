package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/nullseed/coffee/services"
)

type CoffeeHandler struct {
	memserv   services.MemberService
	statsserv services.StatsService
}

func NewCoffeeHandler(m services.MemberService, s services.StatsService) CoffeeHandler {
	return CoffeeHandler{memserv: m, statsserv: s}
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

	chosenMember, err := h.memserv.GetRandomMember()
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

	err := h.statsserv.Increment(member)
	if err != nil {
		log.Printf("error incrementing member stats: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeResponse(w, "<!here> Coffee's ready! ☕")
}

func (h CoffeeHandler) handleCoffeeStats(w http.ResponseWriter) {
	log.Println("Handling coffee stats")

	stats, err := h.statsserv.Get()
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
	id          string
	name        string
	coffeesMade int
}

type byCoffeesMade []memberStats

func (a byCoffeesMade) Len() int {
	return len(a)
}

func (a byCoffeesMade) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a byCoffeesMade) Less(i, j int) bool {
	return a[i].coffeesMade < a[j].coffeesMade
}

func (h CoffeeHandler) getMemberStatsWithNames(stats map[string]int) ([]string, error) {
	members := []string{}
	for k := range stats {
		members = append(members, k)
	}

	names, err := h.memserv.GetMemberNames(members)
	if err != nil {
		return nil, err
	}

	var results []memberStats
	for k, v := range stats {
		name, _ := names[k]

		results = append(results, memberStats{
			id:          k,
			name:        name,
			coffeesMade: v,
		})
	}

	sort.Sort(sort.Reverse(byCoffeesMade(results)))

	var output []string
	for i, r := range results {
		pattern := "%v: %v"

		if i == 0 {
			pattern += " :trophy:"
		}

		if r.coffeesMade == 69 {
			pattern += " :archer:"
		}

		output = append(output, fmt.Sprintf(pattern, r.name, r.coffeesMade))
	}

	return output, nil
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
