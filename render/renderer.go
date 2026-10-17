package render

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

type Renderer interface {
	Draw(*io.Writer, []CellDiff)
}

type TerminalRenderer struct {
	Mode     ColorMode
	Mask    *Color
	Escape   string
	builder  strings.Builder
}

type ColorMode int

const (
	CMTruecolor ColorMode = iota
	CMXterm 
	CMAnsi
)

func (r *TerminalRenderer) Draw(output io.Writer, updates []CellDiff) {
	// because the control sequence for changing the color of terminal output takes many characters,
	// it's necessary to sort changes by line, and then group them by color to output one whole sequence
	sort.SliceStable(updates, func(i, j int) bool {
		if updates[i].Y != updates[j].Y {
			return updates[i].Y < updates[j].Y
		}

		return updates[i].X < updates[j].X
	})

	var fg, bg Color
	var fgSet, bgSet bool
	for _, diff := range updates {
		if diff.Cell == nil {
			continue
		}

		fgChange := !fgSet || fg != diff.Cell.Foreground
		bgChange := !bgSet || bg != diff.Cell.Background

		if fgChange || bgChange {
			fg, fgSet = diff.Cell.Foreground, true
			bg, bgSet = diff.Cell.Background, true

			fmt.Fprintf(&r.builder, "%s[%d;%dH", r.Escape, diff.Y+1, diff.X+1)
			fmt.Fprintf(&r.builder, "%s[0", r.Escape)

			switch r.Mode {
			case CMTruecolor:
				fmt.Fprintf(&r.builder, ";38;2;%d;%d;%d", fg.Red, fg.Green, fg.Blue)
			case CMXterm:
				fmt.Fprintf(&r.builder, ";38;5;%d", fg.ToXtermColor())
			case CMAnsi:
				fmt.Fprintf(&r.builder, ";%d", 30+fg.ToAnsiColor())
			}

			if r.Mask == nil || !colorsEqual(r.Mask, &bg) {
				switch r.Mode {
				case CMTruecolor:
					fmt.Fprintf(&r.builder, ";48;2;%d;%d;%d", bg.Red, bg.Green, bg.Blue)
				case CMXterm:
					fmt.Fprintf(&r.builder, ";48;5;%d", bg.ToXtermColor())
				case CMAnsi:
					fmt.Fprintf(&r.builder, ";%d", 40+bg.ToAnsiColor())
				}
			}

			r.builder.WriteByte('m')
		}

		r.builder.WriteRune(diff.Cell.Char)
	}

	fmt.Fprintf(&r.builder, "%s[0m", r.Escape)
	_, _ = io.WriteString(output, r.builder.String())
	r.builder.Reset()

	if f, ok := output.(interface{ Flush() error }); ok {
		_ = f.Flush()
	}
}

func colorsEqual(a, b *Color) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}
