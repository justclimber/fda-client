package codeeditor

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/justclimber/fda-lang/fdalang"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	img "image"
	"image/color"
)

const identWidth = 3

type CodeColor struct {
	colorDefault color.Color
	colorIdent   color.Color
	colorKeyword color.Color
	colorConst   color.Color
	colorSymbols color.Color
	colorInvalid color.Color
	colorType    color.Color
}

type CodeFont struct {
	face   font.Face
	height int
	width  int
}

func (cf *CodeFont) buildMetricsCache() {
	cf.height = cf.face.Metrics().Height.Floor()
	_, w, _ := cf.face.GlyphBounds('a')
	cf.width = w.Round()
}

type CodeDrawer struct {
	screen *ebiten.Image
	origin img.Point
	ident  int
	col    int
	row    int
	font   CodeFont
	colors CodeColor
}

func (cd *CodeDrawer) reset(screen *ebiten.Image, p img.Point) {
	cd.screen = screen
	cd.ident = 0
	cd.origin = p
	cd.col = 0
	cd.row = 0
	cd.colors = CodeColor{
		colorDefault: colornames.Aqua,
		colorIdent:   colornames.Blue,
		colorKeyword: colornames.Aquamarine,
		colorConst:   colornames.Green,
		colorSymbols: colornames.Whitesmoke,
		colorInvalid: colornames.Red,
		colorType:    colornames.Aqua,
	}
}

func (cd *CodeDrawer) newLine(ident int) {
	cd.col = 0
	cd.row++
	cd.ident += ident
}

func (cd *CodeDrawer) drawLinesTokens(screen *ebiten.Image, code [][]fdalang.Token, rect img.Rectangle) {
	if code == nil {
		return
	}
	rect = rect.Add(img.Point{X: 0, Y: 20})
	cd.reset(screen, rect.Min)
	for _, tokens := range code {
		cd.drawTokens(tokens)
		cd.newLine(0)
	}
}

func (cd *CodeDrawer) tokenColorMap() map[fdalang.TokenID]color.Color {
	return map[fdalang.TokenID]color.Color{
		fdalang.TokenNumInt:   cd.colors.colorConst,
		fdalang.TokenNumFloat: cd.colors.colorConst,
		fdalang.TokenTrue:     cd.colors.colorConst,
		fdalang.TokenFalse:    cd.colors.colorConst,
		fdalang.TokenIdent:    cd.colors.colorIdent,
		fdalang.TokenStruct:   cd.colors.colorKeyword,
		fdalang.TokenIf:       cd.colors.colorKeyword,
		fdalang.TokenElse:     cd.colors.colorKeyword,
		fdalang.TokenSwitch:   cd.colors.colorKeyword,
		fdalang.TokenCase:     cd.colors.colorKeyword,
		fdalang.TokenDefault:  cd.colors.colorKeyword,
		fdalang.TokenFunction: cd.colors.colorKeyword,
		fdalang.TokenReturn:   cd.colors.colorKeyword,
		fdalang.TokenEnum:     cd.colors.colorKeyword,
		fdalang.TokenType:     cd.colors.colorType,
		fdalang.TokenInvalid:  cd.colors.colorInvalid,
	}
}

func (cd *CodeDrawer) colorForToken(t fdalang.TokenID) color.Color {
	c, ok := cd.tokenColorMap()[t]
	if !ok {
		return cd.colors.colorSymbols
	}
	return c
}

func (cd *CodeDrawer) drawTokens(tokens []fdalang.Token) {
	if len(tokens) == 0 {
		return
	}
	spaceBeforeTokens := map[fdalang.TokenID]bool{
		fdalang.TokenPlus:       true,
		fdalang.TokenMinus:      true,
		fdalang.TokenAssignment: true,
	}
	spaceAfterTokens := map[fdalang.TokenID]bool{
		fdalang.TokenPlus:       true,
		fdalang.TokenMinus:      true,
		fdalang.TokenAssignment: true,
		fdalang.TokenComma:      true,
	}

	prevTokenID := fdalang.TokenSOL
	prevSpacePotential := false
	if tokens[0].ID == fdalang.TokenRBrace && cd.ident > 0 {
		cd.ident--
	}
	for _, t := range tokens {
		spaceBeforeTokensStr := ""
		spaceAfterTokensStr := ""
		spacePotential := t.ID == fdalang.TokenIdent ||
			fdalang.TokensKeywords()[t.ID] ||
			fdalang.TokensConstants()[t.ID]
		if _, ok := spaceBeforeTokens[t.ID]; ok || (spacePotential && prevSpacePotential) {
			spaceBeforeTokensStr = " "
		}
		if _, ok := spaceAfterTokens[t.ID]; ok {
			spaceAfterTokensStr = " "
		}
		cd.drawText(spaceBeforeTokensStr+t.Value+spaceAfterTokensStr, cd.colorForToken(t.ID))
		prevTokenID = t.ID
		prevSpacePotential = spacePotential
	}
	if prevTokenID == fdalang.TokenLBrace {
		cd.ident++
	}
}

func (cd *CodeDrawer) drawText(txt string, clr color.Color) {
	curPosX := cd.origin.X + (cd.ident*identWidth+cd.col)*cd.font.width
	curPosY := cd.origin.Y + cd.row*cd.font.height
	text.Draw(cd.screen, txt, cd.font.face, curPosX, curPosY, clr)
	cd.col += len(txt)
}
