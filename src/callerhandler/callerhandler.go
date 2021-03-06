package callerhandler

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gorilla/schema"
	_ "io/ioutil"
	"net/http"
	"twiml"
)

//start and end of the xml sent to Twilio
const (
	start = `<?xml version="1.0" encoding="UTF-8"?><Response>`
	end   = `</Response>`
)

//CallerWrapper which holds a channel the interface Thingy
type CallerWrapper struct {
	Callerid chan twiml.Thingy
}

type context struct {
	b *bytes.Buffer
	r *http.Request
}

//Holds information about the user calling in
type VoiceRequest struct {
    CallSid       string `json:"CallSid"`
    AccountSid    string `json:"AccountSid"`
    From          string `json:"From"`
    To            string `json:"To"`
    CallStatus    string `json:"CallStatus"`
    ApiVersion    string `json:"ApiVersion"`
    Direction     string `json:"Direction"`
    ForwardedFrom string `json:"ForwardedFrom"`
    CallerName    string `json:"CallerName"`
    FromCity      string `json:"FromCity"`
    FromState     string `json:"FromState"`
    FromZip       string `json:"FromZip"`
    FromCountry   string `json:"FromCountry"`
    ToCity        string `json:"ToCity"`
    ToState       string `json:"ToState"`
    ToZip         string `json:"ToZip"`
    ToCountry     string `json:"ToCountry"`
}

func (c CallerWrapper) CallerHandler(w http.ResponseWriter, r *http.Request) {
	//Check if it's a POST call, if not return immediately
	if r.Method != "POST" {
		return
	}

	//Error handling
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	// Unmarshal Form data
	var request VoiceRequest
	decoder := schema.NewDecoder()
	decoder.Decode(&request, r.Form)

	// Store the information from the request into ElasticSearch for analytics
	bytesLine, err := json.Marshal(request)
	url := "http://twilio.axyjo.com:9200/hackathon/logs/"
	body := bytes.NewBuffer(bytesLine)
	_, err = http.Post(url, "application/json", body)
	if err != nil {
		panic(err)
	}

	if request.CallStatus == "completed" {
		return
	}

	//Creates a new Buffer with the initial start xml string
	b := bytes.NewBufferString(start)

	//Marshal the say_repsonse

	//say_response := &twiml.Say{Voice: "female", Language: "en", Loop: 1, Text: "Welcome to Dial Star, There are currently  " + cityName[0]}
	//str, err := xml.Marshal(say_response)
	////Error checking..
	//if err != nil {
	//	panic(err)
	//}
	//Append the Say block to the buffer
	//b.Write(str)
	//fmt.Println(r.Form["CallSid"][0] + " - " + " call initiated from " + r.Form["From"][0])
	//Play Cowbell until user is matched
	response := &twiml.Play{Text: "http://com.twilio.music.classical.s3.amazonaws.com/ClockworkWaltz.mp3", Loop: "1"}
	//Marshal the response
	str, err := xml.Marshal(response)
	if err != nil {
		panic(err)
	}
	//Write the reponse to the Buffer
	b.Write(str)
	//Write the end of the xml to the Buffer
	b.WriteString(end)
	//Write the Buffer to the http.ResponseWriter
	b.WriteTo(w)
	//adds the Thingy to the channel with the users's CallSid, City, and Queue flag set to true.
	c.Callerid <- twiml.Thingy{request.CallSid, request.FromCity, true, request.From}
	fmt.Println(request.CallSid + " - Queued")
}
