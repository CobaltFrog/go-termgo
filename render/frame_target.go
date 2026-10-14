package render

import (
	"fmt"
	"os"
)

type FrameTarget interface {
	BeginFrame()
	Change(x, y int, z uint16, c rune, fg, bg Color)
	EndFrame() []CellDiff
	Swap()
}

type RenderTarget struct {
	front *FrameBuffer 
	back *FrameBuffer
}

func NewRenderTarget(w, h int) (*RenderTarget, error) {
	if w <= 0 || h <= 0 {
		return &RenderTarget{}, fmt.Errorf("frame size is zero or negative")
	}

	p_fr, err := newFrameBuffer(w, h)
	if  err != nil {
		return &RenderTarget{}, fmt.Errorf("render target: front buffer: %v\n", err)
	}

	p_bk, err := newFrameBuffer(w, h)
	if ; err != nil {
		return &RenderTarget{}, fmt.Errorf("render target: back buffer: %v\n", err)
	}

	return &RenderTarget {
		front: p_fr,
		back: p_bk,
	}, nil
}

func (rt *RenderTarget) BeginFrame() {
	rt.back.Clear()
}

func (rt *RenderTarget) Change(x, y int, z uint16, c rune, fg, bg Color) {
	rt.back.UpdateCell(x, y, z, c, fg, bg)
}

func (rt *RenderTarget) EndFrame() []CellDiff {
	res := make([]CellDiff, rt.front.Height * rt.front.Width)
	for y := range rt.back.Height {
		for x := range rt.front.Width {
			old, err := rt.front.GetCell(x, y);
			if  err != nil {
				fmt.Fprintf(os.Stderr, "cell access: %v\n", err)
				return res
			}

			nw, err := rt.back.GetCell(x, y);
			if  err != nil {
				fmt.Fprintf(os.Stderr, "cell access: %v\n", err)
				return res
			}

			if old.Char       != nw.Char       ||
			   old.Foreground != nw.Foreground ||
			   old.Background != nw.Background ||
			   old.Depth      <  nw.Depth {

				res[y * rt.front.Width + x].Cell = &nw
			}
		}
	}

	return res
}

func (rt *RenderTarget) Swap() {
	rt.back, rt.front = rt.front, rt.back
}
