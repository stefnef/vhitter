package main

import (
	"log"
	"net/http"
)

func main() {
	index()
}

func index() {
	fs := http.FileServer(http.Dir("./ui/static"))
	http.Handle("/", fs)

	log.Println("Listening on :3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
