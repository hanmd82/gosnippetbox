package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// define a struct to hold application-wide dependencies
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// note that flag.String returns a pointer to the value, not the value itself
	addrPtr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// initialise a new instance of application struct
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// use the application struct's methods as the handler functions
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	srv := &http.Server{
		Addr:     *addrPtr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Printf("Starting server on %s", *addrPtr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
