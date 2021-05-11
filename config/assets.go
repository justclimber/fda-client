package config

import (
	"github.com/justclimber/ebitenui/image"
	"github.com/justclimber/ebitenui/widget"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
)

type Assets struct {
	NineSlices map[NineSlicesEnum]*image.NineSlice
	Fonts      map[FontsEnum]font.Face
	Prefabs    Prefabs
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

type Prefabs struct {
	AppPanel AppPanel
}

func NewPrefabs(assets *Assets, config *Config) Prefabs {
	return Prefabs{
		AppPanel: NewAppPanel(assets, config),
	}
}

type AppPanel struct {
	ListOpts []widget.ListOpt
}

func NewAppPanel(assets *Assets, config *Config) AppPanel {
	appPanel := AppPanel{}

	listColors := &widget.ListEntryColor{
		Unselected:                 colornames.Gray,
		Selected:                   colornames.Aqua,
		DisabledUnselected:         colornames.Gray,
		DisabledSelected:           colornames.Gray,
		SelectedBackground:         colornames.Darkgray,
		DisabledSelectedBackground: colornames.Darkgray,
	}

	listScrollContainerImages := &widget.ScrollContainerImage{
		Idle:     assets.NineSlices[ImgListIdle],
		Disabled: assets.NineSlices[ImgListDisabled],
		Mask:     assets.NineSlices[ImgListMask],
	}

	appPanel.ListOpts = []widget.ListOpt{
		widget.ListOpts.ScrollContainerOpts(widget.ScrollContainerOpts.Image(listScrollContainerImages)),
		widget.ListOpts.EntryColor(listColors),
		widget.ListOpts.EntryFontFace(assets.Fonts[FntDefault]),
		widget.ListOpts.EntryTextPadding(widget.Insets{5, 18, 7, 5}),
		widget.ListOpts.HideHorizontalSlider(),
		widget.ListOpts.HideVerticalSlider(),
		widget.ListOpts.IsMulti(),
		widget.ListOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		}))),
	}
	return appPanel
}

func (a AppPanel) Make(
	entries []interface{},
	labelFunc func(interface{}) string,
	selectedHandler func(*widget.ListEntrySelectedEventArgs),
) *widget.List {
	return widget.NewList(append(a.ListOpts, []widget.ListOpt{
		widget.ListOpts.Entries(entries),
		widget.ListOpts.EntryLabelFunc(labelFunc),
		widget.ListOpts.EntrySelectedHandler(selectedHandler),
	}...)...)
}
