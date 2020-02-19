package main

import (
	dom "github.com/dfirebaugh/game-of-life-wasm/wasm/dom"
	gol "github.com/dfirebaugh/game-of-life-wasm/wasm/life"
)

func main() {
	g := gol.New()
	d := dom.New(g)

	g.RenderCb = func() {
		d.Render()
	}
	g.ResetCb = func() {
		d.Reset()
	}

	g.Init()
}
