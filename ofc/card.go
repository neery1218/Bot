package ofc

type Card string

func (c Card) IsValid() bool {
	return len(string(c)) == 2
}
