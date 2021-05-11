package main

import (
	"github.com/justclimber/ebitenui"
	"github.com/justclimber/ebitenui/widget"
	"github.com/justclimber/fda-client/config"
	"golang.org/x/image/colornames"
)

func (s *SceneMain) setupUI() error {
	rootContainer := widget.NewContainer(
		"root",
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),
			widget.GridLayoutOpts.Padding(widget.NewInsetsSimple(20)),
			widget.GridLayoutOpts.Spacing(0, 20))),
	)
	s.ui = &ebitenui.UI{Container: rootContainer}

	rootContainer.AddChild(widget.NewText(
		widget.TextOpts.Text(
			"Header",
			s.g.assets.Fonts[config.FntDefault],
			s.g.config.Style.WindowsPanel.FontColor,
		),
	))

	mainContainer := widget.NewContainer(
		"main",
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
		)),
	)
	rootContainer.AddChild(mainContainer)

	footerContainer := widget.NewContainer("footer", widget.ContainerOpts.Layout(widget.NewRowLayout()))
	rootContainer.AddChild(footerContainer)

	footerText := widget.NewText(
		widget.TextOpts.Text("Footer", s.g.assets.Fonts[config.FntDefault], s.g.config.Style.WindowsPanel.FontColor))

	footerContainer.AddChild(footerText)

	apps := []*app{
		s.testButtonApp(),
		s.testImageCallbackApp(),
	}

	am := newAppManager(
		s.ui,
		apps,
		s.g.assets.NineSlices[config.ImgWindow],
		widget.NewInsetsSimple(5),
		15,
		s.g.assets.Fonts[config.FntDefault],
		colornames.White)

	// @todo: get colors from config
	listColors := &widget.ListEntryColor{
		Unselected:                 colornames.Gray,
		Selected:                   colornames.Aqua,
		DisabledUnselected:         colornames.Gray,
		DisabledSelected:           colornames.Gray,
		SelectedBackground:         colornames.Darkgray,
		DisabledSelectedBackground: colornames.Darkgray,
	}

	listScrollContainerImages := &widget.ScrollContainerImage{
		Idle:     s.g.assets.NineSlices[config.ImgListIdle],
		Disabled: s.g.assets.NineSlices[config.ImgListDisabled],
		Mask:     s.g.assets.NineSlices[config.ImgListMask],
	}

	appList := widget.NewList(
		widget.ListOpts.Entries(am.appLinks()),
		widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			return e.(appLink).app.title
		}),
		widget.ListOpts.ScrollContainerOpts(widget.ScrollContainerOpts.Image(listScrollContainerImages)),
		widget.ListOpts.EntryColor(listColors),
		widget.ListOpts.EntryFontFace(s.g.assets.Fonts[config.FntDefault]),
		widget.ListOpts.EntryTextPadding(widget.NewInsetsSimple(15)),
		widget.ListOpts.HideHorizontalSlider(),
		widget.ListOpts.HideVerticalSlider(),
		widget.ListOpts.IsMulti(),
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			am.appToggle(args.Entry.(appLink).app)
		}),
		widget.ListOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		}))),
	)
	mainContainer.AddChild(appList)

	return nil
}

func (s *SceneMain) testButtonApp() *app {
	img := s.g.assets.NineSlices[config.ImgButton]
	buttonImage := &widget.ButtonImage{
		Idle:     img,
		Hover:    img,
		Pressed:  img,
		Disabled: img,
	}
	buttonColor := &widget.ButtonTextColor{
		Idle:     colornames.Aqua,
		Disabled: colornames.Aqua,
	}
	b := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("test button", s.g.assets.Fonts[config.FntDefault], buttonColor),
		widget.ButtonOpts.TextPadding(widget.Insets{7, 15, 15, 7}),
	)

	return &app{
		title:   "test button app",
		content: b,
	}
}

func (s *SceneMain) testImageCallbackApp() *app {
	g := widget.NewGraphic(
		widget.GraphicOpts.Callback(s.drawHistoryPlayerCallback),
		widget.GraphicOpts.Image(s.historyPlayerImage),
	)
	return &app{
		title:   "test callback app",
		content: g,
	}
}