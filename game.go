package main

import (
	"errors"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justclimber/fda-client/config"
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
	assets       *config.Assets
	config       *config.Config
}

func (g *Game) Update() error {
	now := time.Now()
	dt := now.Sub(g.lastTime)
	g.lastTime = now
	switchToScene, err := g.CurrentScene.Update(dt)
	if switchToScene != nil {
		err = g.SwitchScene(*switchToScene)
		if err != nil {
			return err
		}
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
		assets: config.NewAssets(),
	}
	g.LoadScenes()
	if err := g.SwitchScene(GameSceneStart); err != nil {
		panic(err)
	}
	return g
}

func (g *Game) LoadScenes() {
	g.scenes = make(map[SceneID]Scene)
	g.scenes[GameSceneStart] = newSceneStart(g.assets, g.config)
}

func (g *Game) SwitchScene(s SceneID) error {
	scene, loaded := g.scenes[s]
	if !loaded {
		switch s {
		case GameSceneMain:
			// @fixme: ugly hack!
			g.config = g.scenes[GameSceneStart].(*SceneStart).config
			scene = newSceneMain(g.assets, g.config)
			g.scenes[GameSceneMain] = scene
		default:
			return errors.New("Undefined or unloaded scene: " + string(s))
		}
	}
	g.CurrentScene = scene
	return nil
}
