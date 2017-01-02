package tbregion

import (
	runewidth "github.com/mattn/go-runewidth"
	termbox "github.com/nsf/termbox-go"
)

type XY [2]int

func (xy XY) X() int {
	return xy[0]
}
func (xy XY) Y() int {
	return xy[1]
}

func (xy XY) Set(x, y int) XY {
	xy[0] = x
	xy[1] = y
	return xy
}

func (xy XY) Add(x, y int) XY {
	xy[0] += x
	xy[1] += y
	return xy
}

type Region struct {
	Hidden   bool
	Cells    [][]termbox.Cell // inner storage of cell
	regions  []*Region        // child regions
	width    int              // width of the region
	height   int              // height of the region
	position XY               // x, y
	parent   *Region          // parent region
	dirty    bool             // bool marking if this region is dirty
}

// Create a new region
func NewRegion(width, height int, defaultCell termbox.Cell) *Region {
	region := &Region{width: width, height: height, position: [2]int{0, 0}, parent: nil}
	region.Cells = make([][]termbox.Cell, height)
	for y := 0; y < height; y++ {
		region.Cells[y] = make([]termbox.Cell, width)
		for x := 0; x < width; x++ {
			region.Cells[y][x] = defaultCell
		}
	}
	region.regions = make([]*Region, 0)
	region.dirty = true
	return region
}

// Create a new region inside region
func (r *Region) NewRegion(width, height int, defaultCell termbox.Cell) *Region {
	region := NewRegion(width, height, defaultCell)
	region.parent = r
	r.regions = append(r.regions, region)
	r.dirty = true
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

	if r.dirty {
		for rY := range r.Cells {
			// for rY := 0; rY < r.height; rY++ {
			actualY := startY + rY
			// hide the cell if out of the bound
			if width != -1 && (actualY < actualMinY || actualY > actualMaxY) {
				continue
			}
			for rX := range r.Cells[rY] {
				// for rX := 0; rX < r.width; rX++ {
				actualX := startX + rX
				// hide the cell if out of the bound
				if width != -1 && (actualX < actualMinX || actualX > actualMaxX) {
					continue
				}
				termbox.SetCell(actualX, actualY, r.Cells[rY][rX].Ch, r.Cells[rY][rX].Fg, r.Cells[rY][rX].Bg)
			}
		}
	}
	// draw the child regions
	for _, region := range r.regions {
		if r.dirty { // if parent is dirty, we need to force all child to redraw
			region.dirty = true
		}
		region.Draw(startX, startY, r.width, r.height)
	}
	r.dirty = false
}

// Draw a thin border in this region
// See DrawThinBorder(Borderable)
func (r *Region) DrawThinBorder() {
	DrawThinBorder(r)
}

// Get the Size of the region.
func (r *Region) GetSize() XY {
	return XY{r.width, r.height}
}

// Get the Position of the region
func (r *Region) GetPosition() XY {
	return r.position
}

// Setting position of this region with respect to parent.
func (r *Region) SetPosition(xy XY) {
	r.position = xy
	r.dirty = true
	if r.parent != nil {
		r.parent.dirty = true
	}
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
	r.Cells[y][x] = termbox.Cell{Ch: ru, Fg: fg, Bg: bg}
	r.dirty = true
}

func (r *Region) SetText(x, y int, str string, fg, bg termbox.Attribute) {
	drawX := x
	for _, value := range str {
		r.SetCell(drawX, y, value, fg, bg)
		drawX += runewidth.RuneWidth(value)
	}
}

// Set the rune value at this position.
func (r *Region) SetRune(x, y int, ru rune) {
	if r.IsOutOfBound(x, y) {
		return
	}
	r.Cells[y][x].Ch = ru
	r.dirty = true
}

// Set the foreground value at this position.
func (r *Region) SetForeground(x, y int, fg termbox.Attribute) {
	if r.IsOutOfBound(x, y) {
		return
	}
	r.Cells[y][x].Fg = fg
	r.dirty = true
}

// Set the background value at this position.
func (r *Region) SetBackground(x, y int, bg termbox.Attribute) {
	if r.IsOutOfBound(x, y) {
		return
	}
	r.Cells[y][x].Bg = bg
	r.dirty = true
}

// Fill the region with data
func (r *Region) Fill(ru rune, fg, bg termbox.Attribute) {
	for x := 0; x < r.width; x++ {
		for y := 0; y < r.height; y++ {
			r.Cells[y][x] = termbox.Cell{Ch: ru, Fg: fg, Bg: bg}
		}
	}
	r.dirty = true
}
