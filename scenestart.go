package main

import (
	"encoding/json"
	"errors"
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
const fontsDirPath = "assets/fonts/"

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

func (s *SceneStart) loadAssetsUpdate(imgIndex int, fontIndex int, log []string) SceneStateUpdateFunc {
	return func(dt time.Duration) (SceneStateUpdateFunc, SceneStateDrawFunc, bool, error) {
		if imgIndex == len(config.GetAvailableImages()) && fontIndex == len(config.GetAvailableFonts()) {
			duration := 3 * time.Second
			return s.serverConnectingUpdate(duration), s.serverConnectingDraw(duration), true, nil
		}
		if imgIndex < len(config.GetAvailableImages()) {
			imgToLoad := config.GetAvailableImages()[imgIndex]
			filePath := imgsDirPath + string(imgToLoad) + ".png"
			// @todo: get width and height from config
			img, err := loadImageNineSlice(filePath, 5, 5)
			if err != nil {
				return s.error("loading image", err)
			}
			s.assets.NineSlices[imgToLoad] = img
			log = append(log, "image "+string(imgToLoad)+" loaded")
			imgIndex++
			return s.loadAssetsUpdate(imgIndex, fontIndex, log), s.loadDraw(log), false, nil
		}

		if fontIndex < len(config.GetAvailableFonts()) {
			fontToLoad := config.GetAvailableFonts()[fontIndex]
			fInfo, ok := s.config.Fonts[fontToLoad]
			if !ok {
				return s.error("inconsistent config for fonts", errors.New("need config for "+string(fontToLoad)))
			}
			face, err := loadFont(fontsDirPath+fInfo.FaceFile, fInfo.Size)
			if err != nil {
				return s.error("loading font "+string(fontToLoad), err)
			}
			s.assets.Fonts[fontToLoad] = face
			log = append(log, "font "+string(fontToLoad)+" loaded")
			fontIndex++
			return s.loadAssetsUpdate(imgIndex, fontIndex, log), s.loadDraw(log), false, nil
		}
		return s.staySameState()
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
