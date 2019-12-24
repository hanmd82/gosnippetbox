package main

import (
    "log"
    "net/http"
)

// Define a home handler function which writes a byte slice containing
// "Hello from Snippetbox" as the response body.
func home(w http.ResponseWriter, r *http.Request) {
    // Check if the current request URL path exactly matches "/". If it doesn't, use
    // the http.NotFound() function to send a 404 response to the client, and return.
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }

    w.Write([]byte("Hello from Snippetbox"))
}

// Add a showSnippet handler function.
func showSnippet(w http.ResponseWriter, r *http.Request) {
    w.Write([] byte("Display a specific snippet..."))
}

// Add a createSnippet handler function.
func createSnippet(w http.ResponseWriter, r *http.Request) {
    w.Write([] byte("Create a new snippet..."))
}

func main() {
    // Use the http.NewServeMux() function to initialize a new servemux, then
    // register the home function as the handler for the "/" URL pattern.
    mux := http.NewServeMux()
    mux.HandleFunc("/", home)

    // Register the two new handler functions and corresponding URL patterns with the servemux
    mux.HandleFunc("/snippet", showSnippet)
    mux.HandleFunc("/snippet/create", createSnippet)

    // Use the http.ListenAndServe() function to start a new web server, passing in
    // two parameters: the TCP network address to listen on (in this case ":4000")
    // and the servemux just created. If http.ListenAndServe() returns an error
    // use the log.Fatal() function to log the error message and exit. Any error
    // returned by http.ListenAndServe() is always non-nil.
    log.Println("Starting server on :4000")
    err := http.ListenAndServe(":4000", mux)
    log.Fatal(err)
}
