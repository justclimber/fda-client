package config

import (
	"github.com/justclimber/ebitenui/image"
	"golang.org/x/image/font"
)

type Assets struct {
	NineSlices map[NineSlicesEnum]*image.NineSlice
	Fonts      map[FontsEnum]*font.Face
}

func NewAssets() *Assets {
	return &Assets{
		NineSlices: map[NineSlicesEnum]*image.NineSlice{},
		Fonts:      map[FontsEnum]*font.Face{},
	}
}

type NineSlicesEnum string

const (
	imgWindow NineSlicesEnum = "window"
	imgButton NineSlicesEnum = "button"
)

func GetAvailableImages() []NineSlicesEnum {
	return []NineSlicesEnum{
		imgWindow,
		imgButton,
	}
}

type FontsEnum string

const (
	fntDefault FontsEnum = "default"
	fntCode    FontsEnum = "code"
)

func GetAvailableFonts() []FontsEnum {
	return []FontsEnum{
		fntDefault,
		fntCode,
	}
}
