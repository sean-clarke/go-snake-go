package main


import (
	"fmt"
	"math/rand"
)


var (
	randomModifier float32 = 1 // range of 0-2 determining entropy's effect on snake's movement
	hungerModifier float32 = 1 // rating of 0-2 determining hunger's effect on snake's movement
	fearModifier float32 = 1 // rating of 0-2 determining pessimism's effect on snake's movement
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
				} else {
					fmt.Printf("%d ", matrix.Matrix[y][x].Tenure)
				}
			}
			fmt.Printf("%s", "\n")
		}
		fmt.Println("")
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
	for idx := range data.Board.Food {
		food := &data.Board.Food[idx]
		matrix.Matrix[food.Y][food.X].Food = true
	}

	// set tenure / matrix's heads
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
			var c int = 0
			var f bool = false
			var v float32 = 1

			// Give edge & corner squares a lower base value 
			if x == 0 || x == width - 1 {
				v -= 0.25
			}
			if y == 0 || y == height - 1 {
				v -= 0.25
			}

			// Add random value to square value
			if randomModifier != 0 {
				v += rand.Float32() * randomModifier
			}

			matrix.Matrix[y][x] = Square{
				Tenure: c,
				Food: f,
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
func (matrix Matrix) rateSquare(y int, x int, origin Direction, distance int, depth int) float32 {
	if x == -1 || x == matrix.Width || y == -1 || y == matrix.Height {
		return 0
	}
	if depth == 0 {
		return matrix.Matrix[y][x].Base
	}

	var nodes = map[Direction]Direction{
		Up: Down,
		Left: Right,
		Right: Left,
		Down: Up,
	}
	delete(nodes, origin)

	var rating float32 = 0

	for direction, opposite := range nodes {
		rating += matrix.Matrix[y][x].Base * matrix.rateSquare(
			y + directionShift[direction][0],
			x + directionShift[direction][1],
			opposite,
			distance + 1,
			depth - 1,
		)
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

	matrix.initMatrix(data)
	matrix.printMatrix(false) // print matrix values
	matrix.printMatrix(true) // print matrix object repr

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
	var confidence float32 = 0

	for direction, opposite := range directions {
		var c = matrix.rateSquare(
			y0 + directionShift[direction][0],
			x0 + directionShift[direction][1],
			opposite,
			1, // maybe should be 0
			3,
		)
		if c > confidence {
			next = direction
			confidence = c
		}
	}

	return next
}