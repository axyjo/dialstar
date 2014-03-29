package main

import (
	"callerhandler"
	"net/http"
)

//This is just all the handlers and shit
func main() {
	http.HandleFunc("/caller/", callerhandler.CallerHandler)
	http.ListenAndServe(":3000", nil)
}
