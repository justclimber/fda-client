package main

import (
	"encoding/json"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/justclimber/fda-client/config"
	_ "image/png"
	"io/ioutil"
	"strings"
	"time"
)

const configFileName = "config.json"
const imgsDirPath = "assets/images/"

type SceneStart struct {
	SceneState
	assets *config.Assets
	config *config.Config
}

func newSceneStart(assets *config.Assets, config *config.Config) *SceneStart {
	s := &SceneStart{
		assets: assets,
		config: config,
	}
	s.stateUpdateFunc = s.loadConfigUpdate()
	return s
}

func (s *SceneStart) loadConfigUpdate() SceneStateUpdateFunc {
	return func(dt time.Duration) (SceneStateUpdateFunc, SceneStateDrawFunc, bool, error) {
		jsonfile, err := ioutil.ReadFile(configFileName)
		if err != nil {
			return s.error("loading config file", err)
		}
		var c config.Config
		err = json.Unmarshal(jsonfile, &c)
		if err != nil {
			return s.error("decoding json config file", err)
		}
		s.config = &c
		return s.loadAssetsUpdate(0, 0, []string{}), s.loadDraw([]string{"config loaded"}), false, nil
	}
}

func (s *SceneStart) loadAssetsUpdate(imgsLoaded int, fontsLoaded int, log []string) SceneStateUpdateFunc {
	return func(dt time.Duration) (SceneStateUpdateFunc, SceneStateDrawFunc, bool, error) {
		if imgsLoaded == len(config.GetAvailableImages()) {
			duration := 3 * time.Second
			return s.serverConnectingUpdate(duration), s.serverConnectingDraw(duration), true, nil
		}
		if imgsLoaded < len(config.GetAvailableImages()) {
			imgToLoad := config.GetAvailableImages()[imgsLoaded]
			filePath := imgsDirPath + string(imgToLoad) + ".png"
			// @todo: get width and height from config
			img, err := loadImageNineSlice(filePath, 5, 5)
			if err != nil {
				return s.error("loading image", err)
			}
			s.assets.NineSlices[imgToLoad] = img
			log = append(log, string(imgToLoad) + " loaded")
			imgsLoaded++
		}
		return s.loadAssetsUpdate(imgsLoaded, fontsLoaded, log), s.loadDraw(log), false, nil
	}
}

func (s *SceneStart) loadDraw(log []string) SceneStateDrawFunc {
	return func(screen *ebiten.Image) {
		ebitenutil.DebugPrint(screen, strings.Join(log, "\n"))
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
