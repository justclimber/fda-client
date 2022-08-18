package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func RunOnEbiten(game *Game) {
	ebiten.SetWindowSize(1300, 800)
	ebiten.SetWindowTitle("FDA Game Prototype")
	ebiten.SetVsyncEnabled(false)
	ebiten.SetWindowResizable(true)
	//ebiten.SetMaxTPS(1)

	err := ebiten.RunGame(game)
	if err != nil {
		log.Print(err)
	}
}
