package spritesheet

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/faiface/pixel"
)

func LoadSpriteSheet(path string, row, col int) (func(int) pixel.Picture, error) {
	// Open file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Load Image
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	// Check if tile is square
	b := img.Bounds()
	if b.Max.X/col != b.Max.Y/row {
		fmt.Println("width/col = ", b.Max.X, ", height/row = ", b.Max.Y)
		return nil, errors.New(fmt.Sprintf("Invalid dimensions (%d, %d) for sprite sheet %s\n", row, col, path))
	}

	tileSize := b.Max.X / col

	return func(i int) pixel.Picture {
		if i < 0 || i >= row*col {
			panic("Index out of bounds for sprite sheet")
		}
		r := i / col
		c := i % col

		subImage := img.(interface {
			SubImage(r image.Rectangle) image.Image
		}).SubImage(image.Rect(c*tileSize, r*tileSize, (c+1)*tileSize, (r+1)*tileSize))
		return pixel.PictureDataFromImage(subImage)
	}, nil
}

func LoadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func GetPlayBGPic() pixel.Picture {
	blackImg := image.NewRGBA(image.Rect(0, 0, 200, 400))
	for x := 0; x < 200; x++ {
		for y := 0; y < 400; y++ {
			blackImg.SetRGBA(x, y, color.RGBA{0x00, 0x00, 0x00, 0xA0})
		}
	}

	blackPic := pixel.PictureDataFromImage(blackImg)
	return blackPic
}

func GetNextPieceBGPic() pixel.Picture {
	blackImg := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			blackImg.SetRGBA(x, y, color.RGBA{0x00, 0x00, 0x00, 0xA0})
		}
	}
	blackPic := pixel.PictureDataFromImage(blackImg)
	return blackPic
}
