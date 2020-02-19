package life

import (
	"reflect"
	"time"
)

type cell struct {
	Alive       bool
	Neighbors   int
	coordinates coords
}

type coords struct {
	x int
	y int
}

// TODO: add canvas support

// Game holds the state of our cells
type Game struct {
	IsPaused   bool
	Generation int
	Message    string
	speed      time.Duration
	Cells      [SIZE][SIZE]cell
	RenderCb   func()
	ResetCb    func()
}

// New Returns a new Game instance
func New() Game {
	g := &Game{}
	g.speed = TICKSPEED
	return *g
}

func (g *Game) Init() {

	g.reset()
	g.startPolling()

}

func (g *Game) reset() {
	g.ResetCb()
	g.IsPaused = true
	g.Generation = 0
	g.Message = ""
	g.updateCells(randomAlive)
	g.RenderCb()
	logger("reset!")
}

func (g *Game) ClearGrid() {

	g.updateCells(func(c cell, xy coords) cell {
		c.Alive = false
		g.countNeighbors(c, xy)
		return c
	})

	g.IsPaused = true
	g.Generation = 0
	g.Message = "cleared"
	logger("cleared")
	g.RenderCb()
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
			if g.Cells[nX][nY].Alive {
				neighborCount++
			}
		}
	}
	c.Neighbors = neighborCount
	return c
}

func (g *Game) checkRules(c cell, xy coords) cell {
	if c.Alive {
		if c.Neighbors < 2 || c.Neighbors > 3 {
			c.Alive = false
		}
	} else if c.Neighbors == 3 {
		c.Alive = true
	}
	return c
}

// iterate steps through the graph and modifies the
// cell based on the operation passed in
func (g *Game) updateCells(changeCell func(c cell, xy coords) cell) {
	var i, j int
	for i = 0; i < SIZE; i++ {
		for j = 0; j < SIZE; j++ {
			g.Cells[i][j] = changeCell(g.Cells[i][j], coords{i, j})
		}
	}
}

func randomAlive(c cell, xy coords) cell {
	c.Alive = d2(int64(xy.x + xy.y))
	return c
}

func (g *Game) Generate() {
	tmpCells := g.Cells

	g.updateCells(g.countNeighbors)
	g.updateCells(g.checkRules)
	if reflect.DeepEqual(tmpCells, g.Cells) {
		g.Message = "graph did not change - pausing..."
		g.IsPaused = false
		g.RenderCb()
		return
	}
	g.Generation++
	g.RenderCb()
}

// startPolling creates an infinite loop which is important because it prevents the go code from exiting
func (g *Game) startPolling() {
	for {
		if !g.IsPaused {
			go g.Generate()
		}

		time.Sleep(g.speed * time.Second)
	}
}

func (g *Game) TogglePause() {
	g.IsPaused = !g.IsPaused
}
