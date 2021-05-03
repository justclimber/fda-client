package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"time"
)

type SceneStateUpdateFunc func(dt time.Duration) (
	nextStateUpdate SceneStateUpdateFunc,
	nextStateDraw SceneStateDrawFunc,
	rerun bool,
	err error,
)
type SceneStateDrawFunc func(screen *ebiten.Image)

type SceneState struct {
	stateUpdateFunc SceneStateUpdateFunc
	stateDrawFunc   SceneStateDrawFunc
	switchToScene   *SceneID
}

func (s *SceneState) runState(dt time.Duration) (*SceneID, error) {
	newStateUpdate, newStateDraw, rerun, err := s.stateUpdateFunc(dt)
	if err != nil {
		return nil, err
	}
	if newStateUpdate == nil {
		if s.switchToScene != nil {
			return s.switchToScene, nil
		}
		return nil, nil
	}
	s.stateUpdateFunc = newStateUpdate
	s.stateDrawFunc = newStateDraw
	if rerun {
		return s.runState(dt)
	}
	return nil, nil
}

func (s *SceneState) Draw(screen *ebiten.Image) {
	s.stateDrawFunc(screen)
}

func (s *SceneState) Update(dt time.Duration) (*SceneID, error) {
	return s.runState(dt)
}

func (s *SceneState) SwitchToScene(newScene SceneID) (
	nextStateUpdate SceneStateUpdateFunc,
	nextStateDraw SceneStateDrawFunc,
	rerun bool,
	err error,
) {
	s.switchToScene = &newScene
	return nil, nil, false, nil
}
