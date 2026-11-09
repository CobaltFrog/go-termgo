package render

import (
	"fmt"
	"io"
	"strings"
)

type Renderer interface {
	Draw(io.Writer, []CellDiff)
}

type ColorMode byte

const (
	CMTCLR ColorMode = iota
	CMXTRM 
	CMANSI
)

type modeSeq struct {
	mode GraphicsMode
	set byte
	reset byte
}
var flagTab = [8]modeSeq {
	{GMBOLD, 1, 22},
	{GMDIMD, 2, 22},
	{GMITAL, 3, 23},
	{GMUNDR, 4, 24},
	{GMBLNK, 5, 25},
	{GMRVRS, 7, 27},
	{GMHIDN, 8, 28},
	{GMSTRK, 9, 29},
}

type TerminalRenderer struct {
	Mode     ColorMode
	Mask    *Color // mask for default terminal bg color
	Repositioned bool // cursor position unknown
	builder  strings.Builder
	cursorX, cursorY uint16
}

func (r *TerminalRenderer) Draw(out io.Writer, diffs []CellDiff) {
	// because the control sequence for changing the color of terminal output takes many characters,
	// it's necessary to sort changes by line, and then group them by color to output one whole sequence
	var fg, bg Color
	var oldgm GraphicsMode
	var first = true
	for _, diff := range diffs {
		gmch := first || oldgm != diff.Cell.Style
		fgch := first || !colorsEqual(&fg, &diff.Cell.Foreground)
		bgch := first || !colorsEqual(&bg, &diff.Cell.Background)

		if r.Repositioned || r.cursorX != diff.X || r.cursorY != diff.Y {
			fmt.Fprintf(&r.builder, "\x1b[%d;%dH", diff.Y + 1, diff.X + 1)
			r.cursorX, r.cursorY = diff.X, diff.Y
			r.Repositioned = false
		}

		// if next cell has the same mode, move to next
		if !(gmch || fgch || bgch) { 
			r.builder.WriteRune(diff.Cell.Char)
			r.cursorX++
			continue
		}

		var ctlsec []string
		writeMode := func(s string) {
			ctlsec = append(ctlsec, s)
		}

		if first {
			// reset modes for terminal
			// necessary for the output in terminal to be equal to that in the diff
			writeMode("0")
			first = false
		}

		// start mode sequence
		if gmch {
			for _, gmf := range flagTab {
				oldHas := oldgm.Has(gmf.mode)
				newHas := diff.Cell.Style.Has(gmf.mode)
				if oldHas && !newHas { // remove previous modes
					writeMode(fmt.Sprint(gmf.reset))
					// return dimmed mode if it was reset by GMBOLD
					if gmf.reset == 22 && diff.Cell.Style.Has(GMDIMD) {
						writeMode("2")
					}
					// same fore bold mode (if it waw reset by GMDIMD)
					if gmf.reset == 22 && diff.Cell.Style.Has(GMBOLD) {
						writeMode("1")
					}
				} else if !oldHas && newHas { // add new if they are needed
					writeMode(fmt.Sprint(gmf.set))
				}
			}
		}
		if fgch {
			c := diff.Cell.Foreground
			switch r.Mode {
			case CMTCLR:
				writeMode(fmt.Sprintf("38;2;%d;%d;%d", c.Red, c.Green, c.Blue))
			case CMXTRM:
				writeMode(fmt.Sprintf("38;5;%d", c.ToXtermColor()))
			case CMANSI:
				writeMode(fmt.Sprint(30 + c.ToAnsiColor()))
			}

		}
		if bgch {
			c := diff.Cell.Background
			if r.Mask == nil || !colorsEqual(r.Mask, &c) {
				switch r.Mode {
				case CMTCLR:
					writeMode(fmt.Sprintf("48;2;%d;%d;%d", c.Red, c.Green, c.Blue))
				case CMXTRM:
					writeMode(fmt.Sprintf("48;5;%d", c.ToXtermColor()))
				case CMANSI:
					writeMode(fmt.Sprint(40 + c.ToAnsiColor()))
				}
			} else {
				writeMode("49")
			}
		}
		// end mode sequence, append cell's rune
		if len(ctlsec) > 0 {
			r.builder.WriteString("\x1b[")
			r.builder.WriteString(strings.Join(ctlsec, ";"))
			r.builder.WriteRune('m')
		}
		r.builder.WriteRune(diff.Cell.Char)
		r.cursorX++

		fg, bg, oldgm = diff.Cell.Foreground, diff.Cell.Background, diff.Cell.Style
	}

	_, _ = io.WriteString(out, r.builder.String())
	r.builder.Reset()
	r.Repositioned = true

	if f, ok := out.(interface{ Flush() error }); ok {
		_ = f.Flush()
	}
}

func colorsEqual(a, b *Color) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}
