package main

import (
	"log"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

func main() {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	boxStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorReset)
	
	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
		os.Exit(1)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
		os.Exit(1)
	}
	
	s.SetStyle(defStyle)
	s.Clear()
	
	updateCh := make(chan struct{}, 1)
	
	updateCh <- struct{}{}
	
	go func() {
		for {
			select {
			case <-updateCh:
			case <-time.After(1 * time.Second):
			}
			
			s.Clear()
			
			xmax, ymax := s.Size()
			middleX, middleY := xmax/2, ymax/2
			
			// Draw tip
			tip := "Press CTRL+C to quit..."
			drawText(s, xmax-len(tip), ymax-1, xmax, ymax, tcell.StyleDefault.Foreground(tcell.ColorDarkGray).Background(tcell.ColorReset), tip)
			
			// Draw time
			currentTime := time.Now().Format("3:04:05 PM")
			drawText(s, middleX-len(currentTime)/2, middleY, xmax, ymax, boxStyle, currentTime)
			
			s.Show()
		}
	}()
	
	quit := func() {
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()
	
	// Event loop
	for {
		// Poll event
		ev := s.PollEvent()
		
		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
			select {
			case updateCh <- struct{}{}:
			default:
			}
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return
			} else if ev.Key() == tcell.KeyCtrlL {
				s.Sync()
				updateCh <- struct{}{}
			} else if ev.Rune() == 'C' || ev.Rune() == 'c' {
				s.Clear()
				updateCh <- struct{}{}
			}
		}
	}
}
