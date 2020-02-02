package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	// note that flag.String returns a pointer to the value, not the value itself
	addrPtr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Printf("Starting server on %s", *addrPtr)
	err := http.ListenAndServe(*addrPtr, mux)
	log.Fatal(err)
}
