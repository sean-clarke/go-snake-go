package main

import (
	"fmt"
)

var (
	randomModifier float32 = 1 // range of 0-2 determining entropy's effect on snake's movement
	hungerModifier float32 = 1 // rating of 0-2 determining hunger's effect on snake's movement
	fearModifier float32 = 1 // rating of 0-2 determining pessimism's effect on snake's movement
	directionShift = map[Direction][]int {
		"up": {0, -1},
		"left" : {-1, 0},
		"right" : {1, 0},
		"down" : {0, 1},
	}
)

/*
	evalBoard
	paramaters:
		*

	returns:
		*
*/
func (matrix Matrix) initMatrix(data Req) {
	var width, height int = data.Board.Width, data.Board.Height
	for x := range matrix.Matrix {
		for y := range matrix.Matrix[x] {
			var t string = "O"
			var c int = 0
			var f bool = false
			var v float32 = 1
			if x == 0 || x == width - 1 {
				v -= 0.25
			}
			if y == 0 || y == height - 1 {
				v -= 0.25
			}

			matrix.Matrix[x][y] = Square{
				X: x,
				Y: y,
				Tag: t,
				Tenure: c,
				Food: f,
				Base: v,
			}
		}
	}
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
func rateSquare(x int, y int, width int, height int, origin Direction, distance int, depth int, matrix Matrix) float32 {
	if x == -1 || x == width || y == -1 || y == height {
		return 0
	}
	if depth == 0 {
		return matrix.Matrix[x][y].Base
	}

	var nodes = map[Direction]Direction{
		"up": "down",
		"left": "right",
		"right": "left",
		"down": "up",
	}
	delete(nodes, origin)

	var rating float32 = 0

	for direction, opposite := range nodes {
		rating += rateSquare(
			x + directionShift[direction][0],
			y + directionShift[direction][1],
			width,
			height,
			opposite,
			distance + 1,
			depth - 1,
			matrix,
		)
	}

	return rating
}

/*
	step
	paramaters:
		game Req

	returns:
		direction Direction
*/
func step(data Req) Direction {
	var matrix = Matrix{ make([][]Square, data.Board.Width) }
	for i := range matrix.Matrix {
    	matrix.Matrix[i] = make([]Square, data.Board.Height)
	}

	matrix.initMatrix(data)

	for w := range matrix.Matrix {
		for h := range matrix.Matrix[w] {
			fmt.Printf("%.2f ", matrix.Matrix[w][h].Base)
		}
		fmt.Printf("%s", "\n")
	}
	fmt.Println("")

	var directions = map[Direction]Direction{
		"up": "down",
		"left": "right",
		"right": "left",
		"down": "up",
	}
	var x0, x1, y0, y1 int = data.You.Body[0].X, data.You.Body[1].X, data.You.Body[0].Y, data.You.Body[1].Y
	
	if x0 < x1 {
		delete(directions, "right")
	} else if x0 > x1 {
		delete(directions, "left")
	} else if y0 < y1 {
		delete(directions, "down")
	} else if y0 > y1 {
		delete(directions, "up")
	}

	var next Direction
	var confidence float32 = 0

	for direction, opposite := range directions {
		var c = rateSquare(
			x0 + directionShift[direction][0],
			y0 + directionShift[direction][1],
			data.Board.Width,
			data.Board.Height,
			opposite,
			1, // maybe should be 0
			3,
			matrix,
		)
		if c > confidence {
			next = direction
			confidence = c
		}
	}

	return next
}