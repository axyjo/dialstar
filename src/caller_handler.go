package caller_handler

import (
	"net/http"
	"html/template"
)


func caller_handler(w. http.ResponseWriter, r *http.Request){
	t, _ := template.Parsefiles("names.xml")
	t.Execute(t, _)
}