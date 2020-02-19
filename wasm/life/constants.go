package life

// PRINTLOG - should we print the console messages?
const PRINTLOG = false

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
