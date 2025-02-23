package main

import (
	"fmt"
	"log"
	"net/http"
)

type Operands struct {
	Operand1 string `json:"a"`
	Operand2 string `json:"b"`
}

func main() {
	// only on test.
	fmt.Println("backend calculator using golan.")

	//
	server := http.NewServeMux()
	port := ":5000"

	server.HandleFunc("GET /add", Add)
	server.HandleFunc("/sub", Sub)
	server.HandleFunc("/multi", Multi)
	server.HandleFunc("/div", Div)

	fmt.Println("server...")
	log.Fatal(http.ListenAndServe(port, server))

}
