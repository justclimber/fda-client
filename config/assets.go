package config

import (
	"github.com/justclimber/ebitenui/image"
	"github.com/justclimber/ebitenui/widget"
	"github.com/justclimber/fda-client/codeeditor"
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
	ImgWindow             NineSlicesEnum = "window"
	ImgButton             NineSlicesEnum = "button"
	ImgListIdle           NineSlicesEnum = "list-idle"
	ImgListDisabled       NineSlicesEnum = "list-disabled"
	ImgListMask           NineSlicesEnum = "list-mask"
	ImgCodeEditorIdle     NineSlicesEnum = "codeeditor-idle"
	ImgCodeEditorDisabled NineSlicesEnum = "codeeditor-disabled"
	ImgCodeEditorHover    NineSlicesEnum = "codeeditor-hover"
)

func GetAvailableImages() []NineSlicesEnum {
	return []NineSlicesEnum{
		ImgWindow,
		ImgButton,
		ImgListIdle,
		ImgListDisabled,
		ImgListMask,
		ImgCodeEditorIdle,
		ImgCodeEditorDisabled,
		ImgCodeEditorHover,
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
	AppPanel         AppPanel
	Window           Window
	CodeEditorPrefab CodeEditorPrefab
}

func NewPrefabs(assets *Assets, config *Config) Prefabs {
	return Prefabs{
		AppPanel:         NewAppPanel(assets, config),
		Window:           NewWindow(assets, config),
		CodeEditorPrefab: NewCodeEditorPrefab(assets, config),
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

type Window struct {
	bgImage    *image.NineSlice
	WindowOpts []widget.WindowOpt
	fontFace   font.Face
}

func NewWindow(assets *Assets, config *Config) Window {
	return Window{
		bgImage:  assets.NineSlices[ImgWindow],
		fontFace: assets.Fonts[FntDefault],
	}
}

func (w Window) Make(title string, content widget.PreferredSizeLocateableWidget) *widget.Window {
	content.GetWidget().LayoutData = widget.AnchorLayoutData{
		StretchHorizontal:  true,
		StretchVertical:    true,
	}

	c := widget.NewContainer(
		"app "+title,
		widget.ContainerOpts.BackgroundImage(w.bgImage),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(15)),
		)),
	)

	c.AddChild(content)

	mc := widget.NewContainer(
		"app "+title+" movable",
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical)),
		),
	)

	mc.AddChild(widget.NewText(
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.TextOpts.Text(title, w.fontFace, colornames.White),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionCenter),
	))

	return widget.NewWindow(
		widget.WindowOpts.Movable(mc),
		widget.WindowOpts.Contents(c),
	)
}

type CodeEditorPrefab struct {
	opts []codeeditor.Opt
}

func NewCodeEditorPrefab(assets *Assets, config *Config) CodeEditorPrefab {
	bgImage := &codeeditor.BgImage{
		Idle:     assets.NineSlices[ImgCodeEditorIdle],
		Disabled: assets.NineSlices[ImgCodeEditorDisabled],
	}

	colors := &codeeditor.Colors{
		Idle:          colornames.White,
		Disabled:      colornames.Gray,
		Cursor:        colornames.Aqua,
		DisabledCaret: colornames.Gray,
	}
	return CodeEditorPrefab{
		opts: []codeeditor.Opt{
			codeeditor.Opts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			})),
			codeeditor.Opts.BgImage(bgImage),
			codeeditor.Opts.Placeholder("Enter code here"),
			codeeditor.Opts.Face(assets.Fonts[FntCode]),
			codeeditor.Opts.CursorOpts(
				codeeditor.CursorOpts.Size(assets.Fonts[FntCode], 2),
			),
			codeeditor.Opts.Colors(colors),
			codeeditor.Opts.Padding(widget.Insets{
				Left:   13,
				Right:  13,
				Top:    7,
				Bottom: 7,
			}),
		},
	}
}

func (c CodeEditorPrefab) Make() *codeeditor.CodeEditor {
	return codeeditor.NewCodeEditor(c.opts...)
}
