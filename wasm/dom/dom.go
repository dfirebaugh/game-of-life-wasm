package dom

import (
	"fmt"
	"syscall/js"

	gol "github.com/dfirebaugh/game-of-life-wasm/wasm/life"
)

// CELLSIZE size of cells in pixels
const CELLSIZE = 15

// SHOWNEIGHBORS render how many neighbors?
const SHOWNEIGHBORS = false

// CELLBORDERSIZE is how big the border for each cell should be
const CELLBORDERSIZE = 1

/*GRIDWIDTH width of the grid is determined by SIZE (row width) * the size of a Cell with an
accomidation for border size of each cell (i.e. the left border and the right border)
*/
const GRIDWIDTH = gol.SIZE*CELLSIZE + ((CELLBORDERSIZE * gol.SIZE) * 2)

// Nodes - contains nodes on the DOM
type Nodes struct {
}

var g *gol.Game
var document js.Value

// New returns a new Nodes instance
func New(game *gol.Game) Nodes {
	document = js.Global().Get("document")
	n := &Nodes{}

	g = game

	n.Reset()

	return *n
}

// Render - update the display
// we do this by rerending the message and updateing the Grid
func (n Nodes) Render() {
	n.setButtons()
	n.updateGrid()
	n.renderMessage()
}

// Reset sets up new DOM nodes and adds the first generation to the Grid
func (n Nodes) Reset() {
	js.Global().Set("GRID_SIZE", gol.SIZE)

	Grid := document.Call("getElementById", "gridContainer")

	n.createButtons()

	Grid.Set("style", fmt.Sprintf("width: %dpx", GRIDWIDTH))

	Grid.Set("innerHTML", "") // reset childNodes

	for y := 0; y < gol.SIZE; y++ {
		rowNode := document.Call("createElement", "row")
		for x := 0; x < gol.SIZE; x++ {
			// declaring new variables to break reference to iterator
			cellX := x
			cellY := y

			// handle a cell click
			cellClickCb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				g.ToggleCell(cellX, cellY)
				return nil
			})

			cellNode := document.Call("createElement", "cell")
			cellNode.Set("id", fmt.Sprintf("cell-%d-%d", x, y))

			// TODO: move away from pixels and use something more dynamic
			cellNode.Set("style", fmt.Sprintf("width: %dpx; height: %dpx; border: solid %dpx", CELLSIZE, CELLSIZE, CELLBORDERSIZE))
			cellNode.Call("addEventListener", "click", cellClickCb)
			rowNode.Call("appendChild", cellNode)
		}
		Grid.Call("appendChild", rowNode)
	}

	msgContainer := document.Call("getElementById", "messageContainer")
	msgContainer.Set("innerHTML", "") // clear out message container

	genLabel := document.Call("createElement", "generation-label")
	generation := document.Call("createElement", "generation")
	message := document.Call("createElement", "message")

	generation.Set("id", "generation")
	message.Set("id", "message")
	genLabel.Set("innerHTML", "Generation: ")

	msgContainer.Call("appendChild", genLabel)
	msgContainer.Call("appendChild", generation)
	msgContainer.Call("appendChild", message)
}

// registerCallbacks defines callbacks that are used within JS code
func (n Nodes) registerCallbacks() {
	runCb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		g.TogglePause()
		n.setButtons()

		// runCb.Release() // release the function if the button will not be clicked again
		return nil
	})

	generateCb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		g.Generate()
		// generateCb.Release() // release the function if the button will not be clicked again
		return nil
	})

	resetCb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		g.Reset()

		return nil
	})

	clearCb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		g.ClearGrid()

		return nil
	})

	PlayBtn := document.Call("getElementById", "runBtn")
	ResetBtn := document.Call("getElementById", "resetBtn")
	GenBtn := document.Call("getElementById", "generateBtn")
	ClearBtn := document.Call("getElementById", "ClearBtn")

	PlayBtn.Call("addEventListener", "click", runCb)
	GenBtn.Call("addEventListener", "click", generateCb)
	ResetBtn.Call("addEventListener", "click", resetCb)
	ClearBtn.Call("addEventListener", "click", clearCb)
}

func (n Nodes) createButtons() {
	btnContainer := document.Call("getElementById", "btnContainer")
	PlayBtn := document.Call("createElement", "button")
	ResetBtn := document.Call("createElement", "button")
	GenBtn := document.Call("createElement", "button")
	ClearBtn := document.Call("createElement", "button")

	btnContainer.Set("innerHTML", "") // clear container

	PlayBtn.Set("id", "runBtn")
	ResetBtn.Set("innerHTML", "Reset")
	ResetBtn.Set("id", "resetBtn")
	GenBtn.Set("innerHTML", "Generate")
	GenBtn.Set("id", "generateBtn")
	ClearBtn.Set("innerHTML", "Clear")
	ClearBtn.Set("id", "ClearBtn")

	btnContainer.Call("appendChild", PlayBtn)
	btnContainer.Call("appendChild", ResetBtn)
	btnContainer.Call("appendChild", GenBtn)
	btnContainer.Call("appendChild", ClearBtn)

	n.registerCallbacks()
}

func (n Nodes) updateGrid() {

	// TODO: this would be better if we only change what needs to be changed
	for y, row := range g.Cells {
		for x, cell := range row {
			domCell := document.Call("getElementById", fmt.Sprintf("cell-%d-%d", x, y))
			domCell.Call("setAttribute", "alive", cell.Alive)
			if SHOWNEIGHBORS {
				domCell.Set("innerHTML", cell.Neighbors)
			}
		}
	}
}

func (n Nodes) renderMessage() {
	generation := document.Call("getElementById", "generation")
	message := document.Call("getElementById", "message")

	generation.Set("innerHTML", g.Generation)
	message.Set("innerHTML", g.Message)
}

func (n Nodes) setButtons() {
	PlayBtn := document.Call("getElementById", "runBtn")

	// TODO: only update this when it changes
	if !g.IsPaused {
		PlayBtn.Set("innerHTML", "Pause")
	} else {
		PlayBtn.Set("innerHTML", "Play")
	}
}
