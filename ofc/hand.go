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

func (h *Hand) NumCards() int {
	return len(h.Top) + len(h.Mid) + len(h.Bot)
}

func (h *Hand) Empty() bool {
	return len(h.Top) == 0 && len(h.Mid) == 0 && len(h.Bot) == 0
}

func (h *Hand) IsValid() bool {
	// FIXME: check total hand size (must be one of 0,5,7,9,11,13)
	return len(h.Top) <= 3 && len(h.Mid) <= 5 && len(h.Bot) <= 5
}

func (h *Hand) FindCard(card Card) (Position, bool) {
	for _, c := range h.Top {
		if c.Val == card.Val {
			return Top, true
		}
	}

	for _, c := range h.Mid {
		if c.Val == card.Val {
			return Mid, true
		}
	}

	for _, c := range h.Bot {
		if c.Val == card.Val {
			return Bot, true
		}
	}

	return Top, false
}
