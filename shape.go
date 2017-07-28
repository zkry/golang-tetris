package main

func moveShape(r, c int, s Shape) Shape {
	var newShape Shape
	for i := 0; i < 4; i++ {
		newShape[i].row = s[i].row + r
		newShape[i].col = s[i].col + c
	}
	return newShape
}

func moveShapeDown(s Shape) Shape {
	return moveShape(-1, 0, s)
}

func moveShapeRight(s Shape) Shape {
	return moveShape(0, 1, s)
}

func moveShapeLeft(s Shape) Shape {
	return moveShape(0, -1, s)
}

func isGameOver(s Shape) bool {
	for i := 0; i < 4; i++ {
		if s[i].row >= 20 {
			return true
		}
	}
	return false
}

func getShapeWidth(s Shape) int {
	maxWidth := 0
	for i := 1; i < 4; i++ {
		w := s[i].col - s[0].col
		if w > maxWidth {
			maxWidth = w
		}
	}
	return maxWidth
}

func getShapeHeight(s Shape) int {
	maxHeight := -1
	minHeight := 22
	for i := 0; i < 4; i++ {
		if s[i].row < minHeight {
			minHeight = s[i].row
		}
		if s[i].row > maxHeight {
			maxHeight = s[i].row
		}
	}
	return maxHeight - minHeight
}

func rotateShape(s Shape) Shape {
	var retShape Shape
	pivot := s[1]
	retShape[1] = pivot
	for i := 0; i < 4; i++ {
		// Index 1 is the pivot point
		if i == 1 {
			continue
		}
		dRow := pivot.row - s[i].row
		dCol := pivot.col - s[i].col
		retShape[i].row = pivot.row + (dCol * -1)
		retShape[i].col = pivot.col + (dRow)
	}
	return retShape
}

func getShapeFromPiece(p Piece) Shape {
	var retShape Shape
	switch p {
	case LPiece:
		retShape = Shape{
			Point{row: 1, col: 0},
			Point{row: 1, col: 1},
			Point{row: 1, col: 2},
			Point{row: 0, col: 0},
		}
	case IPiece:
		retShape = Shape{
			Point{row: 1, col: 0},
			Point{row: 1, col: 1},
			Point{row: 1, col: 2},
			Point{row: 1, col: 3},
		}
	case OPiece:
		retShape = Shape{
			Point{row: 1, col: 0},
			Point{row: 1, col: 1},
			Point{row: 0, col: 0},
			Point{row: 0, col: 1},
		}
	case TPiece:
		retShape = Shape{
			Point{row: 1, col: 0},
			Point{row: 1, col: 1},
			Point{row: 1, col: 2},
			Point{row: 0, col: 1},
		}
	case SPiece:
		retShape = Shape{
			Point{row: 0, col: 0},
			Point{row: 0, col: 1},
			Point{row: 1, col: 1},
			Point{row: 1, col: 2},
		}
	case ZPiece:
		retShape = Shape{
			Point{row: 1, col: 0},
			Point{row: 1, col: 1},
			Point{row: 0, col: 1},
			Point{row: 0, col: 2},
		}
	case JPiece:
		retShape = Shape{
			Point{row: 1, col: 0},
			Point{row: 0, col: 1},
			Point{row: 0, col: 0},
			Point{row: 0, col: 2},
		}
	default:
		panic("getShapeFromPiece(Piece): Invalid piece entered")
	}
	return retShape

}
