package tbregion

// TOP RIGHT BOTTOM LEFT
const Thin_0011 = '┐'
const Thin_1001 = '┘'
const Thin_0110 = '┌'
const Thin_1100 = '└'
const Thin_1111 = '┼'
const Thin_0101 = '─'
const Thin_1110 = '├'
const Thin_1011 = '┤'
const Thin_1101 = '┴'
const Thin_0111 = '┬'
const Thin_1010 = '│'
const Thin_0000 = ' '
const Thin_0001 = ' '
const Thin_0010 = ' '
const Thin_0100 = ' '
const Thin_1000 = ' '

var thin [16]rune

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

func InitBorder() {
	thin = [...]rune{
		Thin_0000, Thin_0001, Thin_0010, Thin_0011,
		Thin_0100, Thin_0101, Thin_0110, Thin_0111,
		Thin_1000, Thin_1001, Thin_1010, Thin_1011,
		Thin_1100, Thin_1101, Thin_1110, Thin_1111,
	}
}

type Borderable interface {
	GetSize() XY
	SetRune(x, y int, ru rune)
}

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

// Draw Line based on a [][]bool, starting from x, y
func GetThinLines(borderable Borderable, connections [][]bool) [][]rune {
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

func DrawThinLines(borderable Borderable, x, y int, connections [][]bool) {
	size := borderable.GetSize()
	lines := GetThinLines(borderable, connections)
	for lY, _ := range lines {
		actualY := lY + y
		if actualY < 0 || actualY >= size.Y {
			continue
		}
		for lX, _ := range lines[lY] {
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

func DrawHThinLine(startX, startY, length int, borderable Borderable) {
	if length < 0 {
		return
	}
	ru := GetThinLine(false, true, false, true)
	for x := startX; x < startX+length; x++ {
		borderable.SetRune(x, startY, ru)
	}
}

func DrawVThinLine(startX, startY, length int, borderable Borderable) {
	if length < 0 {
		return
	}
	ru := GetThinLine(true, false, true, false)
	for y := startY; y < startY+length; y++ {
		borderable.SetRune(startX, y, ru)
	}
}
