package tbregion

import (
	runewidth "github.com/mattn/go-runewidth"
	termbox "github.com/nsf/termbox-go"
)

// XY is a struct for x, y position
type XY struct {
	X int
	Y int
}

// Set the position of the XY and return itself
func (xy *XY) Set(x, y int) *XY {
	xy.X = x
	xy.Y = y
	return xy
}

// Add x y to a XY and return itself
func (xy *XY) Add(x, y int) *XY {
	xy.X += x
	xy.Y += y
	return xy
}

// Region is a representation of a region on the screen
type Region struct {
	Hidden   bool             // bool if the region is hidden
	Cells    [][]termbox.Cell // inner storage of cell
	regions  []*Region        // child regions
	width    int              // width of the region
	height   int              // height of the region
	position XY               // x, y
	parent   *Region          // parent region
	dirty    bool             // bool marking if this region is dirty
}

// NewRegion creates a new Region struct
// cells is a single cell stating the default value of all the cell in the region
func NewRegion(width, height int, cells ...termbox.Cell) *Region {
	region := &Region{width: width, height: height, position: XY{X: 0, Y: 0}, parent: nil}
	region.Cells = make([][]termbox.Cell, height)
	defaultCell := termbox.Cell{Ch: ' ', Fg: termbox.ColorWhite, Bg: termbox.ColorBlack}
	if len(cells) > 0 {
		defaultCell = cells[0]
	}
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

// NewRegion create a new region inside region
func (r *Region) NewRegion(width, height int, cells ...termbox.Cell) *Region {
	region := NewRegion(width, height, cells...)
	region.parent = r
	r.regions = append(r.regions, region)
	r.dirty = true
	return region
}

// RemoveRegion remove a region from itself if it contains the region
// return bool if the region was removed, false otherwise
func (r *Region) RemoveRegion(region *Region) bool {
	index := r.GetRegionIndex(region)
	if index == -1 {
		return false
	}
	// Do we need to optimize ?
	r.regions = append(r.regions[:index], r.regions[index+1:]...)
	return true
}

// RemoveAllRegions removes all regions
func (r *Region) RemoveAllRegions() {
	r.regions = make([]*Region, 0)
}

// Close closes the region.
func (r *Region) Close() {
	if r.parent == nil {
		return
	}
	r.parent.RemoveRegion(r)
	// any clean up ?
}

// GetRegionIndex return the index of another region in this region
// return -1 if not found
func (r *Region) GetRegionIndex(region *Region) int {
	for ind, value := range r.regions {
		if value == region {
			return ind
		}
	}
	return -1
}

// MarkForRedraw marks this region for redraw
func (r *Region) MarkForRedraw() {
	r.dirty = true
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
	startX := x + r.position.X
	startY := y + r.position.Y

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

// DrawThinBorder draws a thin border in the border cell of this region
// See DrawThinBorder(Borderable)
func (r *Region) DrawThinBorder() {
	DrawThinBorder(r)
}

// GetSize returns the size of the region
func (r *Region) GetSize() XY {
	return XY{r.width, r.height}
}

// GetPosition returns the position of the region
func (r *Region) GetPosition() XY {
	return r.position
}

// SetPosition sets the position of this region with relative the parent region
func (r *Region) SetPosition(x, y int) {
	r.position = XY{X: x, Y: y}
	r.dirty = true
	if r.parent != nil {
		r.parent.dirty = true
	}
}

// IsOutOfBound returns if this x, y value is out of bound of this region
func (r *Region) IsOutOfBound(x, y int) bool {
	if x < 0 || x >= r.width {
		return true
	}
	if y < 0 || y >= r.height {
		return true
	}
	return false
}

// SetCell sets the cell in position x, y of this region
// atttributes are optional
// first attribute is foreground, second attribute is background
func (r *Region) SetCell(x, y int, ru rune, attributes ...termbox.Attribute) {
	if r.IsOutOfBound(x, y) {
		return
	}
	r.Cells[y][x].Ch = ru
	if len(attributes) >= 2 {
		r.Cells[y][x].Fg = attributes[0]
		r.Cells[y][x].Bg = attributes[1]
	} else if len(attributes) == 1 {
		r.Cells[y][x].Fg = attributes[0]
	}
	r.dirty = true
}

// SetText set the text starting from x, y of this region
// atttributes are optional
// first attribute is foreground, second attribute is background
func (r *Region) SetText(x, y int, str string, attributes ...termbox.Attribute) {
	drawX := x
	for _, value := range str {
		r.SetCell(drawX, y, value, attributes...)
		drawX += runewidth.RuneWidth(value)
	}
}

// SetTextCenter set the text at the center of row y of this region
// atttributes are optional
// first attribute is foreground, second attribute is background
func (r *Region) SetTextCenter(y int, str string, attributes ...termbox.Attribute) {
	// calculate the start point to draw
	drawX := r.width/2 - runewidth.StringWidth(str)/2
	for _, value := range str {
		r.SetCell(drawX, y, value, attributes...)
		drawX += runewidth.RuneWidth(value)
	}
}

// SetRune set the rune value at position x, y
func (r *Region) SetRune(x, y int, ru rune) {
	if r.IsOutOfBound(x, y) {
		return
	}
	r.Cells[y][x].Ch = ru
	r.dirty = true
}

// SetForeground set the foreground value at this position.
func (r *Region) SetForeground(x, y int, fg termbox.Attribute) {
	if r.IsOutOfBound(x, y) {
		return
	}
	r.Cells[y][x].Fg = fg
	r.dirty = true
}

// SetBackground Set the background value at this position.
func (r *Region) SetBackground(x, y int, bg termbox.Attribute) {
	if r.IsOutOfBound(x, y) {
		return
	}
	r.Cells[y][x].Bg = bg
	r.dirty = true
}

// Fill set all the cell in this region to a specific value
func (r *Region) Fill(ru rune, attributes ...termbox.Attribute) {
	for x := 0; x < r.width; x++ {
		for y := 0; y < r.height; y++ {
			r.SetCell(x, y, ru, attributes...)
		}
	}
	r.dirty = true
}

// InitRegion initialize the library
func InitRegion() error {
	initThinLines()
	return nil
}
