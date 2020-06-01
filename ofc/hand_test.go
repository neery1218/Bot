package ofc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func CreateCards(cardsStr []string) []Card {
	cards := make([]Card, 0)
	for _, c := range cardsStr {
		cards = append(cards, Card{c, Coord{0, 0}})
	}
	return cards
}

func TestIsValidTrue(t *testing.T) {
	h := Hand{
		CreateCards([]string{"Ah"}),
		CreateCards([]string{"2d", "3d"}),
		CreateCards([]string{"9s", "Ts"})}
	assert.True(t, h.IsValid())
}

func TestIsValidFalse(t *testing.T) {
	h := Hand{
		CreateCards([]string{"Ah", "Ad", "8d", "4d"}),
		CreateCards([]string{"2d", "3d"}),
		CreateCards([]string{"9s", "Ts"})}

	assert.False(t, h.IsValid())
}

func TestCards(t *testing.T) {
	h := Hand{
		CreateCards([]string{"Ah"}),
		CreateCards([]string{"2d", "3d"}),
		CreateCards([]string{"9s", "Ts"})}

	assert.Equal(t, len(h.Cards()), 5)
}
