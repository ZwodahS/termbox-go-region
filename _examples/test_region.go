package main

import (
	"time"

	termbox "github.com/nsf/termbox-go"
	tbregion "github.com/zwodahs/termbox-go-region"
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	tbregion.InitRegion()
	defer termbox.Close()

	events := make(chan termbox.Event)
	go func() {
		for {
			events <- termbox.PollEvent()
		}
	}()

	update := make(chan int)
	go func() {
		for {
			update <- 1
			time.Sleep(1 * time.Second / 3)
		}
	}()

	termbox.SetInputMode(termbox.InputEsc)
	termbox.Flush()

	mainRegion := tbregion.NewRegion(100, 100, termbox.Cell{Ch: ' ', Fg: termbox.ColorDefault, Bg: termbox.ColorDefault})

	region1 := mainRegion.NewRegion(20, 10, termbox.Cell{Ch: ' ', Fg: termbox.ColorDefault, Bg: termbox.ColorRed})
	region1.DrawThinBorder()

	region2 := region1.NewRegion(10, 5, termbox.Cell{Ch: ' ', Fg: termbox.ColorDefault, Bg: termbox.ColorBlue})
	region2.DrawThinBorder()

	region3 := mainRegion.NewRegion(50, 2, termbox.Cell{Ch: ' ', Fg: termbox.ColorRed, Bg: termbox.ColorBlue})
	region3.SetText(0, 0, "This is a testing text", termbox.ColorDefault, termbox.ColorDefault)
	region3.SetPosition(0, 20)
	region4 := mainRegion.NewRegion(50, 2, termbox.Cell{Ch: ' ', Fg: termbox.ColorRed, Bg: termbox.ColorBlue})
	region4.SetText(0, 0, "This is a testing text", termbox.ColorDefault)
	region4.SetPosition(0, 24)
	region5 := mainRegion.NewRegion(50, 2, termbox.Cell{Ch: ' ', Fg: termbox.ColorRed, Bg: termbox.ColorBlue})
	region5.SetText(0, 0, "This is a testing text")
	region5.SetPosition(0, 28)

	selectedRegion := &region1

loop:
	for {
		select {
		case e := <-events:
			switch e.Type {
			case termbox.EventKey:
				switch e.Key {
				case termbox.KeyEsc:
					break loop
				case termbox.KeyArrowDown:
					position := (*selectedRegion).GetPosition()
					position.Add(0, 1)
					(*selectedRegion).SetPosition(position.X, position.Y)
				case termbox.KeyArrowUp:
					position := (*selectedRegion).GetPosition()
					position.Add(0, -1)
					(*selectedRegion).SetPosition(position.X, position.Y)
				case termbox.KeyArrowLeft:
					position := (*selectedRegion).GetPosition()
					position.Add(-1, 0)
					(*selectedRegion).SetPosition(position.X, position.Y)
				case termbox.KeyArrowRight:
					position := (*selectedRegion).GetPosition()
					position.Add(1, 0)
					(*selectedRegion).SetPosition(position.X, position.Y)
				default:
					switch e.Ch {
					case 'h':
						region2.Hidden = !region2.Hidden
					case '1':
						selectedRegion = &region1
					case '2':
						selectedRegion = &region2
					case '3':
						selectedRegion = &region3
					case '4':
						selectedRegion = &region4
					case '5':
						selectedRegion = &region5
					}
				}
			}
			mainRegion.Draw()
			termbox.Flush()
		case <-update:
			mainRegion.Draw()
			termbox.Flush()
		}
	}
}
