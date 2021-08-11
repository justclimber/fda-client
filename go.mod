module github.com/justclimber/fda-client

go 1.16

replace (
	github.com/justclimber/ebitenui => ../ebitenui/
	github.com/justclimber/fda-lang => ../fda-lang/
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/hajimehoshi/ebiten/v2 v2.1.0
	github.com/justclimber/ebitenui v0.0.0-20210501081741-62b0ef890536
	github.com/justclimber/fda-lang v0.0.0-20210514065505-c963f114ac9b // indirect
	github.com/justclimber/marslang v0.0.0-20210308090925-50aa8f139075
	golang.org/x/exp v0.0.0-20210503015746-b3083d562e1d // indirect
	golang.org/x/image v0.0.0-20210220032944-ac19c3e999fb
	golang.org/x/sys v0.0.0-20210503173754-0981d6026fa6 // indirect
)
