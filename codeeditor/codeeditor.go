package codeeditor

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justclimber/ebitenui/event"
	"github.com/justclimber/ebitenui/image"
	"github.com/justclimber/ebitenui/input"
	"github.com/justclimber/ebitenui/widget"
	"github.com/justclimber/marslang/ast"
	"github.com/justclimber/marslang/lexer"
	"github.com/justclimber/marslang/parser"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	img "image"
	"image/color"
	"log"
	"math"
	"strings"
	"sync/atomic"
	"time"
)

type CodeEditor struct {
	ChangedEvent *event.Event

	Lines []string
	Index []ast.INode
	Code  *ast.StatementsBlock

	codeDrawer CodeDrawer

	widgetOpts      []widget.WidgetOpt
	cursorOpts      []CursorOpt
	image           *BgImage
	colors          *Colors
	padding         widget.Insets
	face            font.Face
	repeatDelay     time.Duration
	repeatInterval  time.Duration
	placeholderText string

	init           *widget.MultiOnce
	commandToFunc  map[controlCommand]codeEditorCommandFunc
	widget         *widget.Widget
	cursor         *Cursor
	text           *widget.Text
	renderBuf      *image.MaskedRenderBuffer
	mask           *image.NineSlice
	cursorPosition img.Point
	state          stateFunc
	scrollOffset   int
	focused        bool
	lastInputText  string
}

type Opt func(t *CodeEditor)

type Options struct {
}

type ChangedEventArgs struct {
	CodeEditor *CodeEditor
	InputText  string
}

type ChangedHandlerFunc func(args *ChangedEventArgs)

type BgImage struct {
	Idle     *image.NineSlice
	Disabled *image.NineSlice
}

type Colors struct {
	Idle          color.Color
	Disabled      color.Color
	Cursor        color.Color
	DisabledCaret color.Color
}

type stateFunc func() (stateFunc, bool)

type controlCommand int

type codeEditorCommandFunc func()

var Opts Options

const (
	textInputGoLeft = controlCommand(iota + 1)
	textInputGoRight
	textInputGoUp
	textInputGoDown
	textInputGoStart
	textInputGoEnd
	textInputBackspace
	textInputDelete
	testInputEnter
)

var keyToCommand = map[ebiten.Key]controlCommand{
	ebiten.KeyLeft:      textInputGoLeft,
	ebiten.KeyRight:     textInputGoRight,
	ebiten.KeyHome:      textInputGoStart,
	ebiten.KeyEnd:       textInputGoEnd,
	ebiten.KeyBackspace: textInputBackspace,
	ebiten.KeyDelete:    textInputDelete,
	ebiten.KeyEnter:     testInputEnter,
	ebiten.KeyUp:        textInputGoUp,
	ebiten.KeyDown:      textInputGoDown,
}

func NewCodeEditor(opts ...Opt) *CodeEditor {
	t := &CodeEditor{
		ChangedEvent:   &event.Event{},
		Lines:          []string{""},
		repeatDelay:    300 * time.Millisecond,
		repeatInterval: 35 * time.Millisecond,

		init:          &widget.MultiOnce{},
		commandToFunc: map[controlCommand]codeEditorCommandFunc{},
		renderBuf:     image.NewMaskedRenderBuffer(),
		codeDrawer:    CodeDrawer{},
	}
	t.state = t.idleState(true)

	t.commandToFunc[textInputGoLeft] = t.doGoLeft
	t.commandToFunc[textInputGoRight] = t.doGoRight
	t.commandToFunc[textInputGoStart] = t.doGoStart
	t.commandToFunc[textInputGoEnd] = t.doGoEnd
	t.commandToFunc[textInputBackspace] = t.doBackspace
	t.commandToFunc[textInputDelete] = t.doDelete
	t.commandToFunc[testInputEnter] = t.doEnter
	t.commandToFunc[textInputGoUp] = t.doGoUp
	t.commandToFunc[textInputGoDown] = t.doGoDown

	t.init.Append(t.createWidget)

	for _, o := range opts {
		o(t)
	}

	return t
}

func (o Options) WidgetOpts(opts ...widget.WidgetOpt) Opt {
	return func(t *CodeEditor) {
		t.widgetOpts = append(t.widgetOpts, opts...)
	}
}

