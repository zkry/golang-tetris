package main

import (
	"fmt"
	_ "image/png"
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"

	ss "github.com/zkry/blockfall/spritesheet"
)

func main() {
	pixelgl.Run(run)
}

// BoardRows is the height of the game board in terms of blocks
const BoardRows = 22

// BoardCols is the width of the game board in terms of blocks
const BoardCols = 10

// Point represents a coordinate on the game board with Point{row:0, col:0}
// representing the bottom left
type Point struct {
	row int
	col int
}

// Board is an array containing the entire game board pieces.
type Board [22][10]Block

// Block represents the color of the block
type Block int

// Different values a point on the grid can hold
const (
	Empty Block = iota
	Goluboy
	Siniy
	Pink
	Purple
	Red
	Yellow
	Green
	Gray
	GoluboySpecial
	SiniySpecial
	PinkSpecial
	PurpleSpecial
	RedSpecial
	YellowSpecial
	GreenSpecial
	GraySpecial
)

// Piece is a constant for a shape of piece. There are 7 classic pieces like L, and O
type Piece int

// Various values that the pieces can be
const (
	IPiece Piece = iota
	JPiece
	LPiece
	OPiece
	SPiece
	TPiece
	ZPiece
)

// Shape is a type containing four points, which represents the four points
// making a contiguous 'piece'.
type Shape [4]Point

const levelLength = 60.0 // Time it takes for game to speed up
const speedUpRate = 0.1  // Every new level, the amount the game speeds up by

var gameBoard Board
var activeShape Shape // The shape that the player controls
var currentPiece Piece
var gravityTimer float64
var baseSpeed float64 = 0.8
var gravitySpeed float64 = 0.8
var levelUpTimer float64 = levelLength
var gameOver bool = false
var leftRightDelay float64
var moveCounter int
var score int
var nextPiece Piece

var blockGen func(int) pixel.Picture
var bgImgSprite pixel.Sprite
var gameBGSprite pixel.Sprite
var nextPieceBGSprite pixel.Sprite

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Blockfall",
		Bounds: pixel.R(0, 0, 765, 450),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Matriax on opengameart.org
	blockGen, err = ss.LoadSpriteSheet("blocks.png", 2, 8)
	if err != nil {
		panic(err)
	}

	// Background image, by ansimuz on opengameart.org
	bgPic, err := ss.LoadPicture("parallax-mountain-bg.png")
	if err != nil {
		panic(err)
	}
	bgImgSprite = *pixel.NewSprite(bgPic, bgPic.Bounds())

	// Game Background
	blackPic := ss.GetPlayBGPic()
	gameBGSprite = *pixel.NewSprite(blackPic, blackPic.Bounds())

	// Next Piece BG
	nextPiecePic := ss.GetNextPieceBGPic()
	nextPieceBGSprite = *pixel.NewSprite(nextPiecePic, nextPiecePic.Bounds())

	nextPiece = Piece(rand.Intn(7))
	gameBoard.addPiece()
	last := time.Now()
	for !win.Closed() && !gameOver {
		// Perform time processing events
		dt := time.Since(last).Seconds()
		last = time.Now()
		gravityTimer += dt
		levelUpTimer -= dt

		// Time Functions:
		// Gravity
		if gravityTimer > gravitySpeed {
			gravityTimer -= gravitySpeed
			didCollide := gameBoard.applyGravity()
			if !didCollide {
				if gameBoard.isTouchingFloor() {
					gravityTimer -= gravitySpeed // Add extra time when touching floor
				}
			} else {
				score += 10
			}
		}

		// Delay for left/right movement
		if leftRightDelay > 0.0 {
			leftRightDelay = math.Max(leftRightDelay-dt, 0.0)
		}

		// Speed up
		if levelUpTimer <= 0 {
			if baseSpeed > 0.2 {
				baseSpeed = math.Max(baseSpeed-speedUpRate, 0.2)
			}
			levelUpTimer = levelLength
			gravitySpeed = baseSpeed
		}

		if win.Pressed(pixelgl.KeyRight) && leftRightDelay == 0.0 {
			gameBoard.movePiece(1)
			if moveCounter > 0 {
				leftRightDelay = 0.1
			} else {
				leftRightDelay = 0.5
			}
			moveCounter++
		}
		if win.Pressed(pixelgl.KeyLeft) && leftRightDelay == 0.0 {
			gameBoard.movePiece(-1)
			if moveCounter > 0 {
				leftRightDelay = 0.1
			} else {
				leftRightDelay = 0.5
			}
			moveCounter++
		}
		if win.JustPressed(pixelgl.KeyDown) {
			gravitySpeed = 0.08 // TODO: Code could result in bugs
			if gravityTimer > 0.08 {
				gravityTimer = 0.08
			}
		}
		if win.JustReleased(pixelgl.KeyDown) {
			gravitySpeed = baseSpeed // TODO: Code could result in bugs
		}
		if win.JustPressed(pixelgl.KeyUp) {
			gameBoard.rotatePiece()
			if gameBoard.isTouchingFloor() {
				gravityTimer = 0 // Make gravity more forgiving when moving pieces
			}
		}
		if win.JustPressed(pixelgl.KeySpace) {
			gameBoard.instafall()
			score += 12
		}
		if !win.Pressed(pixelgl.KeyRight) && !win.Pressed(pixelgl.KeyLeft) {
			moveCounter = 0
			leftRightDelay = 0
		}

		win.Clear(colornames.Black)
		displayBG(win)

		displayScore(win)
		gameBoard.displayBoard(win)
		win.Update()
	}
}

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

