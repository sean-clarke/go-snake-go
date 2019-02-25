package main

var (
	randomModifier int = 1
	hungerModifier int = 1
)


/*
	initBoard
	paramaters:
		*

	returns:
		*
*/
func initBoard(height int, width int) {

	return
}

/*
	rateSquare
	paramaters:
		depth int

	returns:
		rating int
*/
func rateSquare(depth int) int {
	return 0
}

/*
	step
	paramaters:
		game Req

	returns:
		direction Direction
*/
func step(game Req) Direction {

	if game.You.Body[0].Y == 8 {
		if game.You.Body[0].X == 8 {
			return Left
		}

		return Right
	}

	return Down
}