package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func RunOnEbiten(game *Game)  {
	ebiten.SetWindowSize(1300, 800)
	ebiten.SetWindowTitle("FDA Game Prototype")
	ebiten.SetVsyncEnabled(false)
	ebiten.SetWindowResizable(true)


	err := ebiten.RunGame(game)
	if err != nil {
		log.Print(err)
	}
}