func (b *Board) isTouchingFloor() bool {
	blockType := b[activeShape[0].row][activeShape[0].col]
	b.drawPiece(activeShape, Empty)
	isTouching := b.checkCollision(moveShapeDown(activeShape))
	b.drawPiece(activeShape, blockType)
	return isTouching
}

func (b *Board) rotatePiece() {
	if currentPiece == OPiece {
		return
	}
	blockType := b[activeShape[0].row][activeShape[0].col]
	// Erase Piece
	b.drawPiece(activeShape, Empty)

	newShape := rotateShape(activeShape)
	if b.checkCollision(newShape) {
		if !b.checkCollision(moveShapeRight(newShape)) {
			newShape = moveShapeRight(newShape)
		} else if !b.checkCollision(moveShapeLeft(newShape)) {
			newShape = moveShapeLeft(newShape)
			// TODO: Add up case
		} else {
			b.drawPiece(activeShape, blockType)
			return
		}
	}
	activeShape = newShape
	b.drawPiece(activeShape, blockType)
}

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

func (b *Board) drawPiece(s Shape, t Block) {
	for i := 0; i < 4; i++ {
		b[activeShape[i].row][activeShape[i].col] = t
	}
}

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

func (b *Board) applyGravity() bool {
	blockType := b[activeShape[0].row][activeShape[0].col]
	// Erase old piece
	b.drawPiece(activeShape, Empty)

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

func (b *Board) instafall() {
	collide := false
	for !collide {
		collide = b.applyGravity()
	}
}

func (b *Board) checkRowCompletion(s Shape) {
	// Ony the rows of the shape can be filled
	rowWasDeleted := true
	// Since when we delete a row it can be shifted down, repeatedly try
	// to delete a row until no more deletes can be made
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
			}
		}
	}
}

func (b *Board) deleteRow(row int) {
	for r := row; r < 21; r++ {
		for c := 0; c < 10; c++ {
			b[r][c] = b[r+1][c]
		}
	}
}

func (b *Board) setPiece(r, c int, val Block) {
	b[r][c] = val
}

func (b *Board) fillShape(s Shape, val Block) {
	for i := 0; i < 4; i++ {
		b.setPiece(s[i].row, s[i].col, val)
	}
}

func (b *Board) addPiece() {
	var s Shape
	switch t := nextPiece; t {
	case LPiece:
		c := rand.Intn(8)
		s = Shape{
			Point{row: 21, col: c},
			Point{row: 21, col: c + 1},
			Point{row: 21, col: c + 2},
			Point{row: 20, col: c},
		}
		b.fillShape(s, Goluboy)
		currentPiece = LPiece
	case IPiece:
		c := rand.Intn(7)
		s = Shape{
			Point{row: 21, col: c},
			Point{row: 21, col: c + 1},
			Point{row: 21, col: c + 2},
			Point{row: 21, col: c + 3},
		}
		b.fillShape(s, Siniy)
		currentPiece = IPiece
	case OPiece:
		c := rand.Intn(9)
		s = Shape{
			Point{row: 21, col: c},
			Point{row: 21, col: c + 1},
			Point{row: 20, col: c},
			Point{row: 20, col: c + 1},
		}
		b.fillShape(s, Pink)
		currentPiece = OPiece
	case TPiece:
		c := rand.Intn(8)
		s = Shape{
			Point{row: 21, col: c},
			Point{row: 21, col: c + 1},
			Point{row: 21, col: c + 2},
			Point{row: 20, col: c + 1},
		}
		b.fillShape(s, Purple)
		currentPiece = TPiece
	case SPiece:
		c := rand.Intn(8)
		s = Shape{
			Point{row: 20, col: c},
			Point{row: 20, col: c + 1},
			Point{row: 21, col: c + 1},
			Point{row: 21, col: c + 2},
		}
		b.fillShape(s, Red)
		currentPiece = SPiece
	case ZPiece:
		c := rand.Intn(8)
		s = Shape{
			Point{row: 21, col: c},
			Point{row: 21, col: c + 1},
			Point{row: 20, col: c + 1},
			Point{row: 20, col: c + 2},
		}
		b.fillShape(s, Yellow)
		currentPiece = ZPiece
	case JPiece:
		c := rand.Intn(8)
		s = Shape{
			Point{row: 21, col: c},
			Point{row: 20, col: c + 1},
			Point{row: 20, col: c},
			Point{row: 20, col: c + 2},
		}
		b.fillShape(s, Green)
		currentPiece = JPiece
	default:
		panic("addPiece(): Invalid piece entered")
	}
	activeShape = s
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

func displayScore(win *pixelgl.Window) {
	// Text Generator
	textLocX := 500.0
	textLocY := 400.0
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	scoreTxt := text.New(pixel.V(textLocX, textLocY), basicAtlas)
	fmt.Fprintf(scoreTxt, "Score: %d", score)
	scoreTxt.Draw(win, pixel.IM)
}

func displayBG(win *pixelgl.Window) {
	bgImgSprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
	gameBGSprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
	nextPieceBGSprite.Draw(win, pixel.IM.Moved(pixel.V(182, 225)))
}

func block2spriteIdx(b Block) int {
	return int(b) - 1
}
