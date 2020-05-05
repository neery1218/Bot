package automaton

import (
	"image"
	"image/color"
)

type Bitmap struct {
	ColorArray    [][]color.Color
	ColorModelVal color.Model
}

func (bp *Bitmap) ColorModel() color.Model {
	return bp.ColorModelVal
}

func (bp *Bitmap) Bounds() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: len(bp.ColorArray[0]), Y: len(bp.ColorArray)}}
}

func (bp *Bitmap) At(x, y int) color.Color {
	return bp.ColorArray[y][x]
}
