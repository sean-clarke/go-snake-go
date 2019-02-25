package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/start", Start)
	http.HandleFunc("/move", Move)
	http.HandleFunc("/end", End)
	http.HandleFunc("/ping", Ping)

	http.ListenAndServe(":8080", nil)
}
