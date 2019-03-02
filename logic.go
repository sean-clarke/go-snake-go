package main

/*
*	exp
*		Standard exponent function
*	parameters:
*		base float64
*		power float 64
*	returns:
*		float64
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
*	expandDirections
*		expands encoded directions object into a slice of directions
*	parameters:
*		encoded int
*	returns:
*		[]Direction
*/
func expandDirections(encoded int) []Direction {
	var decoded []Direction
	if encoded % int(Up) == 0 {
		decoded = []Direction{Up}
	}
	if encoded % int(Left) == 0 {
		decoded = append(decoded, Left)
	}
	if encoded % int(Right) == 0 {
		decoded = append(decoded, Right)
	}
	if encoded % int(Down) == 0 {
		decoded = append(decoded, Down)
	}
	return decoded
}
/*
*	flip
*		Returns the flipped 
*	parameters:
*		start Position
*		dir Direction
*	returns:
*		Position
*/
func flip(dir Direction) Direction {
	switch dir {
	case Up:
		return Down
	case Left:
		return Right
	case Right:
		return Left
	case Down:
		return Up
	default:
		return Down
	}
}

/*
*	move
*		Gets the postition from start after moving in dir direction
*	parameters:
*		start Position
*		dir Direction
*	returns:
*		Position
*/
func move(start Position, dir Direction) Position {
	switch dir {
	case Up:
		return Position{start.Y - 1, start.X}
	case Left:
		return Position{start.Y, start.X - 1}
	case Right:
		return Position{start.Y, start.X + 1}
	case Down:
		return Position{start.Y + 1, start.X}
	default:
		return start
	}
}

/*
*	getNeighbours
*		Gets the postitions from start after moving in dir direction
*	parameters:
*		home Position
*		directions []Direction
*	returns:
*		[]Position
*/
func getNeighbours(home Position, directions []Direction, depth int, limit int) []Position {
	neighbours := []Position{}
	// out of bounds
	if home.X < 0 || home.X == limit || home.Y < 0 || home.Y == limit {} else if depth == 1 {
		for _, direction := range directions {
			neighbour := move(home, direction)
			if neighbour.X >= 0 && neighbour.X < limit && neighbour.Y >= 0 && neighbour.Y < limit {
				neighbours = append(neighbours, neighbour)
			}
		}
	} else if depth > 1 {
		// recursive case
		for _, direction := range directions {
			neighboursNeighbours := getNeighbours(move(home, direction), expandDirections(210 / int(flip(direction))), depth - 1, limit)
			for _, neighbour := range neighboursNeighbours {
				neighbours = append(neighbours, neighbour)
			}
		}
	}

	return neighbours
}

