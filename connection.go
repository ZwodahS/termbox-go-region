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
