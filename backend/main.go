package main

import (
	"log"
	"meet-me/auth"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	auth.Github_Routes(mux)

	if err := http.ListenAndServe(":8000", mux); err != nil {
		log.Fatal(err)
	}
}
