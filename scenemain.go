package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"time"
)

type SceneMain struct {
	SceneState
}

func newSceneMain() *SceneMain {
	s := &SceneMain{}
	s.stateUpdateFunc = s.stateUpdateIdle
	s.stateDrawFunc = s.stateDrawIdle
	return s
}

func (s *SceneMain) stateUpdateIdle(dt time.Duration) (SceneStateUpdateFunc, SceneStateDrawFunc, bool, error) {
	return nil, nil, false, nil
}

func (s *SceneMain) stateDrawIdle(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Yay! we in the game now!")
}
