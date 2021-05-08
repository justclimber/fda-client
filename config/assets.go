package config

import (
	"github.com/justclimber/ebitenui/image"
	"golang.org/x/image/font"
)

type Assets struct {
	NineSlices map[NineSlicesEnum]*image.NineSlice
	Fonts      map[FontsEnum]font.Face
}

func NewAssets() *Assets {
	return &Assets{
		NineSlices: map[NineSlicesEnum]*image.NineSlice{},
		Fonts:      map[FontsEnum]font.Face{},
	}
}

type NineSlicesEnum string

const (
	ImgWindow       NineSlicesEnum = "window"
	ImgButton       NineSlicesEnum = "button"
	ImgListIdle     NineSlicesEnum = "list-idle"
	ImgListDisabled NineSlicesEnum = "list-disabled"
	ImgListMask     NineSlicesEnum = "list-mask"
)

func GetAvailableImages() []NineSlicesEnum {
	return []NineSlicesEnum{
		ImgWindow,
		ImgButton,
		ImgListIdle,
		ImgListDisabled,
		ImgListMask,
	}
}

type FontsEnum string

const (
	FntDefault FontsEnum = "default"
	FntCode    FontsEnum = "code"
)

func GetAvailableFonts() []FontsEnum {
	return []FontsEnum{
		FntDefault,
		FntCode,
	}
}
