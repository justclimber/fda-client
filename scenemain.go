package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/justclimber/ebitenui"
	"log"
	"time"
)

type SceneMain struct {
	SceneState
	g *Game
	ui *ebitenui.UI
}

func newSceneMain(g *Game) *SceneMain {
	s := &SceneMain{g: g}
	s.stateUpdateFunc = s.setupUpdate
	s.stateDrawFunc = s.idleDraw

	return s
}

func (s *SceneMain) setupUpdate(dt time.Duration) (SceneStateUpdateFunc, SceneStateDrawFunc, bool, error) {
	s.setupUI()
	log.Print(123123)
	return s.idleUpdate, s.idleDraw, true, nil
}

func (s *SceneMain) idleUpdate(dt time.Duration) (SceneStateUpdateFunc, SceneStateDrawFunc, bool, error) {
	s.ui.Update()
	return s.staySameState()
}

func (s *SceneMain) idleDraw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Yay! we in the game now!")
	if s.ui != nil {
		s.ui.Draw(screen)
	}
}
