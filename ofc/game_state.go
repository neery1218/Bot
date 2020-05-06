package ofc

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type OfcError struct {
	Msg string
}

func (e *OfcError) Error() string {
	return fmt.Sprintf("OfcError: %s", e.Msg)
}

type Coord struct {
	X int
	Y int
}

type GameState struct {
	MyHand     Hand
	OtherHands []Hand         // must be length 1 or 2
	Pull       map[Card]Coord // 0 <= len(Pull) <= 3
	DeadCards  []Card
}

func (gs *GameState) AllCards() []Card {
	allCards := make([]Card, 0)

	allCards = append(allCards, gs.MyHand.Cards()...)

	for _, otherHand := range gs.OtherHands {
		allCards = append(allCards, otherHand.Cards()...)
	}

	for card := range gs.Pull {
		allCards = append(allCards, card)

	}

	allCards = append(allCards, gs.DeadCards...)
	return allCards
}

func (gameState *GameState) IsValid() (bool, error) {
	// MyHand
	if !gameState.MyHand.IsValid() {
		return false, &OfcError{fmt.Sprintf("MyHand is invalid! %+v", gameState.MyHand)}
	}

	// OtherHands
	if !(1 <= len(gameState.OtherHands) && len(gameState.OtherHands) <= 2) {
		return false,
			&OfcError{fmt.Sprintf("OtherHands is incorrect size! %v", len(gameState.OtherHands))}
	}

	for i, otherHand := range gameState.OtherHands {
		if !otherHand.IsValid() {
			return false,
				&OfcError{fmt.Sprintf("OtherHands %v is invalid! %+v", i, otherHand)}

		}
	}

	// Pull
	if len(gameState.Pull) != 0 && len(gameState.Pull) != 3 {
		return false,
			&OfcError{"Pull length must be 0 or 3!"}

	}

	for _, coord := range gameState.Pull {
		if coord.X < 0 || coord.Y < 0 {
			return false,
				&OfcError{"Coordinates have to be >= 0!"}
		}
	}

	// check all cards for uniqueness, validity
	uniqueCards := make(map[Card]bool)
	for _, c := range gameState.AllCards() {
		if !c.IsValid() {
			return false, &OfcError{fmt.Sprintf("invalid card %v!", c)}
		}

		if _, exists := uniqueCards[c]; exists {
			return false, &OfcError{fmt.Sprintf("Duplicate card %v!", c)}
		}
		uniqueCards[c] = true
	}

	return true, nil
}

func parseGameStateFromJson(str string) (*GameState, error) {
	gameState := GameState{}
	if err := json.Unmarshal([]byte(str), &gameState); err != nil {
		return nil, err
	}

	if valid, err := gameState.IsValid(); !valid {
		return nil, err
	}

	return &gameState, nil
}

func StateChanged(gsNew, gsOld *GameState) bool {
	return !reflect.DeepEqual(gsNew, gsOld) // FIXME: apparently this is bad
}

func (gs *GameState) DecisionRequired() bool {
	return len(gs.Pull) > 0
}
