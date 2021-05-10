package config

import (
	"github.com/justclimber/ebitenui/widget"
	"golang.org/x/image/colornames"
	"image/color"
	"net"
)

type Config struct {
	NineSlicesParams map[NineSlicesEnum]NiceSlicesParams
	Fonts            map[FontsEnum]Font `json:"fonts"`
	Style            Style              `json:"style"`
	Server           Server             `json:"server"`
}

type NiceSlicesParams struct {
	Centered bool
	W        [3]int
	H        [3]int
}

type Font struct {
	FaceFile string  `json:"face_file"`
	Size     float64 `json:"size"`
}

type Style struct {
	WindowsPanel WindowsPanel `json:"windows_panel"`
}

type WindowsPanel struct {
	Width     int           `json:"width"`
	Padding   widget.Insets `json:"padding"`
	FontColor color.RGBA    `json:"font_color"`
}

type Server struct {
	Ip   net.IP `json:"ip"`
	Port int    `json:"port"`
}

type Tst struct {
	Ip   net.IP `json:"ip"`
	Port int    `json:"port"`
}

func GetConfig() *Config {
	return &Config{
		NineSlicesParams: map[NineSlicesEnum]NiceSlicesParams{
			ImgWindow: {
				Centered: true,
				W:        [3]int{10, 0, 0},
				H:        [3]int{10, 0, 0},
			},
			ImgButton: {
				Centered: true,
				W:        [3]int{12, 0, 0},
				H:        [3]int{12, 0, 0},
			},
			ImgListIdle: {
				Centered: false,
				W:        [3]int{25, 12, 22},
				H:        [3]int{25, 12, 25},
			},
			ImgListDisabled: {
				Centered: false,
				W:        [3]int{25, 12, 22},
				H:        [3]int{25, 12, 25},
			},
			ImgListMask: {
				Centered: false,
				W:        [3]int{26, 10, 23},
				H:        [3]int{26, 10, 26},
			},
		},
		Fonts: map[FontsEnum]Font{
			FntDefault: {
				FaceFile: "NotoSans-Regular.ttf",
				Size:     20,
			},
			FntCode: {
				FaceFile: "DroidSans.ttf",
				Size:     20,
			},
		},
		Style: Style{
			WindowsPanel: WindowsPanel{
				Width:     2000,
				Padding:   widget.NewInsetsSimple(5),
				FontColor: colornames.White,
			},
		},
		Server: Server{
			Ip:   net.ParseIP("127.0.0.1"),
			Port: 4321,
		},
	}
}
