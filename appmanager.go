package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justclimber/ebitenui"
	"github.com/justclimber/ebitenui/widget"
	"github.com/justclimber/fda-client/config"
	gImage "image"
)

type appManager struct {
	apps        []*app
	ui          *ebitenui.UI
	padding     widget.Insets
	spacing     int
}

func newAppManager(
	ui *ebitenui.UI,
	apps []*app,
	padding widget.Insets,
	spacing int,
	windowPrefab config.Window,
) *appManager {
	a := &appManager{
		apps:        apps,
		ui:          ui,
		padding:     padding,
		spacing:     spacing,
	}
	a.initApps(windowPrefab)
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

func (a *app) initWindowBoundsAndPos(am *appManager, ew, eh int) {
	if a.w == 0 {
		a.w, a.h = a.content.PreferredSize()
		a.w += am.padding.Dx()
		headerAndExtraSpaceY := 12
		a.h += am.padding.Dy() + am.spacing + headerAndExtraSpaceY
	}

	r := gImage.Rectangle{gImage.Point{0, 0}, gImage.Point{a.w, a.h}}
	r = r.Add(gImage.Point{(ew - a.w) / 2, (eh - a.h) / 2})
	a.window.SetLocation(r)
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

func (am *appManager) initApps(windowPrefab config.Window) {
	for _, a := range am.apps {
		a.window = windowPrefab.Make(a.title, a.content)

		ew, eh := ebiten.WindowSize()
		a.initWindowBoundsAndPos(am, ew, eh)
	}
}
