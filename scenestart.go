package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"time"
)

type SceneStart struct {
	SceneState
}

func newSceneStart() *SceneStart {
	s := &SceneStart{}
	s.stateUpdateFunc = s.loadConfigUpdate(3 * time.Second)
	return s
}

func (s *SceneStart) loadConfigUpdate(elapsed time.Duration) SceneStateUpdateFunc {
	return func(dt time.Duration) (SceneStateUpdateFunc, SceneStateDrawFunc, bool, error) {
		elapsed -= dt
		if elapsed <= 0 {
			timeToConnect := 3 * time.Second
			return s.serverConnectingUpdate(timeToConnect), s.serverConnectingDraw(timeToConnect), true, nil
		}
		return s.loadConfigUpdate(elapsed), s.loadConfigDraw(elapsed), false, nil
	}
}

func (s *SceneStart) loadConfigDraw(elapsed time.Duration) SceneStateDrawFunc {
	return func(screen *ebiten.Image) {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Config loading... %.1f sec left", elapsed.Seconds()))
	}
}

func (s *SceneStart) serverConnectingUpdate(elapsed time.Duration) SceneStateUpdateFunc {
	return func(dt time.Duration) (SceneStateUpdateFunc, SceneStateDrawFunc, bool, error) {
		elapsed -= dt
		if elapsed <= 0 {
			return s.SwitchToScene(GameSceneMain)
		}
		return s.serverConnectingUpdate(elapsed), s.serverConnectingDraw(elapsed), false, nil
	}
}

func (s *SceneStart) serverConnectingDraw(elapsed time.Duration) SceneStateDrawFunc {
	return func(screen *ebiten.Image) {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Server connecting... %.1f sec left", elapsed.Seconds()))
	}
}
