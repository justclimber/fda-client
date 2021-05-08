package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justclimber/ebitenui"
	"github.com/justclimber/ebitenui/image"
	"github.com/justclimber/ebitenui/widget"
	"golang.org/x/image/font"
	gImage "image"
	"image/color"
)

type appManager struct {
	apps        []*app
	ui          *ebitenui.UI
	bgImage     *image.NineSlice
	padding     widget.Insets
	spacing     int
	face        font.Face
	headerColor color.Color
}

func newAppManager(
	ui *ebitenui.UI,
	apps []*app,
	bgImage *image.NineSlice,
	padding widget.Insets,
	spacing int,
	face font.Face,
	headerColor color.Color,
) *appManager {
	a := &appManager{
		apps:        apps,
		ui:          ui,
		bgImage:     bgImage,
		padding:     padding,
		spacing:     spacing,
		face:        face,
		headerColor: headerColor,
	}
	a.initApps()
	return a
}

type openClosedState int8

const (
	stateClosed openClosedState = iota
	stateOpen
)

type app struct {
	title           string
	content         widget.PreferredSizeLocateableWidget
	window          *widget.Window
	pos             gImage.Point
	w               int
	h               int
	openClosedState openClosedState
	windowCloseFunc ebitenui.RemoveWindowFunc
}

type appLink struct {
	app *app
}

func (am *appManager) appLinks() []interface{} {
	result := make([]interface{}, len(am.apps))
	for i, a := range am.apps {
		result[i] = appLink{app: a}
	}
	return result
}

func (am *appManager) appToggle(app *app) {
	if app.openClosedState == stateOpen {
		app.windowCloseFunc()
		app.openClosedState = stateClosed
		return
	}
	app.windowCloseFunc = am.ui.AddWindow(app.window)
	app.openClosedState = stateOpen
}

func (am *appManager) calcWindowParams(widgetable widget.PreferredSizeLocateableWidget) (int, int, gImage.Point) {
	w, h := widgetable.PreferredSize()
	w += am.padding.Dx()
	// @todo: calculate headerAndExtraSpaceY
	headerAndExtraSpaceY := 12
	h += am.padding.Dy() + am.spacing + headerAndExtraSpaceY
	ew, eh := ebiten.WindowSize()
	x := (ew - w) / 2
	y := (eh - h) / 2
	return w, h, gImage.Point{x, y}
}

func (am *appManager) initApps() {
	for _, app := range am.apps {
		c := widget.NewContainer(
			"app "+app.title,
			widget.ContainerOpts.BackgroundImage(am.bgImage),
			widget.ContainerOpts.Layout(widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Padding(am.padding),
				widget.RowLayoutOpts.Spacing(am.spacing),
			)),
		)

		mc := widget.NewContainer(
			"app "+app.title+" movable",
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
			widget.TextOpts.Text(app.title, am.face, am.headerColor),
			widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionCenter),
		))

		c.AddChild(app.content)

		app.window = widget.NewWindow(
			widget.WindowOpts.Movable(mc),
			widget.WindowOpts.Contents(c),
		)

		app.w, app.h, app.pos = am.calcWindowParams(app.content)

		r := gImage.Rectangle{gImage.Point{0, 0}, gImage.Point{app.w, app.h}}
		r = r.Add(app.pos)
		app.window.SetLocation(r)
	}
}
