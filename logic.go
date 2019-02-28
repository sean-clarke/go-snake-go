package main


import (
	// "fmt"
	"math/rand"
)

// global variables
var (
	hungerModifier float64 // rating of 0-2 determining hunger's effect on snake's movement
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
*//*
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
*/
/*
*	populateMatrix
*		Sets the food and tenure values of the matrix's squares
*	type:
*		matrix Matrix
*	parameters:
*		data Req
*	returns:
*		nil
*	description:
*		lorem ipsum
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
*		matrix Matrix
*	paramaters:
*		data Req
*	returns:
*		nil
*	description:
*		lorem ipsum
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

			// Initialize randomModifier
			var randomModifier float64 = 0.1

			// Increase square value by random value if randomModifier > 0
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
*	exp
*	parameters:
*		base float64
*		power float 64
*	returns:
*		result float64
*	description:
*		lorem ipsum
*/
func exp(base float64, power float64) float64 {
	if power == 0 {
		return 1
	}
	for power > 1 {
		base *= base
		power -= 1
	}
	return base
}

/*
*	rateSquare
*	paramaters:
*		origin Direction
*		depth int
*	returns:
*		rating int
*	description:
*		lorem ipsum
*/
func (matrix *Matrix) rateSquare(y int, x int, origin Direction, distance int, depth int, length int, grownby int, health int, history []Position) Rating {

	// out of bounds
	if x == -1 || x == matrix.Width || y == -1 || y == matrix.Height {
		return Rating{0, distance}
	}

	// currently occupied
	eatenOffset := 0
	if matrix.Matrix[y][x].Self {
		eatenOffset = grownby
	}
	if matrix.Matrix[y][x].Tenure + eatenOffset >= distance {
		return Rating{0, distance}
	}

	// occupied by current path
	for p := range history {
		pos := &history[p]
		if x == pos.X {
			if y == pos.Y {
				return Rating{0, distance}
			}
		}
	}

	// set base value
	health -= 1
	base := matrix.Matrix[y][x].Base
	if matrix.Matrix[y][x].Food {
		grownby += 1
		// to promote moderation, 25 <-> 20, 4 <-> 2
		var hungerModifier float64 = 4 / (exp(2, float64(health) / 25))
		base += float64(100 / (distance * distance)) * 4 * hungerModifier
		health = 100
	}

	// return base value (base case)
	if depth == 0 {
		return Rating{base, distance}
	}

	// add current position to history
	history = append(history, Position{y, x})

	// remove last position in history if tenure is up
	if length < depth && len(history) >= length + grownby {
		history = history[1:]
	}

	// continue search (recursive case)
	var directions = map[Direction]Direction{
		Up: Down,
		Left: Right,
		Right: Left,
		Down: Up,
	}
	delete(directions, origin)

	var rating Rating

	for direction, opposite := range directions {
		node := matrix.rateSquare(
			y + directionShift[direction][0],
			x + directionShift[direction][1],
			opposite,
			distance + 1,
			depth - 1,
			length,
			grownby,
			health,
			history,
		)
		rating.Value += base * node.Value / 3
		if node.Distance > rating.Distance {
			rating.Distance = node.Distance
		}
	}

	return rating
}

/*
*	step
*	paramaters:
*		Req
*	returns:
*		Direction
*	description:
*		lorem ipsum
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

	length := len(data.You.Body)

	// limit depth by snake length
	var localDepth int = 12 // maximum iterations of rateSquare the snake will attempt 
	if length < 50 {
		if localDepth > length + 2 {
			localDepth = length + 2
		}
	} else {
		localDepth += (length - 30) / 18
	}

	// fmt.Printf("%s: ", "depth")
	// fmt.Printf("%d", localDepth)
	// fmt.Println()

	ch := make(chan Packet)
	defer close(ch)
	for direction, opposite := range directions {
		go func(direction Direction, opposite Direction) {
			var rating = matrix.rateSquare(
				y0 + directionShift[direction][0],
				x0 + directionShift[direction][1],
				opposite,
				1,
				localDepth,
				length,
				0,
				data.You.Health,
				[]Position{},
			)
			ch <- Packet{direction, rating}
		}(direction, opposite)	
	}

	reach := 0

	for i := 0; i < len(directions); i++ {
		packet := <-ch
		// fmt.Printf("%s: ", packet.Dir)
		// fmt.Printf("%8f", packet.Rating.Value)
		// fmt.Println()
		// choose highest rated path
		if packet.Rating.Value > confidence {
			next = packet.Dir
			confidence = packet.Rating.Value
		// fallback on longest path if death is inevitable
		} else if confidence == 0 && packet.Rating.Distance > reach {
			reach = packet.Rating.Distance
			next = packet.Dir
		}
	}

	return next
}