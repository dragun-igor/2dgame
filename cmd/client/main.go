package main

import (
	"log"

	"github.com/dragun-igor/hero_knight/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowTitle("Hero Knight")
	screenWidth := 640
	screenHeight := 320
	ebiten.SetWindowSize(screenWidth, screenHeight)
	game := game.NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
