package main

import (
	"github.com/justclimber/ebitenui"
	"github.com/justclimber/ebitenui/widget"
	"github.com/justclimber/fda-client/config"
)

func (s *SceneMain) setupUI() {
	rootContainer := widget.NewContainer(
		"root",
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),
			widget.GridLayoutOpts.Padding(widget.Insets{
				Top:    20,
				Bottom: 20,
			}),
			widget.GridLayoutOpts.Spacing(0, 20))),
	)
	s.ui = &ebitenui.UI{Container: rootContainer}

	rootContainer.AddChild(widget.NewText(
		widget.TextOpts.Text(
			"UserHeader",
			s.assets.Fonts[config.FntDefault],
			s.config.Style.WindowsPanel.FontColor,
		),
	))

	mainContainer := widget.NewContainer(
		"main",
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Padding(widget.Insets{
				Left:  25,
				Right: 25,
			}),
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Stretch([]bool{false, true}, []bool{true}),
			widget.GridLayoutOpts.Spacing(20, 0),
		)))
	rootContainer.AddChild(mainContainer)

	footerContainer := widget.NewContainer("footer", widget.ContainerOpts.Layout(widget.NewRowLayout(
		widget.RowLayoutOpts.Padding(widget.Insets{
			Left:  25,
			Right: 25,
		}),
	)))
	rootContainer.AddChild(footerContainer)

	footerText := widget.NewText(
		widget.TextOpts.Text("footer", s.assets.Fonts[config.FntDefault], s.config.Style.WindowsPanel.FontColor))

	footerContainer.AddChild(footerText)
}
