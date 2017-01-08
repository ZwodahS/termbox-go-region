package tbregion

// Thin (TOP RIGHT BOTTOM LEFT)
const (
	Thin0011 = '┐'
	Thin1001 = '┘'
	Thin0110 = '┌'
	Thin1100 = '└'
	Thin1111 = '┼'
	Thin0101 = '─'
	Thin1110 = '├'
	Thin1011 = '┤'
	Thin1101 = '┴'
	Thin0111 = '┬'
	Thin1010 = '│'
	Thin0000 = ' '
	Thin0001 = ' '
	Thin0010 = ' '
	Thin0100 = ' '
	Thin1000 = ' '
)

var thin [16]rune

// GetThinLine get the correct thin border to draw based on if top, right, bottom, left has connections
func GetThinLine(top, right, bottom, left bool) rune {
	intValue := 0
	if top {
		intValue += 8
	}
	if right {
		intValue += 4
	}
	if bottom {
		intValue += 2
	}
	if left {
		intValue += 1
	}
	return thin[intValue]
}

func initThinLines() {
	thin = [...]rune{
		Thin0000, Thin0001, Thin0010, Thin0011,
		Thin0100, Thin0101, Thin0110, Thin0111,
		Thin1000, Thin1001, Thin1010, Thin1011,
		Thin1100, Thin1101, Thin1110, Thin1111,
	}
}

//Borderable is interface that can border can be drawn on it
type Borderable interface {
	GetSize() XY
	SetRune(x, y int, ru rune)
}

//DrawThinBorder draws thin border around the borderable
func DrawThinBorder(borderable Borderable) {
	size := borderable.GetSize()
	left, right, top, bottom := 0, size.X-1, 0, size.Y-1
	if left == right || top == bottom {
		return
	}
	borderable.SetRune(left, top, GetThinLine(false, true, true, false))
	borderable.SetRune(right, top, GetThinLine(false, false, true, true))
	borderable.SetRune(left, bottom, GetThinLine(true, true, false, false))
	borderable.SetRune(right, bottom, GetThinLine(true, false, false, true))
	DrawHThinLine(left+1, top, size.X-2, borderable)
	DrawHThinLine(left+1, bottom, size.X-2, borderable)
	DrawVThinLine(left, top+1, size.Y-2, borderable)
	DrawVThinLine(right, top+1, size.Y-2, borderable)
}

func getLine(connections [][]bool, x, y int) bool {
	if y < 0 || y >= len(connections) || x < 0 || x >= len(connections[y]) {
		return false
	}
	return connections[y][x]
}

// GetThinLines returns the runes to be draw given the connections
func GetThinLines(connections [][]bool) [][]rune {
	lines := make([][]rune, len(connections))
	for bY := 0; bY < len(connections); bY++ {
		lines[bY] = make([]rune, len(connections[bY]))
		for bX := 0; bX < len(connections[bY]); bX++ {
			lines[bY][bX] = GetThinLine(
				getLine(connections, bX, bY-1),
				getLine(connections, bX+1, bY),
				getLine(connections, bX, bY+1),
				getLine(connections, bX-1, bY),
			)
		}
	}
	return lines
}

// DrawThinLines onto a borderable based on a connections starting from x, y
func DrawThinLines(borderable Borderable, x, y int, connections [][]bool) {
	size := borderable.GetSize()
	lines := GetThinLines(connections)
	for lY := range lines {
		actualY := lY + y
		if actualY < 0 || actualY >= size.Y {
			continue
		}
		for lX := range lines[lY] {
			actualX := lX + x
			if actualX < 0 || actualX >= size.X {
				continue
			}
			if connections[lY][lX] {
				borderable.SetRune(actualX, actualY, lines[lY][lX])
			}
		}
	}
}

// DrawHThinLine draws a horizontal line starting from startX, startY with a specific length
func DrawHThinLine(startX, startY, length int, borderable Borderable) {
	if length < 0 {
		return
	}
	ru := GetThinLine(false, true, false, true)
	for x := startX; x < startX+length; x++ {
		borderable.SetRune(x, startY, ru)
	}
}

// DrawVThinLine draws a vertical line starting from startX, startY with a specific length
func DrawVThinLine(startX, startY, length int, borderable Borderable) {
	if length < 0 {
		return
	}
	ru := GetThinLine(true, false, true, false)
	for y := startY; y < startY+length; y++ {
		borderable.SetRune(startX, y, ru)
	}
}
