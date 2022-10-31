package main

import (
	"log"

	"github.com/dragun-igor/hero_knight/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowTitle("Hero Knight")
	screenWidth := 1600
	screenHeight := 247
	ebiten.SetWindowSize(screenWidth, screenHeight)
	game, err := game.NewGame()
	if err != nil {
		log.Fatal(err)
	}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
