package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	fmt.Println("backend calculator using golan.")

	//
	server := http.NewServeMux()
	port := ":5000"

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi from the server"))
	})

	log.Fatal(http.ListenAndServe(port, server))

}
