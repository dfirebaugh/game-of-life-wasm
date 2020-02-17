package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"syscall/js"
	"time"
)

// SIZE size of the grid
//- i.e. the width and height of the grid
const SIZE = 10

// TICKSPEED is how fast the generations update
const TICKSPEED = 1

type cell struct {
	alive bool
}

type coords struct {
	x int
	y int
}

// Game holds the state of our cells
type Game struct {
	b          bool
	isPaused   bool
	generation int
	message    string
	printlog   bool
	speed      time.Duration
	cells      [SIZE][SIZE]cell
}

// initBoard - setup game
func (g *Game) initGame() {
	js.Global().Set("BOARD_SIZE", SIZE)
	g.printlog = true
	g.isPaused = true
	g.speed = TICKSPEED

	g.reset()

	g.registerCallbacks()
}

func (g *Game) reset() {
	g.generation = 0
	g.updateCells(randomAlive)
	g.logger("reset!")

	g.boardToJS()
	chunk(g.toBinaryStr(true), SIZE)
	g.print()
	g.printToDOM()
}

// registerCallbacks defines callbacks that are used within JS code
// A loop is created to prevent the go code from exiting
func (g *Game) registerCallbacks() {
	var runCb, generateCb js.Func
	runCb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		g.boardToJS()
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

	// wait := make(chan struct{}, 0)
	js.Global().Get("document").Call("getElementById", "runBtn").Call("addEventListener", "click", runCb)
	js.Global().Get("document").Call("getElementById", "generateBtn").Call("addEventListener", "click", generateCb)
	js.Global().Get("document").Call("getElementById", "resetBtn").Call("addEventListener", "click", resetCb)
	g.startPolling()
	// <-wait
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

// toBinaryStr - flatten the cells and convert to a string of binary numbers
// because apparently sharing values to JS of certain types is difficult
//
// If the pretty option is true, the return will include line endings based on the boards SIZE
func (g *Game) toBinaryStr(pretty bool) string {

	var acc []string
	for _, row := range g.cells {
		for _, cell := range row {
			var bitSetVar int8
			if cell.alive == true {
				bitSetVar = 1
			}
			acc = append(acc, strconv.FormatInt(int64(bitSetVar), 2))
		}
	}

	binaryStr := strings.Trim(strings.Join(acc[:], ""), " ")

	if pretty {
		return strings.Join(chunk(binaryStr, SIZE), "\n")
	}

	return binaryStr
}

func (g *Game) boardToJS() {
	js.Global().Set("currentBoard", g.toBinaryStr(false))
}

func (g *Game) checkRules(c cell, xy coords) cell {
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

	if c.alive {
		if neighborCount < 2 || neighborCount > 3 {
			c.alive = false
		}
	} else if neighborCount == 3 {
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

func (g *Game) print() {
	g.logger(g.toBinaryStr(true))
}

func (g *Game) printToDOM() {
	grid := js.Global().Get("document").Call("getElementById", "gridContainer")

	grid.Set("innerHTML", g.toBinaryStr(true))
}

func (g *Game) printMessage() {
	document := js.Global().Get("document")

	genLabel := document.Call("getElementById", "generation-label")
	genLabel.Set("innerHTML", "Generation: ")

	gen := document.Call("getElementById", "generation")
	gen.Set("innerHTML", g.generation)

	if len(g.message) > 0 {
		g.logger("message: ", g.message)
		msg := document.Call("getElementById", "message")
		msg.Set("innerHTML", g.message)
	}
}

func (g *Game) generate() {
	tmpCells := g.cells

	g.updateCells(g.checkRules)
	if reflect.DeepEqual(tmpCells, g.cells) {
		g.message = "graph did not change - pausing..."
		g.isPaused = false
		g.printMessage()
		return
	}
	g.generation++
	g.printMessage()
	g.print()
	g.printToDOM()
}

func (g *Game) startPolling() {
	btn := js.Global().Get("document").Call("getElementById", "runBtn")
	for {
		if !g.isPaused {
			go g.generate()
			btn.Set("innerHTML", "Pause")
		} else {
			btn.Set("innerHTML", "Play")
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
