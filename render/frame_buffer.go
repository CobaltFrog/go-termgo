package render

import "fmt"

type FrameBuffer struct {
	Frame []Cell
	Width  int
	Height int
}

func newFrameBuffer(w, h int) (*FrameBuffer, error) {
	if w <= 0 || h <= 0 {
		return &FrameBuffer{}, fmt.Errorf("buffer size is zero or negative")
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
	for y := range fb.Height {
		for x := range fb.Width {
			fb.Frame[y * fb.Width + x]= EmptyCell()
		}
	}
}

// Changes cell depth, rune and color in provided point in frame buffer
// Returns true if cell values are changed, false otherwise, and on error
func (fb *FrameBuffer) UpdateCell(x, y int, z uint16, c rune, fg, bg Color) (bool, error) {
	if x < 0         ||
	   y < 0         ||
	   x >= fb.Width ||
	   y >= fb.Height {
		return false, fmt.Errorf("point out of frame bounds")
	}

	pos := y * fb.Width + x

	if c  != fb.Frame[pos].Char       ||
	   fg != fb.Frame[pos].Foreground ||
	   bg != fb.Frame[pos].Background ||
	   z  >= fb.Frame[pos].Depth {

		fb.Frame[pos].Char = c;
		fb.Frame[pos].Foreground = fg;
		fb.Frame[pos].Background = bg;
		fb.Frame[pos].Depth = z;

		return true, nil
	}

	return false, nil
}

func (fb *FrameBuffer) GetCell(x, y int) (Cell, error){
	if x < 0          ||
	   y < 0          ||
	   x >= fb.Width  ||
	   y >= fb.Height {
		return Cell{}, fmt.Errorf("point out of frame bounds")
	}

	pos := y * fb.Width + x
	return fb.Frame[pos], nil
}
