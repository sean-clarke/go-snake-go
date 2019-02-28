package main


import (
	"fmt"
	"math/rand"
	"sync"
)

// global variables
var (
	globalDepth int = 13 // maximum iterations of rateSquare the snake will attempt 
	randomModifier float64 = 0.1 // range of 0-2 determining entropy's effect on snake's movement
	hungerModifier float64 // rating of 0-2 determining hunger's effect on snake's movement
	fearModifier float64 = 1 // rating of 0-2 determining pessimism's effect on snake's movement
	directionShift = map[Direction][]int {
		Up: {-1, 0},
		Left: {0, -1},
		Right: {0, 1},
		Down: {1, 0},
	}
)



/*	printMatrix
*		Prints either a value grid or board representation of the matrix
*	type:
*		Matrix
*	parameters:
*		bool - true prints board representation, false prints square values
*	returns:
*		nil
*/
func (matrix Matrix) printMatrix(repr bool) {
	if repr {
		for y := range matrix.Matrix {
			for x := range matrix.Matrix[y] {
				if matrix.Matrix[y][x].Food {
					fmt.Printf("%s ", "F")
				} else if matrix.Matrix[y][x].Tenure > 9 {
					fmt.Printf("%d", matrix.Matrix[y][x].Tenure)
				} else {
					fmt.Printf("%d ", matrix.Matrix[y][x].Tenure)
				}
			}
			fmt.Printf("%s", "\n")
		}
		fmt.Println("")
		return
	}
	for y := range matrix.Matrix {
		for x := range matrix.Matrix[y] {
			fmt.Printf("%.2f    ", matrix.Matrix[y][x].Base)
		}
		fmt.Printf("%s", "\n")
		fmt.Println("")
		fmt.Println("")
	}
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
}

/*
*	populateMatrix
*		Sets the food and tenure values of the matrix's squares
*	type:
*		Matrix
*	parameters:
*		Req
*	returns:
*		nil
*/
func (matrix Matrix) populateMatrix(data Req) {
	// set food
	for i := range data.Board.Food {
		food := &data.Board.Food[i]
		matrix.Matrix[food.Y][food.X].Food = true
	}

	// set tenure / matrix's heads
	for i := range data.Board.Snakes {
		snake := &data.Board.Snakes[i]
		id := snake.ID
		head := snake.Body[0]
		length := len(snake.Body)

		if id != data.You.ID {
			matrix.Heads = append(matrix.Heads, Position{head.Y, head.X})
		}
		matrix.Matrix[head.Y][head.X].Tenure = length - 1

		for p := range snake.Body[1:length] {
			tail := &snake.Body[p]
			self := id == data.You.ID
			matrix.Matrix[tail.Y][tail.X].Tenure = length - 1 - p
			if self {
				matrix.Matrix[tail.Y][tail.X].Self = self
			}
		}
	}
} 

/*
*	initMatrix
*		Creates the matrix's squares and call's populateMatrix
*	type:
*		Matrix
*	paramaters:
*		Req
*	returns:
*		nil
*/
func (matrix Matrix) initMatrix(data Req) {
	var width, height int = data.Board.Width, data.Board.Height
	for y := range matrix.Matrix {
		for x := range matrix.Matrix[y] {
			var v float64 = 1

			// Give edge & corner squares a lower base value 
			if x == 0 || x == width - 1 {
				v -= 0.25
			}
			if y == 0 || y == height - 1 {
				v -= 0.25
			}

			// Add random value to square value
			if randomModifier != 0 {
				v += rand.Float64() * randomModifier
			}

			matrix.Matrix[y][x] = Square{
				Tenure: 0,
				Danger: 0,
				Food: false,
				Self: false,
				Base: v,
			}
		}
	}
	matrix.populateMatrix(data)
}

