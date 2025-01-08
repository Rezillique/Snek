package main

import (
	"math/rand"
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

func drawText(s tcell.Screen, x1, y1, x2, y2 int, text string) {
	row := y1
	col := x1
	style := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
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
			break
		}
		g.snakeBody.Update(GAME_WIDTH, GAME_HEIGHT, longerSnake)
		drawParts(g.Screen, g.snakeBody.Parts, g.FoodPos, g.snakeStyle, defStyle)
		drawText(g.Screen, 1, 1, 8+len(strconv.Itoa(g.Score)), 1, "Score: "+strconv.Itoa(g.Score))

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
	numbers := ""
	for i, num := range g.collected {
		if i > 0 {
			numbers += ","
		}
		numbers += strconv.Itoa(num)
	}
	drawText(g.Screen, GAME_WIDTH/2-25, GAME_HEIGHT/2, GAME_WIDTH/2+25, GAME_HEIGHT/2,
		"Game Over, Score: "+strconv.Itoa(g.Score)+" Numbers: "+numbers+" Play Again? y/n")
	g.Screen.Show()
}
