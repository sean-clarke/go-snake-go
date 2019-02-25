package main

import (
	"net/http"
	"encoding/json"
)

func Index(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("go-snake-go"))
}

func Start(res http.ResponseWriter, req *http.Request) {
	decoded := Req{}
	err := json.NewDecoder(req.Body).Decode(&decoded)
	if err != nil {
		return
	}

	resp := Init{
		Color: "#006666",
		Head: "pixel",
		Tail: "pixel",
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(resp)
	res.Write([]byte("\n"))
}

func Move(res http.ResponseWriter, req *http.Request) {
	decoded := Req{}
	err := json.NewDecoder(req.Body).Decode(&decoded)
	if err != nil {
		return
	}

	direction := step(decoded)

	resp := Resp{Move: direction}
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(resp)
	res.Write([]byte("\n"))
}

func End(res http.ResponseWriter, req *http.Request) {
	return
}

func Ping(res http.ResponseWriter, req *http.Request) {
	return
}