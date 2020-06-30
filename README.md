# Golang Tetris

This is my first experience writing games in Golang. I stumbled upon the 2D game library at [github.com/faiface/pixel](github.com/faiface/pixel) and wanted to try it out. After using it for this game, I have to admit that the library is very well designed and easy to use. Anyone interested in creating games should definitely check it out.

This is a typical Tetris clone. I tried to make the gameplay experience as smooth as possible though the game is lacking bells and whistles like a title-screen. Here are some screen shots of the game:

![A sample example of the program](docs/media/example1.png?raw=true "An example of the program running")
![A sample example of the program](docs/media/example2.png?raw=true "An example of the program running")

## Running the Game

Having Go installed, you can run `go run .` from the root directory to play the game.

## Controls

- Left/Right arrow - Move piece
- Up arrow - Rotate piece
- Down arrow - Fast fall
- Space - Instant drop

## Todo

- [ ] Menus (Opening, game-over, pause)
- [ ] Animation for row clearing
- [ ] Music and sound effects