func (o Options) CursorOpts(opts ...CursorOpt) Opt {
	return func(t *CodeEditor) {
		t.cursorOpts = append(t.cursorOpts, opts...)
	}
}

func (o Options) ChangedHandler(f ChangedHandlerFunc) Opt {
	return func(t *CodeEditor) {
		t.ChangedEvent.AddHandler(func(args interface{}) {
			f(args.(*ChangedEventArgs))
		})
	}
}

func (o Options) BgImage(i *BgImage) Opt {
	return func(t *CodeEditor) {
		t.image = i
	}
}

func (o Options) Colors(c *Colors) Opt {
	return func(t *CodeEditor) {
		t.colors = c
	}
}

func (o Options) Padding(i widget.Insets) Opt {
	return func(t *CodeEditor) {
		t.padding = i
	}
}

func (o Options) Face(f font.Face) Opt {
	return func(t *CodeEditor) {
		t.face = f
		t.codeDrawer.font.face = f
		t.codeDrawer.font.buildMetricsCache()
	}
}

func (o Options) RepeatInterval(i time.Duration) Opt {
	return func(t *CodeEditor) {
		t.repeatInterval = i
	}
}

func (o Options) Placeholder(s string) Opt {
	return func(t *CodeEditor) {
		t.placeholderText = s
	}
}

func (ce *CodeEditor) GetWidget() *widget.Widget {
	ce.init.Do()
	return ce.widget
}

func (ce *CodeEditor) SetLocation(rect img.Rectangle) {
	ce.init.Do()
	ce.widget.Rect = rect
}

func (ce *CodeEditor) PreferredSize() (int, int) {
	ce.init.Do()
	_, h := ce.cursor.PreferredSize()
	return 50, h + ce.padding.Top + ce.padding.Bottom
}

func (ce *CodeEditor) Render(screen *ebiten.Image, def widget.DeferredRenderFunc, debugMode widget.DebugMode) {
	ce.init.Do()

	ce.text.GetWidget().Disabled = ce.widget.Disabled

	//if ce.cursorPosition > len([]rune(ce.InputText)) {
	//	ce.cursorPosition = len([]rune(ce.InputText))
	//}

	for {
		newState, rerun := ce.state()
		if newState != nil {
			ce.state = newState
		}
		if !rerun {
			break
		}
	}

	defer func() {
		ce.lastInputText = strings.Join(ce.Lines, "")
	}()

	if strings.Join(ce.Lines, "") != ce.lastInputText {
		ce.ChangedEvent.Fire(&ChangedEventArgs{
			CodeEditor: ce,
			InputText:  strings.Join(ce.Lines, ""),
		})
	}

	ce.widget.Render(screen, def, debugMode)

	ce.renderImage(screen)
	ce.renderTextAndCaret(screen, def, debugMode)
}

func (ce *CodeEditor) idleState(newKeyOrCommand bool) stateFunc {
	return func() (stateFunc, bool) {
		if !ce.focused {
			return ce.idleState(true), false
		}

		chars := input.InputChars()
		if len(chars) > 0 {
			return ce.charsInputState(chars), true
		}

		st := checkForCommand(ce, newKeyOrCommand)
		if st != nil {
			return st, true
		}

		if input.MouseButtonJustPressedLayer(ebiten.MouseButtonLeft, ce.widget.EffectiveInputLayer()) {
			ce.doGoXY(input.CursorPosition())
		}

		return ce.idleState(true), false
	}
}

func checkForCommand(t *CodeEditor, newKeyOrCommand bool) stateFunc {
	for key, cmd := range keyToCommand {
		if !input.KeyPressed(key) {
			continue
		}

		var delay time.Duration
		if newKeyOrCommand {
			delay = t.repeatDelay
		} else {
			delay = t.repeatInterval
		}

		return t.commandState(cmd, key, delay, nil, nil)
	}

	return nil
}

func (ce *CodeEditor) charsInputState(c []rune) stateFunc {
	return func() (stateFunc, bool) {
		if !ce.widget.Disabled {
			ce.doInsert(c)
			ce.codeChanged()
		}

		ce.cursor.ResetBlinking()

		return ce.idleState(true), false
	}
}

