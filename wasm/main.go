package main

import (
	"github.com/dfirebaugh/game-of-life-wasm/wasm/canvas"
	gol "github.com/dfirebaugh/game-of-life-wasm/wasm/life"
)

func main() {
	g := gol.New()

	// var d gol.Renderer = dom.New(&g)
	// g.Init(d)

	var c gol.Renderer = canvas.New(&g)
	g.Init(c)
}
