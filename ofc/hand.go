package ofc

type Hand struct {
	Top []Card // max length 3
	Mid []Card // max length 5
	Bot []Card // max length 5
}

func (h *Hand) Cards() []Card {
	allCards := make([]Card, 0)
	allCards = append(allCards, h.Top...)
	allCards = append(allCards, h.Mid...)
	allCards = append(allCards, h.Bot...)
	return allCards
}

func (h *Hand) IsValid() bool {
	// FIXME: check total hand size (must be one of 0,5,7,9,11,13)
	return len(h.Top) <= 3 && len(h.Mid) <= 5 && len(h.Bot) <= 5
}