/*
*	rateSquare
*		Recursively rates a square by its child nodes and context in the game 
*	paramaters:
*		pos Position
*		origin Direction
*		distance int
*		depth int
*		length int
*		grownby int
*		health int
*		history []Position{int, int}
*	returns:
*		Rating{float64, int}
*/
func (matrix *Matrix) rateSquare(pos Position, origin Direction, distance int, depth int, length int, longest bool, grownby int, health int, history []Position) Rating {
	var y, x int = pos.Y, pos.X

	// out of bounds
	if x == -1 || x == matrix.Width || y == -1 || y == matrix.Height {
		return Rating{0, distance}
	}

	// forbidden move
	if matrix.Matrix[y][x].Base == 0 {
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
	for h := range history {
		past := &history[h]
		if x == past.X {
			if y == past.Y {
				return Rating{0, distance}
			}
		}
	}

	// set base value
	health -= 1
	base := matrix.Matrix[y][x].Base
	if matrix.Matrix[y][x].Food {
		grownby += 1
		var multiplier float64 = 5
		var divisor float64 = 33
		if longest && length > 8 {
			multiplier = 4
			divisor = 25
		}
		if matrix.Matrix[y][x].Danger < 2 {
			var hungerModifier float64 = multiplier / (exp(2, float64(health) / divisor))
			base += float64(100 / (distance * distance)) * multiplier * hungerModifier
		}
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
	var directions = expandDirections(210 / int(origin))
 
	var rating Rating

	// modify rating by rating of potential future moves
	for _, direction := range directions {
		node := matrix.rateSquare(
			move(pos, direction),
			flip(direction),
			distance + 1,
			depth - 1,
			length,
			longest,
			grownby,
			health,
			history,
		)
		node.Value -= 0.1 / (1 + float64(distance) * float64(distance))
		rating.Value += base * node.Value / 3 
		if node.Distance > rating.Distance {
			rating.Distance = node.Distance
		}
	}

	return rating
}

/*
*	step
*		Entry logic function that returns calculated approximate best next move
*	paramaters:
*		data Req
*	returns:
*		string
*/
func step(data Req) string {
	// set useful variables
	bWidth := data.Board.Width
	bHeight := data.Board.Height
	mId := data.You.ID
	mHead := Position{data.You.Body[0].Y, data.You.Body[0].X}
	mY, mX := mHead.Y, mHead.X
	mLength := len(data.You.Body)
	longest := true

	// move validation checker
	if mLength < 2 {
		return "up"
	}

	var directions []Direction = expandDirections(210)

	// create matrix
	var matrix = Matrix{
		make([][]Square, bHeight),
		bWidth,
		bHeight,
		[]Head{},
		[]Position{},
	}
	var allocation = make([]Square, bHeight * bWidth)
	for i := range matrix.Matrix {
    	matrix.Matrix[i] = allocation[i*bWidth: (i+1)*bWidth]
	}

	// initMatrix
	for y := range matrix.Matrix {
		for x := range matrix.Matrix[y] {
			var v float64 = 1
			var heatmap bool = true

			// Give edge & corner squares a lower base value (and )
			if x == 0 || x == bWidth - 1 {
				v -= 0.125
			} else if heatmap {
				if (y == 2 || y == bHeight - 3) {
					v += 0.33
				} else if (y == 1 || y == 3 || y == bHeight - 2 || y == bHeight - 4) {
					v += 0.16
				}
			}
			if y == 0 || y == bHeight - 1 {
				v -= 0.125
			} else if heatmap {
				if (x == 2 || x == bWidth - 3) {
					v += 0.33
				} else if (x == 1 || x == 3 || x == bWidth - 2 || x == bWidth - 4) {
					v += 0.16
				}
			}

			if mLength > 8 {
				if (x == 0 && y == bHeight / 2) || (x == bWidth - 1 && y == bHeight / 2) || (y == 0 && x == bWidth / 2) || (y == bHeight - 1 && x == bWidth / 2) {
					v += 1
				} else if (x == 0 && y == bHeight / 2 + 1) || (x == 0 && y == bHeight / 2 - 1) || (x == bWidth - 1 && y == bHeight / 2 + 1) || (x == bWidth - 1 && y == bHeight / 2 - 1) || (y == 0 && x == bWidth / 2 + 1) || (y == 0 && x == bWidth / 2 - 1) || (y == bHeight - 1 && x == bWidth / 2 + 1) || (y == bHeight - 1 && x == bWidth / 2 - 1) {
					v += 0.66
				}
			}

			// Initialize randomModifier
			var randomModifier float64 = 0

			// Increase square value by random value if randomModifier > 0
			if randomModifier != 0 {
				//v += rand.Float64() * randomModifier
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

	// populateMatrix
	for i := range data.Board.Food {
		food := &data.Board.Food[i]
		matrix.Matrix[food.Y][food.X].Food = true
	}

	// set tenure / matrix's heads
	for i := range data.Board.Snakes {
		snake := &data.Board.Snakes[i]
		oId := snake.ID
		oY := snake.Body[0].Y
		oX := snake.Body[0].X
		oLength := len(snake.Body)

		if oId != mId {
			matrix.Heads = append(matrix.Heads, Head{Position{oY, oX}, oLength})

			// generate squares next to head
			var pnd int = 210
			if (oY == 0) {
				pnd /= int(Up)
			}
			if oX == 0 {
				pnd /= int(Left)
			}
			if (oX == bWidth - 1) {
				pnd /= int(Right)
			}
			if (oY == bHeight - 1) {
				pnd /= int(Down)
			}
			neighbours := getNeighbours(Position{oY, oX}, expandDirections(pnd), 1, bWidth)
			dangers := getNeighbours(Position{oY, oX}, expandDirections(pnd), 2, bWidth)

			// for squares next to snakes heads...
			if oLength >= mLength {
				longest = false
				// ...if snake is larger than us, set base to ~0
				for neighbour := range neighbours {
					yard := &neighbours[neighbour]
					matrix.Matrix[yard.Y][yard.X].Base = 0.00001
					matrix.Matrix[yard.Y][yard.X].Danger = 2
				}
			} else if mLength > oLength {
				// ...if snake is smaller than us, set danger to -1
				for neighbour := range neighbours {
					yard := &neighbours[neighbour]
					matrix.Matrix[yard.Y][yard.X].Danger = -1
				}
			}

			// for squares two away from snakes heads, decrease base
			for d := range dangers {
				danger := &dangers[d]
				matrix.Matrix[danger.Y][danger.X].Base /= 2
				matrix.Matrix[danger.Y][danger.X].Danger = 1
			}
		} else if mLength > 8 {
			tail := &snake.Body[len(snake.Body) - 1]
			matrix.Matrix[tail.Y][tail.X].Base *= 1 + float64(mLength) / 10
		}
		matrix.Matrix[oY][oX].Tenure = oLength - 1

		for p := range snake.Body[1:oLength] {
			tail := &snake.Body[p]
			self := oId == mId
			matrix.Matrix[tail.Y][tail.X].Tenure = oLength - 1 - p
			if self {
				matrix.Matrix[tail.Y][tail.X].Self = self
			}
		}
	}

	// limit depth by snake length
	var localDepth int = 14
	if mLength < 50 {
		if localDepth > mLength + 2 {
			localDepth = mLength + 2
		}
	} else {
		localDepth += (mLength - 30) / 18
	}

	// concurrently rate potential moves
	ch := make(chan Packet)
	defer close(ch)
	for _, direction := range directions {
		go func(direction Direction) {
			var rating = matrix.rateSquare(
				move(Position{mY, mX}, direction),
				flip(direction),
				1,
				localDepth,
				mLength,
				longest,
				0,
				data.You.Health,
				[]Position{},
			)
			ch <- Packet{direction, rating}
		}(direction)	
	}

	// choose best move
	var next Direction
	var confidence float64 = 0
	var reach int = 0

	for i := 0; i < len(directions); i++ {
		packet := <-ch
		// prefer highest rated move
		if packet.Rating.Value > confidence {
			next = packet.Direction
			confidence = packet.Rating.Value
		// logic for longest path fallback if death is inevitable (ie. confidence == 0)
		} else if confidence == 0 && packet.Rating.Distance > reach {
			next = packet.Direction
			reach = packet.Rating.Distance
		}
	}

	switch next {
	case Up:
		return "up"
	case Left:
		return "left"
	case Right:
		return "right"
	case Down:
		return "down"
	default:
		return "up"
	}
}