package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"syscall/js"
	"time"
)

// SIZE size of the grid
//- i.e. the width and height of the grid
const SIZE = 45

// CELLSIZE size of cells in pixels
const CELLSIZE = 15

// TICKSPEED is how fast the generations update
const TICKSPEED = 1

// CELLBORDERSIZE is how big the border for each cell should be
const CELLBORDERSIZE = 1

// PRINT -- should we print to the log
const PRINT = false

/*GRIDWIDTH width of the grid is determined by SIZE (row width) * the size of a Cell with an
accomidation for border size of each cell (i.e. the left border and the right border)
*/
const GRIDWIDTH = SIZE*CELLSIZE + ((CELLBORDERSIZE * SIZE) * 2)

// SHOWNEIGHBORS - renders neighbors to dom
const SHOWNEIGHBORS = false

type cell struct {
	alive       bool
	neighbors   int
	coordinates coords
}

type coords struct {
	x int
	y int
}

// DOMNodes - contains nodes on the DOM
type DOMNodes struct {
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

// Game holds the state of our cells
type Game struct {
	isPaused   bool
	generation int
	message    string
	printlog   bool
	speed      time.Duration
	cells      [SIZE][SIZE]cell
	dom        DOMNodes
}

// initGame - setup game
func (g *Game) initGame() {
	js.Global().Set("GRID_SIZE", SIZE)
	g.printlog = PRINT
	g.speed = TICKSPEED

	g.reset()

	g.startPolling()
}

func (g *Game) reset() {
	g.isPaused = true
	g.generation = 0
	g.message = ""
	g.updateCells(randomAlive)
	g.initDOMNodes()
	g.render()
	g.logger("reset!")
}

func (g *Game) clearGrid() {

	g.updateCells(func(c cell, xy coords) cell {
		c.alive = false
		g.countNeighbors(c, xy)
		return c
	})

	g.isPaused = true
	g.generation = 0
	g.message = "cleared"
	g.logger("cleared")
	g.render()
}

// registerCallbacks defines callbacks that are used within JS code
func (g *Game) registerCallbacks() {
	var runCb, generateCb js.Func
	runCb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		g.togglePause()
		// runCb.Release() // release the function if the button will not be clicked again
		return nil
	})

	generateCb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		g.generate()
		// generateCb.Release() // release the function if the button will not be clicked again
		return nil
	})

	resetCb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		g.reset()

		return nil
	})

	clearCb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		g.clearGrid()

		return nil
	})

	g.dom.PlayBtn.Call("addEventListener", "click", runCb)
	g.dom.GenBtn.Call("addEventListener", "click", generateCb)
	g.dom.ResetBtn.Call("addEventListener", "click", resetCb)
	g.dom.ClearBtn.Call("addEventListener", "click", clearCb)
}

// chunk splits a string into string slice of a given length
func chunk(str string, n int) []string {
	var acc []string

	for i := 0; i <= len(str)-n; i += n {
		acc = append(acc, string(str[i:i+n]))
	}

	return acc
}

// logger is a simple function to disable logs if Game.printlog is set to false
func (g *Game) logger(a ...interface{}) {
	if g.printlog {
		fmt.Println(strings.Trim(fmt.Sprintf("%v", a), "[]"))
	}
}

func (g *Game) countNeighbors(c cell, xy coords) cell {
	neighbors := []coords{
		{-1, 0}, {-1, 1},
		{1, 0}, {1, -1},
		{0, -1}, {-1, -1},
		{0, 1}, {1, 1},
	}
	neighborCount := 0

	for _, neighbor := range neighbors {
		nX := xy.x + neighbor.x
		nY := xy.y + neighbor.y
		if nX >= 0 && nY >= 0 && nX < SIZE && nY < SIZE {
			if g.cells[nX][nY].alive {
				neighborCount++
			}
		}
	}
	c.neighbors = neighborCount
	return c
}

func (g *Game) checkRules(c cell, xy coords) cell {
	if c.alive {
		if c.neighbors < 2 || c.neighbors > 3 {
			c.alive = false
		}
	} else if c.neighbors == 3 {
		c.alive = true
	}
	return c
}

// iterate steps through the graph and modifies the
// cell based on the operation passed in
func (g *Game) updateCells(changeCell func(c cell, xy coords) cell) {
	var i, j int
	for i = 0; i < SIZE; i++ {
		for j = 0; j < SIZE; j++ {
			g.cells[i][j] = changeCell(g.cells[i][j], coords{i, j})
		}
	}
}

func randomAlive(c cell, xy coords) cell {
	c.alive = d2(int64(xy.x + xy.y))
	return c
}

// flip a d2 to get a random true of false value
func d2(seed int64) bool {
	s1 := rand.NewSource(time.Now().UnixNano() + seed)
	r1 := rand.New(s1)
	return r1.Intn(2) == 1
}

func (g *Game) render() {
	g.renderMessage()
	g.updateDOMGrid()
}

