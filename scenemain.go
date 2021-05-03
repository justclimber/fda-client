package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/justclimber/ebitenui"
	"github.com/justclimber/fda-client/config"
	"log"
	"time"
)

type SceneMain struct {
	SceneState
	assets *config.Assets
	config *config.Config
	ui *ebitenui.UI
}

func newSceneMain(assets *config.Assets, config *config.Config) *SceneMain {
	s := &SceneMain{
		assets: assets,
		config: config,
	}
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
