package codeeditor

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/justclimber/marslang/ast"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	img "image"
	"image/color"
	"strconv"
)

const identWidth = 3

type CodeColor struct {
	colorDefault color.Color
	colorIdent   color.Color
	colorKeyword color.Color
	colorConst   color.Color
	colorSymbols color.Color
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
	}
}

func (cd *CodeDrawer) newLine(ident int) {
	cd.col = 0
	cd.row++
	cd.ident += ident
}

func (cd *CodeDrawer) draw(screen *ebiten.Image, code *ast.StatementsBlock, rect img.Rectangle) {
	if code == nil {
		return
	}
	rect = rect.Add(img.Point{X: 150, Y: 23})
	cd.reset(screen, rect.Min)
	cd.drawStatementBlock(code)
}

func (cd *CodeDrawer) drawStatementBlock(stmts *ast.StatementsBlock) {
	for i, stmt := range stmts.Statements {
		switch astNode := stmt.(type) {
		case *ast.IfStatement:
			cd.drawIfStatement(astNode)
		case *ast.Assignment:
			cd.drawAssignmentStatement(astNode)
			// @todo: other cases
		}
		if i+1 != len(stmts.Statements) {
			cd.newLine(0)
		}
	}
}

func (cd *CodeDrawer) drawIfStatement(ifStmt *ast.IfStatement) {
	cd.drawText("if ", cd.colors.colorKeyword)
	cd.drawExpression(ifStmt.Condition)
	cd.drawText(" {", cd.colors.colorSymbols)
	cd.newLine(1)
	cd.drawStatementBlock(ifStmt.PositiveBranch)
	cd.newLine(-1)
	cd.drawText("}", cd.colors.colorSymbols)
}

func (cd *CodeDrawer) drawExpression(expr ast.IExpression) {
	switch astNode := expr.(type) {
	case *ast.BinExpression:
		cd.drawBinExpression(astNode)
	case *ast.NumInt:
		cd.drawText(strconv.Itoa(int(astNode.Value)), cd.colors.colorConst)
	case *ast.Identifier:
		cd.drawText(astNode.Value, cd.colors.colorConst)
	}
}

func (cd *CodeDrawer) drawBinExpression(expr *ast.BinExpression) {
	cd.drawExpression(expr.Left)
	cd.drawText(fmt.Sprintf(" %s ", expr.Operator), cd.colors.colorSymbols)
	cd.drawExpression(expr.Right)
}

func (cd *CodeDrawer) drawAssignmentStatement(expr *ast.Assignment) {
	cd.drawText(expr.Left.Value, cd.colors.colorIdent)
	cd.drawText(" = ", cd.colors.colorSymbols)
	cd.drawExpression(expr.Value)
}

func (cd *CodeDrawer) drawText(txt string, clr color.Color) {
	curPosX := cd.origin.X + (cd.ident*identWidth+cd.col)*cd.font.width
	curPosY := cd.origin.Y + cd.row*cd.font.height
	text.Draw(cd.screen, txt, cd.font.face, curPosX, curPosY, clr)
	cd.col += len(txt)
}
