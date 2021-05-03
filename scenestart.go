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
	s.stateUpdateFunc = s.stateUpdateLoadConfig(3 * time.Second)
	return s
}

func (s *SceneStart) stateUpdateLoadConfig(elapsed time.Duration) SceneStateUpdateFunc {
	return func(dt time.Duration) (SceneStateUpdateFunc, SceneStateDrawFunc, bool, error) {
		elapsed -= dt
		if elapsed <= 0 {
			timeToConnect := 3 * time.Second
			return s.stateUpdateServerConnecting(timeToConnect), s.stateDrawServerConnecting(timeToConnect), true, nil
		}
		return s.stateUpdateLoadConfig(elapsed), s.stateDrawLoadConfig(elapsed), false, nil
	}
}

func (s *SceneStart) stateDrawLoadConfig(elapsed time.Duration) SceneStateDrawFunc {
	return func(screen *ebiten.Image) {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Config loading... %.1f sec left", elapsed.Seconds()))
	}
}

func (s *SceneStart) stateUpdateServerConnecting(elapsed time.Duration) SceneStateUpdateFunc {
	return func(dt time.Duration) (SceneStateUpdateFunc, SceneStateDrawFunc, bool, error) {
		elapsed -= dt
		if elapsed <= 0 {
			return s.SwitchToScene(GameSceneMain)
		}
		return s.stateUpdateServerConnecting(elapsed), s.stateDrawServerConnecting(elapsed), false, nil
	}
}

func (s *SceneStart) stateDrawServerConnecting(elapsed time.Duration) SceneStateDrawFunc {
	return func(screen *ebiten.Image) {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Server connecting... %.1f sec left", elapsed.Seconds()))
	}
}
