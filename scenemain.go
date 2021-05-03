package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/justclimber/ebitenui"
	"github.com/justclimber/fda-client/config"
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
	s.stateUpdateFunc = s.idleUpdate
	s.stateDrawFunc = s.idleDraw
	return s
}

func (s *SceneMain) idleUpdate(dt time.Duration) (SceneStateUpdateFunc, SceneStateDrawFunc, bool, error) {
	return nil, nil, false, nil
}

func (s *SceneMain) idleDraw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Yay! we in the game now!")
}
