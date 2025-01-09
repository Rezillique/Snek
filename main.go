package main

import (
	"log"
	"os"

	"github.com/gdamore/tcell"
)

func main() {
	screen, err := tcell.NewScreen()

	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)

	snakeStyle := tcell.StyleDefault.Background(tcell.ColorGreen).Foreground(tcell.ColorGreen)
	foodStyle := tcell.StyleDefault.Background(tcell.ColorRed).Foreground(tcell.ColorRed)
	powerStyle := tcell.StyleDefault.Background(tcell.ColorYellow).Foreground(tcell.ColorYellow)

	game := Game{
		Screen:     screen,
		snakeStyle: snakeStyle,
		foodStyle:  foodStyle,
		powerStyle: powerStyle,
		defStyle:   defStyle,
	}

	// Start game loop in a goroutine
	go game.Run()

	// Handle input events
	for {
		switch event := game.Screen.PollEvent().(type) {
		case *tcell.EventResize:
			game.Screen.Sync()
		case *tcell.EventKey:
			if event.Key() == tcell.KeyEscape {
				game.Screen.Fini()
				os.Exit(0)
			}
			if game.GameOver {
				switch event.Key() {
				case tcell.KeyRune:
					switch event.Rune() {
					case 'y', 'Y':
						game.GameOver = false
						go game.Run()
					case 'n', 'N':
						game.Screen.Fini()
						os.Exit(0)
					}
				}
			} else {
				if event.Key() == tcell.KeyUp && game.snakeBody.Yspeed == 0 {
					game.snakeBody.ChangeDir(-1, 0)
				} else if event.Key() == tcell.KeyDown && game.snakeBody.Yspeed == 0 {
					game.snakeBody.ChangeDir(1, 0)
				} else if event.Key() == tcell.KeyLeft && game.snakeBody.Xspeed == 0 {
					game.snakeBody.ChangeDir(0, -1)
				} else if event.Key() == tcell.KeyRight && game.snakeBody.Xspeed == 0 {
					game.snakeBody.ChangeDir(0, 1)
				}
			}
		}
	}
}
