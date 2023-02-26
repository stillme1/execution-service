package main

import (
	"log"
	"net/http"

	"./run"
)

func main() {

	http.HandleFunc("/submit", run.Submit)
	log.Println("Listeing on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
