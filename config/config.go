package config

import (
	"github.com/justclimber/ebitenui/widget"
	"net"
)

type Config struct {
	Fonts  map[FontsEnum]Font `json:"fonts"`
	Style  Style              `json:"style"`
	Server Server             `json:"server"`
}

type Font struct {
	FaceFile string  `json:"face_file"`
	Size     float64 `json:"size"`
}

type Style struct {
	WindowsPanel WindowsPanel `json:"windows_panel"`
}

type WindowsPanel struct {
	Width   int           `json:"width"`
	Padding widget.Insets `json:"padding"`
}

type Server struct {
	Ip   net.IP `json:"ip"`
	Port int    `json:"port"`
}

type Tst struct {
	Ip   net.IP `json:"ip"`
	Port int    `json:"port"`
}
