package render

type GraphicsMode uint8

type Cell struct {
	Char rune 
	Depth uint16
	Style GraphicsMode
	Foreground Color
	Background Color
}

type CellDiff struct {
	X uint16
	Y uint16
	Cell *Cell
}

const GMNONE GraphicsMode = 0
const (
	GMBOLD GraphicsMode = 1 << iota
	GMDIMD
	GMITAL
	GMUNDR
	GMBLNK
	GMRVRS
	GMHIDN
	GMSTRK
)

func EmptyCell() Cell {
	return Cell {
		' ',
		0,
		GMNONE,
		ansi[0],
		ansi[0],
	}
}
func (a GraphicsMode) Has(f GraphicsMode) bool {
	return a&f == f
}

func (a GraphicsMode) HasAny(f GraphicsMode) bool {
	return a&f != 0
}

func (a *GraphicsMode) Set(f GraphicsMode) {
	*a |= f
}

func (a *GraphicsMode) Unset(f GraphicsMode) {
	*a &^= f
}

func (a *GraphicsMode) Toggle(f GraphicsMode) {
	*a ^= f
}
