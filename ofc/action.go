package ofc

type Position int

const (
	Top Position = iota
	Mid
	Bot
)

func (p Position) String() string {
	return [...]string{"Top", "Mid", "Bot"}[p]
}

type Action struct {
	Card Card
	Pos  Position
}
