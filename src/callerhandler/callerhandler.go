package callerhandler

import (
	"bytes"
	"encoding/xml"
	"net/http"
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
	b := bytes.NewBufferString(start)
	response := &Say{Voice: "female", Language: "en", Loop: 0, Text: "Colin"}

	str, err := xml.Marshal(response)
	if err != nil {
		panic(err)
	}
	b.Write(str)
	b.WriteString(end)
	b.WriteTo(w)
}
