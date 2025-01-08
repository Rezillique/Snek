package main

type SnakePart struct {
	X int
	Y int
}

type SnakeBody struct {
	Parts  []SnakePart
	Xspeed int
	Yspeed int
}

func (sb *SnakeBody) ChangeDir(vertical int, horizontal int) {
	sb.Yspeed = vertical
	sb.Xspeed = horizontal
}

func (sb *SnakeBody) Update(width int, height int, longerSnake bool) {
	sb.Parts = append(sb.Parts, sb.Parts[len(sb.Parts)-1].GetUpdatedPart(sb, width, height))
	if !longerSnake {
		sb.Parts = sb.Parts[1:]
	}
}

func (sb *SnakeBody) ResetPos(width int, height int) {
	snakeParts := []SnakePart{
		{
			X: int(width / 2),
			Y: int(width / 2),
		},
		{
			X: int(width/2) + 1,
			Y: int(height / 2),
		},
		{
			X: int(width/2) + 2,
			Y: int(height / 2),
		},
	}
	sb.Parts = snakeParts
	sb.Xspeed = 1
	sb.Yspeed = 0
}

func (sp *SnakePart) GetUpdatedPart(sb *SnakeBody, width int, height int) SnakePart {
	newPart := *sp
	newX := newPart.X + sb.Xspeed
	newY := newPart.Y + sb.Yspeed

	// Keep snake within boundaries
	if newX <= 0 {
		newX = 1
	} else if newX >= width-1 {
		newX = width - 2
	}

	if newY <= 0 {
		newY = 1
	} else if newY >= height-1 {
		newY = height - 2
	}

	newPart.X = newX
	newPart.Y = newY
	return newPart
}