func (ce *CodeEditor) codeChanged() {
	stringToParse := strings.Join(ce.Lines, "\n") + "\n"
	l := lexer.New(stringToParse)
	p, err := parser.New(l)
	if err != nil {
		log.Printf("Lexing error: %s\n", err.Error())
	}
	ce.Code, err = p.Parse()
	if err != nil {
		log.Printf("Parsing error: %s\n", err.Error())
	} else {
		ce.buildIndex()
	}
}

func (ce *CodeEditor) buildIndex() {
	lineNumber := 0
	ce.Index = make([]ast.INode, len(ce.Lines))
	ce.buildIndexForStatementsBlock(ce.Code, lineNumber)
}

func (ce *CodeEditor) buildIndexForStatementsBlock(stmts *ast.StatementsBlock, lineNumber int) int {
	for _, stmt := range stmts.Statements {
		ce.Index[lineNumber] = stmt
		lineNumber++
		switch astNode := stmt.(type) {
		case *ast.Switch:
			for _, c := range astNode.Cases {
				ce.Index[lineNumber] = c
				lineNumber++
				lineNumber = ce.buildIndexForStatementsBlock(c.PositiveBranch, lineNumber)
			}
			if astNode.DefaultBranch != nil {
				lineNumber++
				lineNumber = ce.buildIndexForStatementsBlock(astNode.DefaultBranch, lineNumber)
			}
		case *ast.IfStatement:
			lineNumber = ce.buildIndexForStatementsBlock(astNode.PositiveBranch, lineNumber)
			if astNode.ElseBranch != nil {
				lineNumber++
				lineNumber = ce.buildIndexForStatementsBlock(astNode.ElseBranch, lineNumber)
			}
			// @todo: other cases
		}
	}
	return lineNumber
	//spew.Dump(ce.Index)
}

func (ce *CodeEditor) commandState(cmd controlCommand, key ebiten.Key, delay time.Duration, timer *time.Timer, expired *atomic.Value) stateFunc {
	return func() (stateFunc, bool) {
		if !input.KeyPressed(key) {
			return ce.idleState(true), true
		}

		if timer != nil && expired.Load().(bool) {
			return ce.idleState(false), true
		}

		if timer == nil {
			ce.commandToFunc[cmd]()

			expired = &atomic.Value{}
			expired.Store(false)

			timer = time.AfterFunc(delay, func() {
				expired.Store(true)
			})

			return ce.commandState(cmd, key, delay, timer, expired), false
		}

		return nil, false
	}
}

func (ce *CodeEditor) doInsert(c []rune) {
	s := string(insertChars([]rune(ce.Lines[ce.cursorPosition.Y]), c, ce.cursorPosition.X))

	ce.Lines[ce.cursorPosition.Y] = s
	ce.cursorPosition.X += len(c)
}

func (ce *CodeEditor) doGoLeft() {
	if ce.cursorPosition.X > 0 {
		ce.cursorPosition.X--
	}
	ce.cursor.ResetBlinking()
}

func (ce *CodeEditor) doGoRight() {
	if ce.cursorPosition.X < len([]rune(ce.Lines[ce.cursorPosition.Y])) {
		ce.cursorPosition.X++
	}
	ce.cursor.ResetBlinking()
}

func (ce *CodeEditor) doGoStart() {
	ce.cursorPosition.X = 0
	ce.cursor.ResetBlinking()
}

func (ce *CodeEditor) doGoEnd() {
	ce.cursorPosition.X = len([]rune(ce.Lines[ce.cursorPosition.Y]))
	ce.cursor.ResetBlinking()
}

func (ce *CodeEditor) doGoXY(x int, y int) {
	p := img.Point{x, y}
	if p.In(ce.widget.Rect) {
		tr := ce.padding.Apply(ce.widget.Rect)
		if x < tr.Min.X {
			x = tr.Min.X
		}
		if x > tr.Max.X {
			x = tr.Max.X
		}

		ce.cursorPosition.X = fontStringIndex([]rune(ce.Lines[ce.cursorPosition.Y]), ce.face, x-ce.scrollOffset-tr.Min.X)
		ce.cursor.ResetBlinking()
	}
}

