package main

import (
	"net/http"
)

var users = []User{
	User{1, "Jacopo", "jacopo@gmail.com"},
}

func serveHTTP() {
	mux := http.NewServeMux()
	mux.Handle("/", routeHandler{routes})
	http.ListenAndServe(":8000", mux)
}

func main() {
	serveHTTP()
}

/* TODO
*  - TODOs left
*  - use a DB
 */
