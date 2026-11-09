package render

import (
	"fmt"
	"os"
)

type FrameTarget interface {
	BeginFrame()
	Update(x, y, z uint16, c rune, fg, bg Color, style GraphicsMode)
	EndFrame() []CellDiff
	Swap()
	GetSize() (width, height uint16)
}

type RenderTarget struct {
	front *FrameBuffer 
	back  *FrameBuffer
}

func NewRenderTarget(w, h uint16) (*RenderTarget, error) {
	if w <= 0 || h <= 0 {
		return &RenderTarget{}, fmt.Errorf("frame size is zero or negative")
	}

	p_fr, err := newFrameBuffer(w, h)
	if  err != nil {
		return &RenderTarget{}, fmt.Errorf("front buffer: %v\n", err)
	}

	p_bk, err := newFrameBuffer(w, h)
	if ; err != nil {
		return &RenderTarget{}, fmt.Errorf("back buffer: %v\n", err)
	}

	return &RenderTarget {
		front: p_fr,
		back: p_bk,
	}, nil
}

func (rt *RenderTarget) GetSize() (width, height uint16) {
	return rt.front.Width, rt.front.Height
}

func (rt *RenderTarget) BeginFrame() {
	rt.back.Clear()
}

func (rt *RenderTarget) Update(x, y uint16, c Cell) (bool, error) {
	return rt.back.UpdateCell(x, y, c)
}

func (rt *RenderTarget) GetFrontCell(x, y uint16) (Cell, error) {
	p, err := rt.front.getCell(x, y)
	if err != nil {
		return  Cell{}, err
	}

	return Cell {
		Char: p.Char,
		Depth: p.Depth,
		Foreground: p.Foreground,
		Background: p.Background,
		Style: p.Style,
	}, nil
}

func (rt *RenderTarget) EndFrame() []CellDiff {
	res := make([]CellDiff, rt.front.Height * rt.front.Width)
	for y := range rt.back.Height {
		for x := range rt.front.Width {
			pos := y * rt.front.Width + x
			old, err := rt.front.getCell(x, y);
			if  err != nil {
				fmt.Fprintf(os.Stderr, "end frame: %v\n", err)
				return res
			}

			nw, err := rt.back.getCell(x, y);
			if  err != nil {
				fmt.Fprintf(os.Stderr, "end frame: %v\n", err)
				return res
			}

			if  nw.Depth       >= old.Depth     &&(
				old.Char       != nw.Char       ||
			    old.Style      != nw.Style      ||
			    old.Foreground != nw.Foreground ||
			    old.Background != nw.Background ){

				res[pos].Cell = nw
				res[pos].X = x
				res[pos].Y = y
			}
		}
	}

	return res
}

func (rt *RenderTarget) Swap() {
	rt.back, rt.front = rt.front, rt.back
}