func (ce *CodeEditor) doBackspace() {
	if !ce.widget.Disabled && ce.cursorPosition.X > 0 {
		ce.Lines[ce.cursorPosition.Y] = string(removeChar([]rune(ce.Lines[ce.cursorPosition.Y]), ce.cursorPosition.X-1))
		ce.cursorPosition.X--
	} else if ce.cursorPosition.Y > 0 {
		ce.cursorPosition.Y--
		ce.cursorPosition.X = len(ce.Lines[ce.cursorPosition.Y])
		ce.Lines[ce.cursorPosition.Y] += ce.Lines[ce.cursorPosition.Y+1]
		ce.Lines = append(ce.Lines[:ce.cursorPosition.Y+1], ce.Lines[ce.cursorPosition.Y+2:]...)
	}
	ce.cursor.ResetBlinking()
}

func (ce *CodeEditor) doDelete() {
	if !ce.widget.Disabled && ce.cursorPosition.X < len([]rune(ce.Lines[ce.cursorPosition.Y])) {
		ce.Lines[ce.cursorPosition.Y] = string(removeChar([]rune(ce.Lines[ce.cursorPosition.Y]), ce.cursorPosition.X))
	}
	ce.cursor.ResetBlinking()
}

func (ce *CodeEditor) doEnter() {
	if !ce.widget.Disabled {
		left := ce.Lines[ce.cursorPosition.Y][ce.cursorPosition.X:]
		ce.Lines[ce.cursorPosition.Y] = ce.Lines[ce.cursorPosition.Y][:ce.cursorPosition.X]
		if len(ce.Lines) == ce.cursorPosition.Y+1 {
			ce.Lines = append(ce.Lines, left)
			ce.cursorPosition.Y++
		} else {
			ce.cursorPosition.Y++
			ce.Lines = append(ce.Lines[:ce.cursorPosition.Y+1], ce.Lines[ce.cursorPosition.Y:]...)
			ce.Lines[ce.cursorPosition.Y] = left
		}
		ce.cursorPosition.X = 0
	}
	ce.cursor.ResetBlinking()
}

func (ce *CodeEditor) doGoUp() {
	if !ce.widget.Disabled {
		if ce.cursorPosition.Y > 0 {
			ce.cursorPosition.Y--
			if ce.cursorPosition.X > len(ce.Lines[ce.cursorPosition.Y]) {
				ce.cursorPosition.X = len(ce.Lines[ce.cursorPosition.Y])
			}
		}
	}
	ce.cursor.ResetBlinking()
}

func (ce *CodeEditor) doGoDown() {
	if !ce.widget.Disabled {
		if ce.cursorPosition.Y < len(ce.Lines)-1 {
			ce.cursorPosition.Y++
			if ce.cursorPosition.X > len(ce.Lines[ce.cursorPosition.Y]) {
				ce.cursorPosition.X = len(ce.Lines[ce.cursorPosition.Y])
			}
		}
	}
	ce.cursor.ResetBlinking()
}

func insertChars(r []rune, c []rune, pos int) []rune {
	res := make([]rune, len(r)+len(c))
	copy(res, r[:pos])
	copy(res[pos:], c)
	copy(res[pos+len(c):], r[pos:])
	return res
}

func removeChar(r []rune, pos int) []rune {
	res := make([]rune, len(r)-1)
	copy(res, r[:pos])
	copy(res[pos:], r[pos+1:])
	return res
}

func (ce *CodeEditor) renderImage(screen *ebiten.Image) {
	if ce.image != nil {
		i := ce.image.Idle
		if ce.widget.Disabled && ce.image.Disabled != nil {
			i = ce.image.Disabled
		}

		rect := ce.widget.Rect
		i.Draw(screen, rect.Dx(), rect.Dy(), func(opts *ebiten.DrawImageOptions) {
			opts.GeoM.Translate(float64(rect.Min.X), float64(rect.Min.Y))
		})
	}
}

func (ce *CodeEditor) renderTextAndCaret(screen *ebiten.Image, def widget.DeferredRenderFunc, debugMode widget.DebugMode) {
	ce.renderBuf.Draw(screen,
		func(buf *ebiten.Image) {
			ce.drawTextAndCaret(buf, def, debugMode)
		},
		func(buf *ebiten.Image) {
			rect := ce.widget.Rect
			ce.mask.Draw(buf, rect.Dx()-ce.padding.Left-ce.padding.Right, rect.Dy()-ce.padding.Top-ce.padding.Bottom,
				func(opts *ebiten.DrawImageOptions) {
					opts.GeoM.Translate(float64(rect.Min.X+ce.padding.Left), float64(rect.Min.Y+ce.padding.Top))
					opts.CompositeMode = ebiten.CompositeModeCopy
				})
		})
}

