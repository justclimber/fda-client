package config

import (
	"net"
)

type Config struct {
	NineSlicesParams map[NineSlicesEnum]NiceSlicesParams
	Fonts            map[FontsEnum]Font
	Server           Server
}

type NiceSlicesParams struct {
	Centered bool
	W        [3]int
	H        [3]int
}

type Font struct {
	FaceFile string
	Size     float64
}

type Server struct {
	Ip   net.IP
	Port int
}

type Tst struct {
	Ip   net.IP
	Port int
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
		Server: Server{
			Ip:   net.ParseIP("127.0.0.1"),
			Port: 4321,
		},
	}
}
