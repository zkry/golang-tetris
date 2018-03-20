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

	ss "github.com/zkry/golang-tetris/spritesheet"
)

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

func main() {
	pixelgl.Run(run)
}

// run is the main code for the game. Allows pixelgl to run on main thread
func run() {
	// Initialize the window
	windowWidth := 765.0
	windowHeight := 450.0
	cfg := pixelgl.WindowConfig{
		Title:  "Blockfall",
		Bounds: pixel.R(0, 0, windowWidth, windowHeight),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Load Various Resources:
	// Matriax on opengameart.org
	blockGen, err = ss.LoadSpriteSheet("resources/blocks.png", 2, 8)
	if err != nil {
		panic(err)
	}

	// Background image, by ansimuz on opengameart.org
	bgPic, err := ss.LoadPicture("resources/parallax-mountain-bg.png")
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
	gameBoard.addPiece() // Add initial Piece to game
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

		// Delay for left/right movement. When a key is pressed and holded
		// it should first move when pressed, then after a short wait,
		// it will continuously move. Like when a key is pressed in a text editor
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

		// Keypress Functions
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
			gravitySpeed = 0.08 // TODO: Code could result in bugs if game pause functionality added
			if gravityTimer > 0.08 {
				gravityTimer = 0.08
			}
		}
		if win.JustReleased(pixelgl.KeyDown) {
			gravitySpeed = baseSpeed // TODO: Code could result in bugs if game pause functionality added
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

		// Display Functions
		win.Clear(colornames.Black)
		displayBG(win)
		displayText(win)
		gameBoard.displayBoard(win)
		win.Update()
	}
}

func displayText(win *pixelgl.Window) {
	// Text Generator
	scoreTextLocX := 500.0
	scoreTextLocY := 400.0
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	scoreTxt := text.New(pixel.V(scoreTextLocX, scoreTextLocY), basicAtlas)
	fmt.Fprintf(scoreTxt, "Score: %d", score)
	scoreTxt.Draw(win, pixel.IM.Scaled(scoreTxt.Orig, 2))

	nextPieceTextLocX := 142.0
	nextPieceTextLocY := 285.0
	nextPieceTxt := text.New(pixel.V(nextPieceTextLocX, nextPieceTextLocY), basicAtlas)
	fmt.Fprintf(nextPieceTxt, "Next Piece:")
	nextPieceTxt.Draw(win, pixel.IM)
}

func displayBG(win *pixelgl.Window) {
	// Display various background images
	bgImgSprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
	gameBGSprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
	nextPieceBGSprite.Draw(win, pixel.IM.Moved(pixel.V(182, 225)))

	// Display next block
	baseShape := getShapeFromPiece(nextPiece)
	pic := blockGen(block2spriteIdx(piece2Block(nextPiece)))
	sprite := pixel.NewSprite(pic, pic.Bounds())
	boardBlockSize := 20.0
	scaleFactor := float64(boardBlockSize) / pic.Bounds().Max.Y
	shapeWidth := getShapeWidth(baseShape) + 1
	shapeHeight := 2

	for i := 0; i < 4; i++ {
		r := baseShape[i].row
		c := baseShape[i].col
		x := float64(c)*boardBlockSize + boardBlockSize/2
		y := float64(r)*boardBlockSize + boardBlockSize/2
		sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, scaleFactor).Moved(pixel.V(x+182-(float64(shapeWidth)*10), y+225-(float64(shapeHeight)*10))))
	}
}

// block2spriteIdx associates a blocks color (b Block) with its index in the sprite sheet.
func block2spriteIdx(b Block) int {
	return int(b) - 1
}

// piece2Block associates a pieces shape (Piece) with it's color/image (Block).
func piece2Block(p Piece) Block {
	switch p {
	case LPiece:
		return Goluboy
	case IPiece:
		return Siniy
	case OPiece:
		return Pink
	case TPiece:
		return Purple
	case SPiece:
		return Red
	case ZPiece:
		return Yellow
	case JPiece:
		return Green
	}
	panic("piece2Block: Invalid piece passed in")
	return GraySpecial // Return strange value value
}
