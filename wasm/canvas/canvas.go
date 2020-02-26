package canvas

import (
	"math"
	"syscall/js"

	gol "github.com/dfirebaugh/game-of-life-wasm/wasm/life"
)

// CELLSIZE size of cells in pixels
const CELLSIZE = 15

// CELLBORDERSIZE is how big the border for each cell should be
const CELLBORDERSIZE = 1

/*GRIDWIDTH width of the grid is determined by SIZE (row width) * the size of a Cell with an
accomidation for border size of each cell (i.e. the left border and the right border)
*/
const GRIDWIDTH = gol.SIZE*CELLSIZE + ((CELLBORDERSIZE * gol.SIZE) * 2)

type Canvas struct {
}

var (
	width       float64
	height      float64
	mousePos    [2]float64
	ctx         js.Value
	doc         js.Value
	canvasEl    js.Value
	body        js.Value
	renderFrame js.Func
	g           *gol.Game
	cellsize    int
)

// New instantiates a new instance of the Canvas Renderer
func New(game *gol.Game) Canvas {
	c := &Canvas{}

	g = game

	// Init Canvas stuff
	doc = js.Global().Get("document")

	body = doc.Call("getElementById", "body")
	canvasEl = doc.Call("createElement", "canvas")
	canvasEl.Set("id", "gol-canvas")

	body.Call("appendChild", canvasEl)

	width = doc.Get("body").Get("clientWidth").Float()
	height = doc.Get("body").Get("clientHeight").Float() - 20
	canvasEl.Set("style", "height: 100vh; width: 100vw;")

	ctx = canvasEl.Call("getContext", "2d")

	cellsize = int(math.Round((width / float64(CELLSIZE))) / 3)

	return *c
}

// Render - updates the display
func (c Canvas) Render() {
	go c.updateGrid()
}

// Reset clears state and builds the grid on the Canvas
func (c Canvas) Reset() {

	curBodyW := doc.Get("body").Get("clientWidth").Float()
	curBodyH := doc.Get("body").Get("clientHeight").Float()
	if curBodyW != width || curBodyH != height {
		width, height = curBodyW, curBodyH
		canvasEl.Set("width", width)
		canvasEl.Set("height", height)
	}
}

func (c Canvas) updateGrid() {
	curBodyW := doc.Get("body").Get("clientWidth").Float()
	curBodyH := doc.Get("body").Get("clientHeight").Float()
	if curBodyW != width || curBodyH != height {
		width, height = curBodyW, curBodyH
		canvasEl.Set("width", width)
		canvasEl.Set("height", height)
	}

	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Pull window size to handle resize
		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})
	ctx.Call("clearRect", 0, 0, width, height)
	println(cellsize, cellsize, cellsize, cellsize)

	for y, row := range g.Cells {
		for x, cell := range row {
			ctx.Call("beginPath")
			if cell.Alive {
				ctx.Set("fillStyle", "rebeccapurple")
				ctx.Call("fillRect", cellsize*x, cellsize*y, cellsize, cellsize)
			}
			// ctx.Set("strokeStyle", "green")
			// ctx.Call("rect", cellsize*x, cellsize*y, cellsize, cellsize)
			ctx.Call("stroke")
		}
	}
}
