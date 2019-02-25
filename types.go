package main

type Pos struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Square struct {
	Value int
}

type Direction string

const (
	Up Direction = "up"
	Left Direction = "left"
	Down Direction = "down"
	Right Direction = "right"
)

type Snake struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Health int    `json:"health"`
	Body   []Pos  `json:"body"`
}

type Board struct {
	Height int     `json:"height"`
	Width  int     `json:"width"`
	Food   []Pos   `json:"food"`
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