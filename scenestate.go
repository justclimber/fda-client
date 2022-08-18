package main

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
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

func (s *SceneState) runState(dt time.Duration) error {
	newStateUpdate, newStateDraw, rerun, err := s.stateUpdateFunc(dt)
	if err != nil {
		return err
	}
	if newStateUpdate == nil {
		return nil
	}
	s.stateUpdateFunc = newStateUpdate
	s.stateDrawFunc = newStateDraw
	if rerun {
		return s.runState(dt)
	}
	return nil
}

func (s *SceneState) Draw(screen *ebiten.Image) {
	s.stateDrawFunc(screen)
}

func (s *SceneState) Update(dt time.Duration) error {
	return s.runState(dt)
}

func (s *SceneState) error(msg string, err error) (SceneStateUpdateFunc, SceneStateDrawFunc, bool, error) {
	return nil, nil, false, fmt.Errorf("%s: %w", msg, err)
}

func (s *SceneState) staySameState() (SceneStateUpdateFunc, SceneStateDrawFunc, bool, error) {
	return nil, nil, false, nil
}
