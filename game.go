package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/gdamore/tcell"
)

type Game struct {
	Screen     tcell.Screen
	snakeBody  SnakeBody
	FoodPos    SnakePart
	WallParts  []SnakePart
	Score      int
	GameOver   bool
	snakeStyle tcell.Style
	foodStyle  tcell.Style
	powerStyle tcell.Style
	wallStyle  tcell.Style
	gridValues map[SnakePart]int // Store grid numbers
	collected  []int             // Store collected numbers
	defStyle   tcell.Style
}

const (
	GAME_WIDTH  = 32
	GAME_HEIGHT = 32
	GAME_SIZE   = 32
	GAME_SPEED  = 80
	GRID_CHAR   = '·' // Unicode middle dot forgrid points
	MAX_NUMBER  = 9   // Maximum random number for grid
)

func drawParts(s tcell.Screen, snakeParts []SnakePart, foodPos SnakePart, snakeStyle tcell.Style, foodStyle tcell.Style) {
	s.SetContent(foodPos.X, foodPos.Y, '\u25CF', nil, foodStyle)
	for _, part := range snakeParts {
		s.SetContent(part.X, part.Y, ' ', nil, snakeStyle)
	}
}

func drawText(s tcell.Screen, x1, y1, x2, y2 int, text string, style tcell.Style) {
	row := y1
	col := x1
	for _, r := range text {
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

func checkCollision(parts []SnakePart, otherPart SnakePart) bool {
	for _, part := range parts {
		if part.X == otherPart.X && part.Y == otherPart.Y {
			return true
		}
	}
	return false
}

func (g *Game) UpdateFoodPos(width int, height int) {
	g.FoodPos.X = rand.Intn(width-2) + 1
	g.FoodPos.Y = rand.Intn(height-2) + 1
	if g.FoodPos.Y == 1 && g.FoodPos.X < 10 || checkCollision(g.WallParts, g.FoodPos) {
		g.UpdateFoodPos(width, height)
	}
}

func (g *Game) AddWalls(width, height int) {
	if g.Score > 0 && g.Score%5 == 0 {
		for i := 0; i < 2; i++ {
			wall := SnakePart{
				X: rand.Intn(width-2) + 1,
				Y: rand.Intn(height-2) + 1,
			}
			if !checkCollision(g.snakeBody.Parts, wall) &&
				wall != g.FoodPos &&
				!(wall.Y == 1 && wall.X < 10) {
				g.WallParts = append(g.WallParts, wall)
			}
		}
	}
}

func (g *Game) generatePassword() string {
	// Simple password generation based on score and timestamp
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%d%d", g.Score, timestamp%1000)
}

func (g *Game) savePassword(password string) error {
	return os.WriteFile("password.txt", []byte(password), 0644)
}

func (g *Game) Run() {
	// Initialize grid values
	g.gridValues = make(map[SnakePart]int)
	g.collected = make([]int, 0)

	// Generate random numbers for grid positions
	for x := 1; x < GAME_SIZE-1; x++ {
		for y := 1; y < GAME_SIZE-1; y++ {
			g.gridValues[SnakePart{X: x, Y: y}] = rand.Intn(MAX_NUMBER) + 1
		}
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	gridStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGray)
	g.wallStyle = tcell.StyleDefault.Background(tcell.ColorRed).Foreground(tcell.ColorRed)
	g.Screen.SetStyle(defStyle)

	// Initialize border walls
	g.WallParts = []SnakePart{}
	// Add top and bottom borders
	for x := 0; x < GAME_WIDTH; x++ {
		g.WallParts = append(g.WallParts, SnakePart{X: x, Y: 0})           // Top border
		g.WallParts = append(g.WallParts, SnakePart{X: x, Y: GAME_HEIGHT}) // Bottom border
	}
	// Add left and right borders
	for y := 0; y < GAME_HEIGHT; y++ {
		g.WallParts = append(g.WallParts, SnakePart{X: 0, Y: y})          // Left border
		g.WallParts = append(g.WallParts, SnakePart{X: GAME_WIDTH, Y: y}) // Right border
	}

	g.snakeBody.ResetPos(GAME_WIDTH, GAME_HEIGHT)
	g.UpdateFoodPos(GAME_WIDTH, GAME_HEIGHT)
	g.GameOver = false
	g.Score = 0
	for {
		longerSnake := false
		g.Screen.Clear()

		// Draw grid first (so it appears behind everything else)
		for x := 1; x < GAME_SIZE; x++ {
			for y := 1; y < GAME_SIZE; y++ {
				// Skip drawing grid where snake, food, or walls are
				if !checkCollision(g.snakeBody.Parts, SnakePart{X: x, Y: y}) &&
					!checkCollision(g.WallParts, SnakePart{X: x, Y: y}) &&
					!(g.FoodPos.X == x && g.FoodPos.Y == y) {
					g.Screen.SetContent(x, y, GRID_CHAR, nil, gridStyle)
				}
			}
		}

		if checkCollision(g.snakeBody.Parts[len(g.snakeBody.Parts)-1:], g.FoodPos) {
			g.UpdateFoodPos(GAME_WIDTH, GAME_HEIGHT)
			longerSnake = true
			g.Score++
			g.snakeStyle = g.powerStyle
			go func() {
				time.Sleep(2 * time.Second)
				g.snakeStyle = tcell.StyleDefault.Background(tcell.ColorGreen).Foreground(tcell.ColorGreen)
			}()
		}
		if checkCollision(g.snakeBody.Parts[:len(g.snakeBody.Parts)-1], g.snakeBody.Parts[len(g.snakeBody.Parts)-1]) {
			g.GameOver = true
			g.DrawGameOver()
			break
		}
		g.snakeBody.Update(GAME_WIDTH, GAME_HEIGHT, longerSnake)
		drawParts(g.Screen, g.snakeBody.Parts, g.FoodPos, g.snakeStyle, defStyle)
		drawText(g.Screen, 1, 1, 8+len(strconv.Itoa(g.Score)), 1, "Score: "+strconv.Itoa(g.Score), defStyle)

		// Draw walls
		for _, wall := range g.WallParts {
			g.Screen.SetContent(wall.X, wall.Y, '\u2588', nil, g.wallStyle)
		}

		// After snake moves, collect number if exists
		head := g.snakeBody.Parts[len(g.snakeBody.Parts)-1]
		if val, exists := g.gridValues[head]; exists {
			g.collected = append(g.collected, val)
			delete(g.gridValues, head) // Remove collected number
		}

		time.Sleep(GAME_SPEED * time.Millisecond)
		g.Screen.Show()

		if checkCollision(g.WallParts, g.snakeBody.Parts[len(g.snakeBody.Parts)-1]) {
			break
		}

		if longerSnake {
			g.AddWalls(GAME_WIDTH, GAME_HEIGHT)
		}
	}
	g.GameOver = true
	// Generate and save new password
	password := g.generatePassword()
	if err := g.savePassword(password); err != nil {
		log.Printf("Failed to save password: %v", err)
	}

}

func (g *Game) DrawGameOver() {
	g.Screen.Clear()

	// Calculate center position
	width, height := g.Screen.Size()
	centerX := width / 2
	centerY := height / 2

	// Draw text
	gameOverText := "GAME OVER"
	scoreText := fmt.Sprintf("Final Score: %d", g.Score)
	continueText := "Press Y to play again, N to quit"

	// Draw game over box
	boxStyle := g.defStyle
	boxWidth := len(continueText) + 4 // Add padding
	boxHeight := 7                    // Height for 3 lines of text plus padding

	// Calculate box position
	boxStartX := centerX - boxWidth/2
	boxStartY := centerY - boxHeight/2

	// Draw box borders and text
	for x := boxStartX; x < boxStartX+boxWidth; x++ {
		for y := boxStartY; y < boxStartY+boxHeight; y++ {
			if x == boxStartX || x == boxStartX+boxWidth-1 ||
				y == boxStartY || y == boxStartY+boxHeight-1 {
				g.Screen.SetContent(x, y, '█', nil, boxStyle)
			}
		}
	}

	// Draw text centered in box
	drawText(g.Screen,
		centerX-len(gameOverText)/2, boxStartY+2,
		centerX+len(gameOverText), boxStartY+2,
		gameOverText, g.defStyle)

	drawText(g.Screen,
		centerX-len(scoreText)/2, boxStartY+3,
		centerX+len(scoreText), boxStartY+3,
		scoreText, g.defStyle)

	drawText(g.Screen,
		centerX-len(continueText)/2, boxStartY+4,
		centerX+len(continueText)/2, boxStartY+4,
		continueText, g.defStyle)

	g.Screen.Show()
}

func (g *Game) DrawMainMenu() {
	g.Screen.Clear()
	width, height := g.Screen.Size()
	centerX := width / 2

	// Draw SNEK header with different colors
	headerText := "SNEK"
	headerColors := []tcell.Color{
		tcell.ColorGreen,
		tcell.ColorYellow,
		tcell.ColorRed,
		tcell.ColorBlue,
	}

	headerY := height / 4
	for i, char := range headerText {
		style := tcell.StyleDefault.
			Background(tcell.ColorBlack).
			Foreground(headerColors[i])
		g.Screen.SetContent(centerX-2+i, headerY, char, nil, style)
	}
}
