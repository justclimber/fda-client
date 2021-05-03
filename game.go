package main

import (
	"errors"
	"github.com/hajimehoshi/ebiten/v2"
	"time"
)

type SceneID string

type Scene interface {
	Draw(screen *ebiten.Image)
	Update(dt time.Duration) (*SceneID, error)
}

const (
	GameSceneStart = "start"
	GameSceneMain  = "main"
)

type Game struct {
	CurrentScene Scene
	scenes       map[SceneID]Scene
	lastTime     time.Time
}

func (g *Game) Update() error {
	now := time.Now()
	dt := now.Sub(g.lastTime)
	g.lastTime = now
	switchToScene, err := g.CurrentScene.Update(dt)
	if switchToScene != nil {
		g.SwitchScene(*switchToScene)
	}
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
	}
	g.LoadScenes()
	if err := g.SwitchScene(GameSceneStart); err != nil {
		panic(err)
	}
	return g
}

func (g *Game) LoadScenes() {
	g.scenes = make(map[SceneID]Scene)
	g.scenes[GameSceneStart] = newSceneStart()
	g.scenes[GameSceneMain] = newSceneMain()
}

func (g *Game) SwitchScene(s SceneID) error {
	scene, ok := g.scenes[s]
	if !ok {
		return errors.New("Undefined or unloaded scene: " + string(s))
	}
	g.CurrentScene = scene
	return nil
}