func (g *Game) createButtons() {
	document := js.Global().Get("document")
	g.dom.btnContainer = document.Call("getElementById", "btnContainer")
	g.dom.btnContainer.Set("innerHTML", "") // clear container

	g.dom.PlayBtn = document.Call("createElement", "button")
	g.dom.PlayBtn.Set("id", "runBtn")
	g.dom.btnContainer.Call("appendChild", g.dom.PlayBtn)

	g.dom.ResetBtn = document.Call("createElement", "button")
	g.dom.ResetBtn.Set("innerHTML", "Reset")
	g.dom.ResetBtn.Set("id", "resetBtn")
	g.dom.btnContainer.Call("appendChild", g.dom.ResetBtn)

	g.dom.GenBtn = document.Call("createElement", "button")
	g.dom.GenBtn.Set("innerHTML", "Generate")
	g.dom.GenBtn.Set("id", "generateBtn")
	g.dom.btnContainer.Call("appendChild", g.dom.GenBtn)

	g.dom.ClearBtn = document.Call("createElement", "button")
	g.dom.ClearBtn.Set("innerHTML", "Clear")
	g.dom.ClearBtn.Set("id", "ClearBtn")
	g.dom.btnContainer.Call("appendChild", g.dom.ClearBtn)

	g.registerCallbacks()
}

func (g *Game) initDOMNodes() {
	document := js.Global().Get("document")

	g.dom.Grid = document.Call("getElementById", "gridContainer")

	g.createButtons()

	g.dom.Grid.Set("style", fmt.Sprintf("width: %dpx", GRIDWIDTH))

	g.dom.Grid.Set("innerHTML", "") // reset childNodes

	for y, row := range g.cells {
		rowNode := document.Call("createElement", "row")
		for x, cell := range row {
			// declaring new variables to break reference to iterator
			cellX := x
			cellY := y

			// handle a cell click
			cellClickCb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				g.logger("cell clicked: ", cellX, cellY, g.cells[cellY][cellX].alive)
				g.cells[cellY][cellX].alive = !g.cells[cellY][cellX].alive
				g.render()
				return nil
			})

			cellNode := document.Call("createElement", "cell")
			cellNode.Call("setAttribute", "alive", cell.alive)
			cellNode.Set("id", fmt.Sprintf("cell-%d-%d", x, y))
			// TODO: move away from pixels and use something more dynamic
			cellNode.Set("style", fmt.Sprintf("width: %dpx; height: %dpx; border: solid %dpx", CELLSIZE, CELLSIZE, CELLBORDERSIZE))
			cellNode.Call("addEventListener", "click", cellClickCb)
			rowNode.Call("appendChild", cellNode)
		}
		g.dom.Grid.Call("appendChild", rowNode)
	}

	g.dom.msgContainer = document.Call("getElementById", "messageContainer")
	g.dom.msgContainer.Set("innerHTML", "") // clear out message container

	g.dom.genLabel = document.Call("createElement", "generation-label")
	g.dom.genLabel.Set("innerHTML", "Generation: 0")

	g.dom.generation = document.Call("createElement", "generation")
	g.dom.generation.Set("innerHTML", g.generation)

	g.logger("message: ", g.message)
	g.dom.message = document.Call("createElement", "message")
	g.dom.message.Set("innerHTML", g.message)

	g.dom.msgContainer.Call("appendChild", g.dom.genLabel)
	g.dom.msgContainer.Call("appendChild", g.dom.generation)
	g.dom.msgContainer.Call("appendChild", g.dom.message)
}

func (g *Game) updateDOMGrid() {
	document := js.Global().Get("document")

	// TODO: this would be better if we only change what needs to be changed
	for y, row := range g.cells {
		for x, cell := range row {
			domCell := document.Call("getElementById", fmt.Sprintf("cell-%d-%d", x, y))
			domCell.Call("setAttribute", "alive", cell.alive)
			if SHOWNEIGHBORS {
				domCell.Set("innerHTML", cell.neighbors)
			}
		}
	}
}

func (g *Game) renderMessage() {
	g.dom.genLabel.Set("innerHTML", "Generation: ")
	g.dom.generation.Set("innerHTML", g.generation)
	g.dom.message.Set("innerHTML", g.message)
}

func (g *Game) generate() {
	tmpCells := g.cells

	g.updateCells(g.countNeighbors)
	g.updateCells(g.checkRules)
	if reflect.DeepEqual(tmpCells, g.cells) {
		g.message = "graph did not change - pausing..."
		g.isPaused = false
		g.render()
		return
	}
	g.generation++
	g.render()
}

// startPolling creates an infinite loop which is important because it prevents the go code from exiting
func (g *Game) startPolling() {
	for {
		// TODO: only update this when it changes
		if !g.isPaused {
			go g.generate()
			g.dom.PlayBtn.Set("innerHTML", "Pause")
		} else {
			g.dom.PlayBtn.Set("innerHTML", "Play")
		}
		time.Sleep(g.speed * time.Second)
	}
}

func (g *Game) togglePause() {
	g.isPaused = !g.isPaused
}

func main() {
	game := &Game{}
	game.initGame()
}
