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

func GetThinConnection(top, right, bottom, left bool) rune {
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
	GetSize() [2]int
	SetRune(x, y int, ru rune)
}

func DrawThinBorder(borderable Borderable) {
	size := borderable.GetSize()
	left, right, top, bottom := 0, size[0]-1, 0, size[1]-1
	if left == right || top == bottom {
		return
	}
	borderable.SetRune(left, top, GetThinConnection(false, true, true, false))
	borderable.SetRune(right, top, GetThinConnection(false, false, true, true))
	borderable.SetRune(left, bottom, GetThinConnection(true, true, false, false))
	borderable.SetRune(right, bottom, GetThinConnection(true, false, false, true))
	DrawHLine(left+1, top, size[0]-2, borderable)
	DrawHLine(left+1, bottom, size[0]-2, borderable)
	DrawVLine(left, top+1, size[1]-2, borderable)
	DrawVLine(right, top+1, size[1]-2, borderable)
}

func getConnection(connections [][]bool, x, y int) bool {
	if x < 0 || x >= len(connections) || y < 0 || y >= len(connections[x]) {
		return false
	}
	return connections[x][y]
}

// Draw Connection based on a [][]bool, starting from x, y
func GetConnections(borderable Borderable, connections [][]bool) [][]rune {
	lines := make([][]rune, len(connections))
	for bX := 0; bX < len(connections); bX++ {
		lines[bX] = make([]rune, len(connections[bX]))
		for bY := 0; bY < len(connections[bX]); bY++ {
			lines[bX][bY] = GetThinConnection(
				getConnection(connections, bX, bY-1),
				getConnection(connections, bX+1, bY),
				getConnection(connections, bX, bY+1),
				getConnection(connections, bX-1, bY),
			)
		}
	}
	return lines
}

func DrawConnections(borderable Borderable, x, y int, connections [][]bool) {
	size := borderable.GetSize()
	lines := GetConnections(borderable, connections)
	for lX := 0; lX < len(lines); lX++ {
		actualX := lX + x
		if actualX < 0 || actualX >= size[0] {
			continue
		}
		for lY := 0; lY < len(lines[lX]); lY++ {
			actualY := lY + y
			if actualY < 0 || actualY >= size[1] {
				continue
			}
			if connections[lX][lY] {
				borderable.SetRune(actualX, actualY, lines[lX][lY])
			}
		}
	}
}

func DrawHLine(startX, startY, length int, borderable Borderable) {
	if length < 0 {
		return
	}
	ru := GetThinConnection(false, true, false, true)
	for x := startX; x < startX+length; x++ {
		borderable.SetRune(x, startY, ru)
	}
}

func DrawVLine(startX, startY, length int, borderable Borderable) {
	if length < 0 {
		return
	}
	ru := GetThinConnection(true, false, true, false)
	for y := startY; y < startY+length; y++ {
		borderable.SetRune(startX, y, ru)
	}
}
