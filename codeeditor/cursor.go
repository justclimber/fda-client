package codeeditor

import (
	img "image"
	"image/color"
	"math"
	"sync/atomic"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justclimber/ebitenui/image"
	"github.com/justclimber/ebitenui/widget"
	"golang.org/x/image/font"
)

type Cursor struct {
	Width int
	Color color.Color

	face          font.Face
	blinkInterval time.Duration

	init    *widget.MultiOnce
	widget  *widget.Widget
	image   *image.NineSlice
	height  int
	state   cursorBlinkState
	visible bool
}

type CursorOpt func(c *Cursor)

type CursorOptions struct {
}

var CursorOpts CursorOptions

type cursorBlinkState func() cursorBlinkState

func NewCursor(opts ...CursorOpt) *Cursor {
	c := &Cursor{
		blinkInterval: 450 * time.Millisecond,

		init: &widget.MultiOnce{},
	}
	c.resetBlinking()

	c.init.Append(c.createWidget)

	for _, o := range opts {
		o(c)
	}

	return c
}

func (o CursorOptions) Color(c color.Color) CursorOpt {
	return func(ca *Cursor) {
		ca.Color = c
	}
}

func (o CursorOptions) Size(face font.Face, width int) CursorOpt {
	return func(c *Cursor) {
		c.face = face
		c.Width = width
	}
}

func (c *Cursor) GetWidget() *widget.Widget {
	c.init.Do()
	return c.widget
}

func (c *Cursor) SetLocation(rect img.Rectangle) {
	c.init.Do()
	c.widget.Rect = rect
}

func (c *Cursor) PreferredSize() (int, int) {
	c.init.Do()
	return c.Width, c.height
}

func (c *Cursor) Render(screen *ebiten.Image, def widget.DeferredRenderFunc, debugMode widget.DebugMode) {
	c.init.Do()

	c.state = c.state()

	c.widget.Render(screen, def, debugMode)

	if !c.visible {
		return
	}

	c.image = image.NewNineSliceColor(c.Color)

	c.image.Draw(screen, c.Width, c.height, func(opts *ebiten.DrawImageOptions) {
		p := c.widget.Rect.Min
		opts.GeoM.Translate(float64(p.X), float64(p.Y))
	})
}

func (c *Cursor) ResetBlinking() {
	c.init.Do()
	c.resetBlinking()
}

func (c *Cursor) resetBlinking() {
	c.state = c.blinkState(true, nil, nil)
}

func (c *Cursor) blinkState(visible bool, timer *time.Timer, expired *atomic.Value) cursorBlinkState {
	return func() cursorBlinkState {
		c.visible = visible

		if timer != nil && expired.Load().(bool) {
			return c.blinkState(!visible, nil, nil)
		}

		if timer == nil {
			expired = &atomic.Value{}
			expired.Store(false)

			timer = time.AfterFunc(c.blinkInterval, func() {
				expired.Store(true)
			})
		}

		return c.blinkState(visible, timer, expired)
	}
}

func (c *Cursor) createWidget() {
	c.widget = widget.NewWidget()

	m := c.face.Metrics()
	c.height = int(math.Round(fixedInt26_6ToFloat64(m.Ascent + m.Descent)))
	c.face = nil
}
