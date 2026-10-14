package render

import (
	"math"
)

type Color struct {
	Red   uint8
	Green uint8
	Blue  uint8
}
var ansi = [16]Color {
    {0x00, 0x00, 0x00}, // 00 Black
	{0xAA, 0x00, 0x00}, // 01 Red
	{0x00, 0xAA, 0x00}, // 02 Green
	{0xAA, 0x55, 0x00}, // 03 Yellow
	{0x00, 0x00, 0xAA}, // 04 Blue
	{0xAA, 0x00, 0xAA}, // 05 Magenta
	{0x00, 0xAA, 0xAA}, // 06 Cyan
	{0xAA, 0xAA, 0xAA}, // 07 White
	{0x55, 0x55, 0x55}, // 08 Bright Black
	{0xFF, 0x55, 0x55}, // 09 Bright Red
	{0x55, 0xFF, 0x55}, // 10 Bright Green
	{0xFF, 0xFF, 0x55}, // 11 Bright Yellow
	{0x55, 0x55, 0xFF}, // 12 Bright Blue
	{0xFF, 0x55, 0xFF}, // 13 Bright Magenta
	{0x55, 0xFF, 0xFF}, // 14 Bright Cyan
	{0xFF, 0xFF, 0xFF}, // 15 Bright White
}

var xterm = genXterm()

func ANSI() [16]Color {
	return ansi
}

func Xterm() [256]Color {
	return xterm
}

// Get closest to taget color index from specified palette
func nearest(trg *Color, pl []Color) uint8 {
	l := uint8(len(pl))

	var bstIdx, bstDist uint8
	bstIdx = 0
	bstDist = 255

	var i uint8 = 0
	for ; i < l; i++ {
		dr := pl[i].Red   - trg.Red
		dg := pl[i].Green - trg.Green
		db := pl[i].Blue  - trg.Blue

		dist := dr * dr + dg * dg + db * db
		if (dist < bstDist) {
			bstDist = dist
			bstIdx = i

			if bstDist == 0 {
				break
			}
		}
	}

	return bstIdx
}

func genXterm() [256]Color {
	var plt [256]Color
	copy(plt[:16], ansi[:])
	
    lvl := [6]uint8{0, 95, 135, 175, 215, 255}
	i := 16
	for r := range 6 {
		for g := range 6 {
			for b := range 6 {
				plt[i] = Color{lvl[r], lvl[g], lvl[b]}
				i++
			}
		}
	}

	for j := range 24 {
		v := uint8(8 + j * 10)
		plt[232+j] = Color{v, v, v}
	}

	return plt
}

func (clr *Color) ToAnsiColor() uint8 {
	return nearest(clr, ansi[:])
}

func (clr *Color) ToXtermColor() uint8 {
	return nearest(clr, xterm[:])
}

func FromHsv(h, s, v float32) Color {
	toByte := func(val float32) uint8 { // float value to byte
		return uint8(math.Round(float64(val * 255)))
	}

	var hue float32 = h / 60.0
	var i int = int(math.Floor(float64(hue)))
	f := hue - float32(i);
    p := v * (1.0 - s);
    q := v * (1.0 - s * f);
    t := v * (1.0 - s * (1.0 - f));

    var res Color
	switch(i) {
	case 0:
		res = Color{toByte(v), toByte(t), toByte(p)}
	case 1:
		res = Color{toByte(q), toByte(v), toByte(p)}
	case 2:
		res = Color{toByte(p), toByte(v), toByte(t)}
	case 3:
		res = Color{toByte(p), toByte(q), toByte(v)}
	case 4:
		res = Color{toByte(t), toByte(p), toByte(v)}
	default:
		res = Color{toByte(v), toByte(p), toByte(q)}
	}

	return res
}