/*
*	rateSquare
*	paramaters:
*		origin Direction
*		depth int
*
*	returns:
*		rating int
*
*	description:
*		lorem ipsum
*/
func (matrix Matrix) rateSquare(y int, x int, origin Direction, distance int, depth int, length int, grownby int, history []Position) float64 {
	// return 0 if out of bounds
	if x == -1 || x == matrix.Width || y == -1 || y == matrix.Height {
		return 0
	}
	// return 0 if occupied square
	eatenOffset := 0
	if matrix.Matrix[y][x].Self {
		eatenOffset = grownby
	}
	if matrix.Matrix[y][x].Tenure + eatenOffset >= distance {
		return 0
	}
	for p := range history {
		pos := &history[p]
		if x == pos.X {
			if y == pos.Y {
				return 0
			}
		}
	}
	history = append(history, Position{y, x})
	if len(history) >= length + grownby {
		history = history[1:]
	}

	base := matrix.Matrix[y][x].Base
	if matrix.Matrix[y][x].Food {
		grownby += 1
		base += float64(100 / (distance * distance)) * hungerModifier
	}

	// base case
	if depth == 0 {
		return base
	}

	// recursive case
	var nodes = map[Direction]Direction{
		Up: Down,
		Left: Right,
		Right: Left,
		Down: Up,
	}
	delete(nodes, origin)

	var rating float64 = 0

	for direction, opposite := range nodes {
		value:= matrix.rateSquare(
			y + directionShift[direction][0],
			x + directionShift[direction][1],
			opposite,
			distance + 1,
			depth - 1,
			length,
			grownby,
			history,
		)
		rating += base * value / 3
	}

	return rating
}

/*
*	step
*	paramaters:
*		Req
*	returns:
*		Direction
*/
func step(data Req) Direction {
	var matrix = Matrix{
		make([][]Square, data.Board.Height),
		data.Board.Width,
		data.Board.Height,
		[]Position{},
		[]Position{},
	}
	var allocation = make([]Square, matrix.Width * matrix.Height)
	for i := range matrix.Matrix {
    	matrix.Matrix[i] = allocation[i*matrix.Width: (i+1)*matrix.Width]
	}

	// fmtfmt.Println()
	// fmt.Printf("%s: ", "turn")
	// fmt.Printf("%d", data.Turn)
	// fmt.Println()
	matrix.initMatrix(data)
	// matrix.printMatrix(false) // print matrix values
	// matrix.printMatrix(true) // print matrix object repr

	var directions = map[Direction]Direction{
		Up: Down,
		Left: Right,
		Right: Left,
		Down: Up,
	}
	var x0, x1, y0, y1 int = data.You.Body[0].X, data.You.Body[1].X, data.You.Body[0].Y, data.You.Body[1].Y
	
	if x0 < x1 {
		delete(directions, Right)
	} else if x0 > x1 {
		delete(directions, Left)
	} else if y0 < y1 {
		delete(directions, Down)
	} else if y0 > y1 {
		delete(directions, Up)
	}

	var next Direction
	var confidence float64 = 0

	if data.You.Health > 80 {
		hungerModifier = 0
	} else if data.You.Health > 60 {
		hungerModifier = 0.25
	} else if data.You.Health > 40 {
		hungerModifier = 0.5
	} else if data.You.Health > 20 {
		hungerModifier = 1
	} else {
		hungerModifier = 2
	}

	length := len(data.You.Body)

	// limit depth by snake length
	localDepth := globalDepth
	if localDepth > length + 5 {
		localDepth = length + 5
	} else if length > 100 {
		localDepth = 16
	} else if length > 80 {
		localDepth = 15
	} else if length > 60 {
		localDepth = 14
	}

	// fmt.Printf("%s: ", "depth")
	// fmt.Printf("%d", localDepth)
	// fmt.Println()

	var wg sync.WaitGroup
	for direction, opposite := range directions {
		wg.Add(1)
		go func(y int, x int, direction Direction, opposite Direction, depth int, length int) {
			var c = matrix.rateSquare(
				y + directionShift[direction][0],
				x + directionShift[direction][1],
				opposite,
				1,
				depth,
				length,
				0,
				[]Position{},
			)
			defer wg.Done()
			// fmt.Printf("%s: ", direction)
			// fmt.Printf("%8f", c)
			// fmt.Println()
			if c > confidence {
				next = direction
				confidence = c
			}
		}(y0, x0, direction, opposite, localDepth, length)
	}
	wg.Wait()

	return next
}