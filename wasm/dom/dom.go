package dom

import (
	"fmt"
	"syscall/js"

	gol "github.com/dfirebaugh/game-of-life-wasm/wasm/life"
)

// Nodes - contains nodes on the DOM
type Nodes struct {
	Grid         js.Value
	PlayBtn      js.Value
	ResetBtn     js.Value
	GenBtn       js.Value
	msgContainer js.Value
	message      js.Value
	genLabel     js.Value
	generation   js.Value
	btnContainer js.Value
	ClearBtn     js.Value
}

var life *gol.Game

// New returns a new Nodes instance
func New(game gol.Game) Nodes {
	n := &Nodes{}
	life = &game

	n.Reset()

	return *n
}

// Reset sets up new DOM nodes and adds the first generation to the Grid
func (n *Nodes) Reset() {
	js.Global().Set("GRID_SIZE", gol.SIZE)

	document := js.Global().Get("document")

	n.Grid = document.Call("getElementById", "gridContainer")

	n.createButtons()

	n.Grid.Set("style", fmt.Sprintf("width: %dpx", gol.GRIDWIDTH))

	n.Grid.Set("innerHTML", "") // reset childNodes

	for y, row := range life.Cells {
		rowNode := document.Call("createElement", "row")
		for x, cell := range row {
			// declaring new variables to break reference to iterator
			cellX := x
			cellY := y

			// handle a cell click
			cellClickCb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				life.Cells[cellY][cellX].Alive = !life.Cells[cellY][cellX].Alive
				n.Render()
				return nil
			})

			cellNode := document.Call("createElement", "cell")
			cellNode.Call("setAttribute", "alive", cell.Alive)
			cellNode.Set("id", fmt.Sprintf("cell-%d-%d", x, y))
			// TODO: move away from pixels and use something more dynamic
			cellNode.Set("style", fmt.Sprintf("width: %dpx; height: %dpx; border: solid %dpx", gol.CELLSIZE, gol.CELLSIZE, gol.CELLBORDERSIZE))
			cellNode.Call("addEventListener", "click", cellClickCb)
			rowNode.Call("appendChild", cellNode)
		}
		n.Grid.Call("appendChild", rowNode)
	}

	n.msgContainer = document.Call("getElementById", "messageContainer")
	n.msgContainer.Set("innerHTML", "") // clear out message container

	n.genLabel = document.Call("createElement", "generation-label")
	n.genLabel.Set("innerHTML", "Generation: 0")

	n.generation = document.Call("createElement", "generation")
	n.generation.Set("innerHTML", life.Generation)

	n.message = document.Call("createElement", "message")
	n.message.Set("innerHTML", life.Message)

	n.msgContainer.Call("appendChild", n.genLabel)
	n.msgContainer.Call("appendChild", n.generation)
	n.msgContainer.Call("appendChild", n.message)
}

// registerCallbacks defines callbacks that are used within JS code
func (n *Nodes) registerCallbacks() {
	var runCb, generateCb js.Func
	runCb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		life.TogglePause()
		// runCb.Release() // release the function if the button will not be clicked again
		return nil
	})

	generateCb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		life.Generate()
		// generateCb.Release() // release the function if the button will not be clicked again
		return nil
	})

	resetCb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		life.ResetCb()

		return nil
	})

	clearCb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		life.ClearGrid()

		return nil
	})

	n.PlayBtn.Call("addEventListener", "click", runCb)
	n.GenBtn.Call("addEventListener", "click", generateCb)
	n.ResetBtn.Call("addEventListener", "click", resetCb)
	n.ClearBtn.Call("addEventListener", "click", clearCb)
}

func (n *Nodes) createButtons() {
	document := js.Global().Get("document")
	n.btnContainer = document.Call("getElementById", "btnContainer")
	n.btnContainer.Set("innerHTML", "") // clear container

	n.PlayBtn = document.Call("createElement", "button")
	n.PlayBtn.Set("id", "runBtn")
	n.btnContainer.Call("appendChild", n.PlayBtn)

	n.ResetBtn = document.Call("createElement", "button")
	n.ResetBtn.Set("innerHTML", "Reset")
	n.ResetBtn.Set("id", "resetBtn")
	n.btnContainer.Call("appendChild", n.ResetBtn)

	n.GenBtn = document.Call("createElement", "button")
	n.GenBtn.Set("innerHTML", "Generate")
	n.GenBtn.Set("id", "generateBtn")
	n.btnContainer.Call("appendChild", n.GenBtn)

	n.ClearBtn = document.Call("createElement", "button")
	n.ClearBtn.Set("innerHTML", "Clear")
	n.ClearBtn.Set("id", "ClearBtn")
	n.btnContainer.Call("appendChild", n.ClearBtn)

	n.registerCallbacks()
}

func (n *Nodes) updateGrid() {
	document := js.Global().Get("document")

	// TODO: this would be better if we only change what needs to be changed
	for y, row := range life.Cells {
		for x, cell := range row {
			domCell := document.Call("getElementById", fmt.Sprintf("cell-%d-%d", x, y))
			domCell.Call("setAttribute", "alive", cell.Alive)
			if gol.SHOWNEIGHBORS {
				domCell.Set("innerHTML", cell.Neighbors)
			}
		}
	}
}

func (n *Nodes) renderMessage() {
	n.genLabel.Set("innerHTML", "Generation: ")
	n.generation.Set("innerHTML", life.Generation)
	n.message.Set("innerHTML", life.Message)
}

// Render - update the display
// we do this by rerending the message and updateing the Grid
func (n *Nodes) Render() {
	n.renderMessage()
	n.setButtons()
	n.updateGrid()
}

func (n *Nodes) setButtons() {
	// TODO: only update this when it changes
	if !life.IsPaused {
		n.PlayBtn.Set("innerHTML", "Pause")
	} else {
		n.PlayBtn.Set("innerHTML", "Play")
	}
}
