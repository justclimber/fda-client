package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/justclimber/ebitenui"
	"github.com/justclimber/ebitenui/widget"
	"sync"
	"time"
)

type SceneMain struct {
	SceneState
	g *Game
	ui *ebitenui.UI
	init sync.Once
}

func newSceneMain(g *Game) *SceneMain {
	s := &SceneMain{g: g}
	s.stateUpdateFunc = s.idleUpdate
	s.stateDrawFunc = s.idleDraw

	return s
}

func (s *SceneMain) Setup() error {
	var err error
	s.init.Do(func() {
		err = s.setupUI()
	})
	return err
}

func (s *SceneMain) idleUpdate(dt time.Duration) (SceneStateUpdateFunc, SceneStateDrawFunc, bool, error) {
	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		s.ui.SetDebugMode(widget.DebugModeNone)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
		s.ui.SetDebugMode(widget.DebugModeBorderOnMouseOver)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		s.ui.SetDebugMode(widget.DebugModeBorderAlwaysShow)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF4) {
		s.ui.SetDebugMode(widget.DebugModeInputLayersAlwaysShow)
	}
	s.ui.Update()
	return s.staySameState()
}

func (s *SceneMain) idleDraw(screen *ebiten.Image) {
	s.ui.Draw(screen)
}
