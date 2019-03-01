package main

type Coordinate struct {
	Y int `json:"y"`
	X int `json:"x"`
}

type Position struct {
	Y int
	X int
}

type Head struct {
	Pos Position
	Length int
}

type Square struct {
	Tenure int
	Danger int
	Food bool
	Self bool
	Base float64
}

type Matrix struct {
	Matrix [][]Square
	Width int
	Height int
	Heads []Head
	Food []Position
}

type Direction string

const (
	Up Direction = "up"
	Left Direction = "left"
	Right Direction = "right"
	Down Direction = "down"
)

type Rating struct {
	Value float64
	Distance int
}

type Packet struct {
	Dir Direction
	Rating Rating
}

type Snake struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Health int    `json:"health"`
	Body   []Coordinate  `json:"body"`
}

type Board struct {
	Height int     `json:"height"`
	Width  int     `json:"width"`
	Food   []Coordinate   `json:"food"`
	Snakes []Snake `json:"snakes"`
}

type Game struct {
	ID string `json:"id"`
}

type Req struct {
	Game  Game  `json:"game"`
	Turn  int   `json:"turn"`
	Board Board `json:"board"`
	You   Snake `json:"you"`
}

type Init struct {
	Color string `json:"color,omitempty"`
	Head string `json:"headType,omitempty"` 
	Tail string `json:"tailType,omitempty"`
}

type Resp struct {
	Move Direction `json:"move"`
}