func (ce *CodeEditor) drawTextAndCaret(screen *ebiten.Image, def widget.DeferredRenderFunc, debugMode widget.DebugMode) {
	rect := ce.widget.Rect
	tr := rect
	tr = tr.Add(img.Point{ce.padding.Left, ce.padding.Top})

	inputStr := strings.Join(ce.Lines, "\n")

	cx := 0
	if ce.focused {
		sub := string([]rune(inputStr)[:ce.cursorPosition.X])
		cx = fontAdvance(sub, ce.face)

		dx := tr.Min.X + ce.scrollOffset + cx + ce.cursor.Width + ce.padding.Right - rect.Max.X
		if dx > 0 {
			ce.scrollOffset -= dx
		}

		dx = tr.Min.X + ce.scrollOffset + cx - ce.padding.Left - rect.Min.X
		if dx < 0 {
			ce.scrollOffset -= dx
		}
	}

	tr = tr.Add(img.Point{ce.scrollOffset, 0})

	ce.text.SetLocation(tr)
	if len(inputStr) > 0 {
		ce.text.Label = inputStr
	} else {
		ce.text.Label = ce.placeholderText
	}
	if ce.widget.Disabled || len(inputStr) == 0 {
		ce.text.Color = ce.colors.Disabled
	} else {
		ce.text.Color = ce.colors.Idle
	}
	ce.text.Render(screen, def, debugMode)
	ce.codeDrawer.draw(screen, ce.Code, ce.widget.Rect)

	if ce.focused {
		if ce.widget.Disabled {
			ce.cursor.Color = ce.colors.DisabledCaret
		} else {
			ce.cursor.Color = ce.colors.Cursor
		}

		cy := ce.face.Metrics().Height.Ceil() * ce.cursorPosition.Y
		tr = tr.Add(img.Point{cx, cy})
		ce.cursor.SetLocation(tr)

		ce.cursor.Render(screen, def, debugMode)
	}
}

func (ce *CodeEditor) Focus(focused bool) {
	ce.init.Do()
	widget.WidgetFireFocusEvent(ce.widget, focused)
	ce.cursor.resetBlinking()
	ce.focused = focused
}

func (ce *CodeEditor) createWidget() {
	ce.widget = widget.NewWidget(ce.widgetOpts...)
	ce.widgetOpts = nil

	ce.cursor = NewCursor(append(ce.cursorOpts, CursorOpts.Color(ce.colors.Cursor))...)
	ce.cursorOpts = nil

	ce.text = widget.NewText(widget.TextOpts.Text("", ce.face, color.White))

	ce.mask = image.NewNineSliceColor(color.RGBA{255, 0, 255, 255})
}

func fontAdvance(s string, f font.Face) int {
	_, a := font.BoundString(f, s)
	return int(math.Round(fixedInt26_6ToFloat64(a)))
}

// fontStringIndex returns an index into r that corresponds closest to pixel position x
// when string(r) is drawn using f. Pixel position x==0 corresponds to r[0].
func fontStringIndex(r []rune, f font.Face, x int) int {
	start := 0
	end := len(r)
	p := 0
loop:
	for {
		p = start + (end-start)/2
		sub := string(r[:p])
		a := fontAdvance(sub, f)

		switch {
		// x is right of advance
		case x > a:
			if p == start {
				break loop
			}

			start = p

			// x is left of advance
		case x < a:
			if end == p {
				break loop
			}

			end = p

			// x matches advance exactly
		default:
			return p
		}
	}

	if len(r) > 0 {
		a1 := fontAdvance(string(r[:p]), f)
		a2 := fontAdvance(string(r[:p+1]), f)
		if math.Abs(float64(x-a2)) < math.Abs(float64(x-a1)) {
			p++
		}
	}

	return p
}

//goland:noinspection GoSnakeCaseUsage
func fixedInt26_6ToFloat64(i fixed.Int26_6) float64 {
	return float64(i) / (1 << 6)
}