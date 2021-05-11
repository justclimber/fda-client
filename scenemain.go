package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/justclimber/ebitenui"
	"github.com/justclimber/ebitenui/widget"
	"golang.org/x/image/colornames"
	"image"
	"sync"
	"time"
)

type SceneMain struct {
	SceneState
	g                  *Game
	ui                 *ebitenui.UI
	init               sync.Once
	historyPlayerImage *ebiten.Image
}

func newSceneMain(g *Game) *SceneMain {
	s := &SceneMain{g: g}
	s.stateUpdateFunc = s.idleUpdate
	s.stateDrawFunc = s.idleDraw
	s.historyPlayerImage = ebiten.NewImage(200, 200)

	return s
}

func (s *SceneMain) OnSwitch() error {
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

func (s *SceneMain) drawHistoryPlayerCallback(image *ebiten.Image, origin image.Rectangle)  {
	ebitenutil.DrawRect(
		image,
		float64(30 + origin.Min.X),
		float64(30 + origin.Min.Y),
		100,
		100,
		colornames.Aqua,
	)
}
