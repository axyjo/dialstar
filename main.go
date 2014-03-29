package main

import (
	"caller_handler"
	"html/template"
	"net/http"
)

//This is just all the handlers and shit
func main() {
	http.HandleFunc("/caller/", caller_handler)
	http.ListenAndServe(":3000", nil)
}
