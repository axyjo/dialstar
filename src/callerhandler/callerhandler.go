package callerhandler

import (
	"bytes"
	"encoding/xml"
	_ "fmt"
	_ "io/ioutil"
	"net/http"
	"twiml"
)

const (
	start = `<?xml version="1.0" encoding="UTF-8"?><Response>`
	end   = `</Response>`
)

type Say struct {
	XMLName  xml.Name `xml:"Say"`
	Voice    string   `xml:"voice,attr"`
	Language string   `xml:"language,attr"`
	Loop     uint     `xml:"loop,attr"`
	Text     string   `xml:",chardata"`
}
type context struct {
	b *bytes.Buffer
	r *http.Request
}

func CallerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	cityName := r.Form["FromCity"]
	//fmt.Println(actual)
	b := bytes.NewBufferString(start)
	say_response := &Say{Voice: "female", Language: "en", Loop: 1, Text: "Colin from " + cityName[0]}

	str, err := xml.Marshal(say_response)
	if err != nil {
		panic(err)
	}
	b.Write(str)
	response := &twiml.Conference{Text: "foobar", EndConferenceOnExit: "true"}
	dial := &twiml.Dial{Conference: *response}
	str, err = xml.Marshal(dial)
	if err != nil {
		panic(err)
	}
	b.Write(str)
	b.WriteString(end)
	b.WriteTo(w)
}
