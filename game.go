package main

import (
	"errors"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/justclimber/fda-client/config"
)

type SceneID string

type Scene interface {
	Draw(screen *ebiten.Image)
	Update(dt time.Duration) error
	OnSwitch() error
}

const (
	GameSceneStart = "start"
	GameSceneMain  = "main"
)

type Game struct {
	CurrentScene Scene
	scenes       map[SceneID]Scene
	lastTime     time.Time
	assets       *config.Assets
	config       *config.Config
}

func (g *Game) Update() error {
	now := time.Now()
	dt := now.Sub(g.lastTime)
	g.lastTime = now
	err := g.CurrentScene.Update(dt)
	return err
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.CurrentScene.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func NewGame() *Game {
	g := &Game{
		lastTime: time.Now(),
		assets:   config.NewAssets(),
	}
	g.LoadScenes()
	if err := g.SwitchScene(GameSceneStart); err != nil {
		panic(err)
	}
	return g
}

func (g *Game) LoadScenes() {
	g.scenes = make(map[SceneID]Scene)
	g.scenes[GameSceneStart] = newSceneStart(g)
	g.scenes[GameSceneMain] = newSceneMain(g)
}

func (g *Game) SwitchScene(s SceneID) error {
	scene, ok := g.scenes[s]
	if !ok {
		return errors.New("Undefined or unloaded scene: " + string(s))
	}
	err := scene.OnSwitch()
	if err != nil {
		return err
	}
	g.CurrentScene = scene
	return nil
}
