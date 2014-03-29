package webui

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func WebHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Incoming web request")

	data, err := ioutil.ReadFile("src/webui/default.html")

	if err != nil {
		panic(err)
	}

	w.Write(data)
}
