package config

import (
	"github.com/justclimber/ebitenui/widget"
	"net"
)

type Config struct {
	//FontNames map[FontsEnum]string `json:"font_names"`
	Style     Style                `json:"style"`
	Server    Server               `json:"server"`
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
