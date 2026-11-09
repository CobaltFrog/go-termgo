package render

import (
	"fmt"
)

type FrameBuffer struct {
	Frame []Cell
	Width  uint16
	Height uint16
}

func newFrameBuffer(w, h uint16) (*FrameBuffer, error) {
	if w <= 0 || h <= 0 {
		return &FrameBuffer{}, fmt.Errorf("frame buffer: size is zero or negative")
	}

	f := make([]Cell, w * h)
	for i := range f {
		f[i] = EmptyCell()
	}

	return &FrameBuffer{
		Frame: f,
		Width: w,
		Height: h,
	}, nil
}

func (fb *FrameBuffer) Clear() {
	emt := EmptyCell()
	for y := range fb.Height {
		for x := range fb.Width {
			cell := &fb.Frame[fb.idx(x, y)]
			cell.Char       = emt.Char
			cell.Depth      = emt.Depth
			cell.Style      = emt.Style
			cell.Foreground = emt.Foreground
			cell.Background = emt.Background
		}
	}
}

// Changes cell depth, rune and color in provided point in frame buffer
// Returns true if cdepth, rune and color in provided point in frame buffell values are changed, false otherwise, and on error
func (fb *FrameBuffer) UpdateCell(x, y uint16, c Cell) (bool, error) {
	if !fb.inBounds(x, y) {
		return  false, fmt.Errorf("put cell: point out of frame bounds")
	}

	pos := fb.idx(x, y)
	if  fb.isUpdated(pos, c) {
		fb.Frame[pos].Char       = c.Char;
		fb.Frame[pos].Depth      = c.Depth;
		fb.Frame[pos].Style      = c.Style
		fb.Frame[pos].Foreground = c.Foreground
		fb.Frame[pos].Background = c.Background;

		return true, nil
	}

	return false, nil
}

func (fb *FrameBuffer) isUpdated(idx uint16, c Cell) bool {
	return c.Depth      >= fb.Frame[idx].Depth  &&  (
		   c.Char       != fb.Frame[idx].Char       ||
		   c.Style      != fb.Frame[idx].Style      ||
		   c.Foreground != fb.Frame[idx].Foreground ||
		   c.Background != fb.Frame[idx].Background )
}

func (fb *FrameBuffer) getCell(x, y uint16) (*Cell, error){
	if !fb.inBounds(x, y) {
		   return nil, fmt.Errorf("get cell: point out of frame bounds")
	}

	return &fb.Frame[fb.idx(x, y)], nil
}

func (fb *FrameBuffer) inBounds(x, y uint16) bool {
	if x > fb.Width ||
	   y > fb.Height {
		return false
	}

	return true
}

func (fb *FrameBuffer) idx(x, y uint16) uint16 {
	return y * fb.Width + x
}
