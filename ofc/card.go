package ofc

type Coord struct {
	X int
	Y int
}

type Card struct {
	Val   string
	Coord Coord
}

func (c Card) IsValid() bool {
	return len(string(c.Val)) == 2
}
