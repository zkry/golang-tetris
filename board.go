package main

import (
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// isTouchingFloor checks if the piece that the user is controlling has a piece
// directly below it. Used to give the user more time when placing block on
// floor
func (b *Board) isTouchingFloor() bool {
	blockType := b[activeShape[0].row][activeShape[0].col]
	b.drawPiece(activeShape, Empty)
	isTouching := b.checkCollision(moveShapeDown(activeShape))
	b.drawPiece(activeShape, blockType)
	return isTouching
}

// rotatePiece rotates the piece that the user is currently moving clockwise by
// 90 degrees. The rotation is made and collision is checked. If the rotation can
// be completed by moving the newly rotated shape, the rotation will also be
// performed. If it is impossible to rotate, does nothing.
func (b *Board) rotatePiece() {
	// The O piece should not be rotated
	if currentPiece == OPiece {
		return
	}
	blockType := b[activeShape[0].row][activeShape[0].col]
	// Erase Piece
	b.drawPiece(activeShape, Empty)

	// Get the new shape and check for it's collision
	newShape := rotateShape(activeShape)
	if b.checkCollision(newShape) {
		if !b.checkCollision(moveShapeRight(newShape)) {
			newShape = moveShapeRight(newShape)
		} else if !b.checkCollision(moveShapeLeft(newShape)) {
			newShape = moveShapeLeft(newShape)
		} else if !b.checkCollision(moveShapeDown(newShape)) {
			newShape = moveShapeDown(newShape)
		} else {
			b.drawPiece(activeShape, blockType)
			return
		}
	}
	activeShape = newShape
	b.drawPiece(activeShape, blockType)
}

// movePiece attemps to move the piece that the user is controlling either
// right or left. +1 signifies a right move while -1 signifies a left move
func (b *Board) movePiece(dir int) {
	blockType := b[activeShape[0].row][activeShape[0].col]

	// Erase old piece
	b.drawPiece(activeShape, Empty)

	// Check collision
	didCollide := b.checkCollision(moveShape(0, dir, activeShape))
	if !didCollide {
		activeShape = moveShape(0, dir, activeShape)
	}
	b.drawPiece(activeShape, blockType)
}

// drawPiece sets the values of a board, b, to a specific block type, t
// according to shape, s.
func (b *Board) drawPiece(s Shape, t Block) {
	for i := 0; i < 4; i++ {
		b[activeShape[i].row][activeShape[i].col] = t
	}
}

// checkCollision checks if at the 4 points of a shape, s, there is
// nothing but Empty value under it and the position of the shape
// is inside the playing board (10x22 (top two rows invisiable)).
func (b Board) checkCollision(s Shape) bool {
	for i := 0; i < 4; i++ {
		r := s[i].row
		c := s[i].col
		if r < 0 || r > 21 || c < 0 || c > 9 || b[r][c] != Empty {
			return true
		}
	}
	return false
}

// applyGravity is the function that moves a piece down. If a collision
// is detected place the piece down and add a new piece. Returns wheather
// a collision was made.
func (b *Board) applyGravity() bool {
	blockType := b[activeShape[0].row][activeShape[0].col]
	// Erase old piece
	b.drawPiece(activeShape, Empty)

	// Does the block collide if it moves down?
	didCollide := b.checkCollision(moveShapeDown(activeShape))

	if !didCollide {
		activeShape = moveShapeDown(activeShape)
	}

	b.drawPiece(activeShape, blockType)

	if didCollide {
		if isGameOver(activeShape) {
			gameOver = true
		}
		b.checkRowCompletion(activeShape)
		b.addPiece() // Replace with random piece
		return true
	}
	return false
}

// instafall calls the applyGravity function until a collision is detected.
func (b *Board) instafall() {
	collide := false
	for !collide {
		collide = b.applyGravity()
	}
}

// checkRowCompletion checks if the rows in a given shape are filled (ie should
// be deleted). If full, deletes the rows.
func (b *Board) checkRowCompletion(s Shape) {
	// Ony the rows of the shape can be filled
	rowWasDeleted := true
	// Since when we delete a row it can be shifted down, repeatedly try
	// to delete a row until no more deletes can be made
	var deleteRowCt int
	for rowWasDeleted {
		rowWasDeleted = false
		for i := 0; i < 4; i++ {
			r := s[i].row
			emptyFound := false
			// Look for empty row
			for c := 0; c < 10; c++ {
				if b[r][c] == Empty {
					emptyFound = true
					continue
				}
			}
			// If no empty cell was found in row delete row
			if !emptyFound {
				b.deleteRow(r)
				rowWasDeleted = true
				score += 200
				deleteRowCt++
			}
		}
	}
	// Bonus score for combos over one
	if deleteRowCt > 1 {
		score += (deleteRowCt - 1) * 200
	}
}

// deleteRow remoes a row by shifting everything above it down by one.
func (b *Board) deleteRow(row int) {
	for r := row; r < 21; r++ {
		for c := 0; c < 10; c++ {
			b[r][c] = b[r+1][c]
		}
	}
}

// setPiece sets a value in the game board to a specific block type.
func (b *Board) setPiece(r, c int, val Block) {
	b[r][c] = val
}

// fillShape sets
func (b *Board) fillShape(s Shape, val Block) {
	for i := 0; i < 4; i++ {
		b.setPiece(s[i].row, s[i].col, val)
	}
}

// addPiece creates a piece at the top of the screen at a random position
// and sets it to the piece that the player is controlling
// (ie activeShape).
func (b *Board) addPiece() {
	var offset int
	if nextPiece == IPiece {
		offset = rand.Intn(7)
	} else if nextPiece == OPiece {
		offset = rand.Intn(9)
	} else {
		offset = rand.Intn(8)
	}
	baseShape := getShapeFromPiece(nextPiece)
	baseShape = moveShape(20, offset, baseShape)
	b.fillShape(baseShape, piece2Block(nextPiece))
	currentPiece = nextPiece
	activeShape = baseShape
	nextPiece = Piece(rand.Intn(7))
}

// displayBoard displays a particular game board with all of its pieces
// onto a given window, win
func (b *Board) displayBoard(win *pixelgl.Window) {
	boardBlockSize := 20.0 //win.Bounds().Max.X / 10
	pic := blockGen(0)
	imgSize := pic.Bounds().Max.X
	scaleFactor := float64(boardBlockSize) / float64(imgSize)

	for col := 0; col < BoardCols; col++ {
		for row := 0; row < BoardRows-2; row++ {
			val := b[row][col]
			if val == Empty {
				continue
			}

			x := float64(col)*boardBlockSize + boardBlockSize/2
			y := float64(row)*boardBlockSize + boardBlockSize/2
			pic := blockGen(block2spriteIdx(val))
			sprite := pixel.NewSprite(pic, pic.Bounds())
			sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, scaleFactor).Moved(pixel.V(x+282, y+25)))
		}
	}

	// Display Shadow
	pieceType := b[activeShape[0].row][activeShape[0].col]
	ghostShape := activeShape
	b.drawPiece(activeShape, Empty)
	for {
		if b.checkCollision(moveShapeDown(ghostShape)) {
			break
		}
		ghostShape = moveShapeDown(ghostShape)
	}
	b.drawPiece(activeShape, pieceType)

	gpic := blockGen(block2spriteIdx(Gray))
	sprite := pixel.NewSprite(gpic, gpic.Bounds())
	for i := 0; i < 4; i++ {
		if b[ghostShape[i].row][ghostShape[i].col] == Empty {
			x := float64(ghostShape[i].col)*boardBlockSize + boardBlockSize/2
			y := float64(ghostShape[i].row)*boardBlockSize + boardBlockSize/2
			sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, scaleFactor/2).Moved(pixel.V(x+282, y+25)))
		}
	}
}
