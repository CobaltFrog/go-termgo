package render

type Cell struct {
	Char rune 
	Depth uint16
	
	Foreground Color
	Background Color
}

type CellDiff struct {
	X int
	Y int
	Cell *Cell
}

var emptyCell = Cell {
	' ',
	0,
	ansi[0],
	ansi[0],
}

func EmptyCell() Cell {
	return emptyCell
}
