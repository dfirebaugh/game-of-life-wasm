package main

import (
	dom "github.com/dfirebaugh/game-of-life-wasm/wasm/dom"
	gol "github.com/dfirebaugh/game-of-life-wasm/wasm/life"
)

func main() {
	g := gol.New()
	var d gol.Renderer = dom.New(&g)

	g.Init(d)

}
