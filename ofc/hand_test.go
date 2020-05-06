package ofc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidTrue(t *testing.T) {
	h := Hand{
		[]Card{"Ah"},
		[]Card{"2d", "3d"},
		[]Card{"9s", "Ts"}}

	assert.True(t, h.IsValid())
}

func TestIsValidFalse(t *testing.T) {
	h := Hand{
		[]Card{"Ah", "Ad", "8d", "4d"},
		[]Card{"2d", "3d"},
		[]Card{"9s", "Ts"}}

	assert.False(t, h.IsValid())
}

func TestCards(t *testing.T) {
	h := Hand{
		[]Card{"Ah"},
		[]Card{"2d", "3d"},
		[]Card{"9s", "Ts"}}

	assert.Equal(t, len(h.Cards()), 5)
}
