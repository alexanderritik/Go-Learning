package main

import (
	"fmt"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found", http.StatusNotFound)
	}

	if r.Method != "Get" {
		http.Error(w, "Method not supported", http.StatusNotFound)
	}

	fmt.Fprint(w, "hello!")
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/form" {
		http.Error(w, "404 not found", http.StatusNotFound)
	}

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm err: %v", err)
	}

	name := r.FormValue("name")
	fmt.Fprintf(w, "Name %s\n", name)
}

func main() {
	fileserver := http.FileServer(http.Dir("./static"))

	http.Handle("/", fileserver)

	http.HandleFunc("/form", formHandler)

	http.HandleFunc("/hello", helloHandler)

	fmt.Printf("Starting server at port 8080\n")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
