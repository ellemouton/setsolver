package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"setsolver"
)

type Card struct {
	ID     int    `json:"id"`
	Colour string `json:"colour"`
	Fill   string `json:"fill"`
	Shape  string `json:"shape"`
	Count  string `json:"count"`
}

type Request struct {
	Cards []Card `json:"cards"`
}

type Response struct {
	Sets [][]int `json:"sets"`
}

type state struct {
	solver setsolver.SetSolver
}

func (s *state) findSets(inputCards []Card) ([][]int, error) {
	cards := unmarshalCards(inputCards)

	sets, err := s.solver.Solve(cards)
	if err != nil {
		return nil, err
	}

	return marshalSets(sets), nil
}

func marshalSets(sets []*setsolver.Set) [][]int {
	marshalledSets := make([][]int, len(sets))

	for i, s := range sets {
		marshalledSets[i] = make([]int, 3)

		marshalledSets[i][0] = s.Cards[0].ID
		marshalledSets[i][1] = s.Cards[1].ID
		marshalledSets[i][2] = s.Cards[2].ID
	}

	return marshalledSets
}

func unmarshalCards(cards []Card) []*setsolver.Card {
	marshalledCards := make([]*setsolver.Card, len(cards))

	for i, c := range cards {
		marshalledCards[i] = unmarshalCard(c)
	}

	return marshalledCards
}

func unmarshalCard(card Card) *setsolver.Card {
	return setsolver.NewDefaultCard(
		card.ID,
		setsolver.VariantName(card.Colour),
		setsolver.VariantName(card.Count),
		setsolver.VariantName(card.Shape),
		setsolver.VariantName(card.Fill),
	)
}

func (s *state) setHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(
			w, "Invalid request method", http.StatusMethodNotAllowed,
		)
	}

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sets, err := s.findSets(req.Cards)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Responding with: %v\n", sets)

	response := Response{
		Sets: sets,
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	solver, err := setsolver.NewDefaultSetSolver()
	if err != nil {
		log.Fatal(err)
	}

	state := &state{solver: solver}

	http.HandleFunc("/set", state.setHandler)
	fmt.Println("SetSolver server is running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
