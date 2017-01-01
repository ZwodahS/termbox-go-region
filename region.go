/**
MIT License

Copyright (c) 2017 ZwodahS

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package tbregion

import (
	termbox "github.com/nsf/termbox-go"
)

type XY [2]int

type Region struct {
	Hidden   bool
	Cells    [][]termbox.Cell // inner storage of cell
	regions  []*Region        // child regions
	width    int              // width of the region
	height   int              // height of the region
	position XY               // x, y
	parent   *Region          // parent region
}

// Create a new region
func NewRegion(width, height int, defaultCell termbox.Cell) *Region {
	region := &Region{width: width, height: height, position: [2]int{0, 0}, parent: nil}
	region.Cells = make([][]termbox.Cell, width)
	for x := 0; x < width; x++ {
		region.Cells[x] = make([]termbox.Cell, height)
		for y := 0; y < height; y++ {
			region.Cells[x][y] = defaultCell
		}
	}
	region.regions = make([]*Region, 0)
	return region
}

// Create a new region inside region
func (r *Region) NewRegion(width, height int, defaultCell termbox.Cell) *Region {
	region := NewRegion(width, height, defaultCell)
	region.parent = r
	r.regions = append(r.regions, region)
	return region
}

// Draw the cells in the region onto the termbox buffer.
// This will call the Draw method in all the child regions.
// Takes 4 arguments, x, y, width, height
func (r *Region) Draw(params ...int) {
	if r.Hidden {
		return
	}
	// Draw the current cell, then draw the regions, so they can be draw ontop of region.
	// Checking for what region is visible might be slower than just drawing over it,
	// so let's just draw over it.

	// x, y in this case is the position of the "parent" region.

	// calculate the position to start drawing
	// Is there a better way to do this
	x, y, width, height := 0, 0, -1, -1
	if len(params) > 0 {
		x = params[0]
	}
	if len(params) > 1 {
		y = params[1]
	}
	if len(params) > 2 {
		width = params[2]
	}
	if len(params) > 3 {
		height = params[3]
	}

	actualMinX, actualMinY, actualMaxX, actualMaxY := -1, -1, -1, -1
	if width != -1 {
		actualMinX = x
		actualMaxX = x + width - 1 //inclusive
	}
	if height != -1 {
		actualMinY = y
		actualMaxY = y + height - 1
	}

	// the actual starts to this region
	startX := x + r.position[0]
	startY := y + r.position[1]

	for rX := 0; rX < r.width; rX++ {
		actualX := startX + rX
		// hide the cell if out of the bound
		if width != -1 && (actualX < actualMinX || actualX > actualMaxX) {
			continue
		}
		for rY := 0; rY < r.height; rY++ {
			actualY := startY + rY
			// hide the cell if out of the bound
			if width != -1 && (actualY < actualMinY || actualY > actualMaxY) {
				continue
			}
			termbox.SetCell(actualX, actualY, r.Cells[rX][rY].Ch, r.Cells[rX][rY].Fg, r.Cells[rX][rY].Bg)
		}
	}
	// draw the child regions
	for _, region := range r.regions {
		region.Draw(startX, startY, r.width, r.height)
	}
}

// Draw a thin border in this region
// See DrawThinBorder(Borderable)
func (r *Region) DrawThinBorder() {
	DrawThinBorder(r)
}

// Get the Size of the region.
func (r *Region) GetSize() [2]int {
	return [2]int{r.width, r.height}
}

// Get the Position of the region
func (r *Region) GetPosition() XY {
	return r.position
}

// Setting position of this region with respect to parent.
func (r *Region) SetPosition(xy XY) {
	r.position = xy
}

// Check if a position is out of bound.
func (r *Region) IsOutOfBound(x, y int) bool {
	if x < 0 || x >= r.width {
		return true
	}
	if y < 0 || y >= r.height {
		return true
	}
	return false
}

// Set the cell value at this position.
func (r *Region) SetCell(x, y int, ru rune, fg, bg termbox.Attribute) {
	if r.IsOutOfBound(x, y) {
		return
	}
	r.Cells[x][y] = termbox.Cell{Ch: ru, Fg: fg, Bg: bg}
}

// Set the rune value at this position.
func (r *Region) SetRune(x, y int, ru rune) {
	if r.IsOutOfBound(x, y) {
		return
	}
	r.Cells[x][y].Ch = ru
}

// Set the foreground value at this position.
func (r *Region) SetForeground(x, y int, fg termbox.Attribute) {
	if r.IsOutOfBound(x, y) {
		return
	}
	r.Cells[x][y].Fg = fg
}

// Set the background value at this position.
func (r *Region) SetBackground(x, y int, bg termbox.Attribute) {
	if r.IsOutOfBound(x, y) {
		return
	}
	r.Cells[x][y].Bg = bg
}

// Fill the region with data
func (r *Region) Fill(ru rune, fg, bg termbox.Attribute) {
	for x := 0; x < r.width; x++ {
		for y := 0; y < r.height; y++ {
			r.Cells[x][y] = termbox.Cell{Ch: ru, Fg: fg, Bg: bg}
		}
	}
}

// Get the cell that is in this position.
// Returns actual reference to the cell, which can be modified directly.
func (r *Region) GetCell(x, y int) *termbox.Cell {
	if r.IsOutOfBound(x, y) {
		return nil
	}
	return &r.Cells[x][y]
}