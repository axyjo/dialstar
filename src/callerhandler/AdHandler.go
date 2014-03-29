package callerhandler

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/gorilla/schema"
	_ "io/ioutil"
	"net/http"
	"twiml"
)


var ad_counter = 0


func AdHandler(w http.ResponseWriter, r *http.Request) {
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
	//Store the city name of the user making the call
	//cityName := r.Form["FromCity"]
	//fmt.Println(actual)
	if request.CallStatus == "completed" {
		return
	}

	//Creates a new Buffer with the initial start xml string
	b := bytes.NewBufferString(start)

	// Queue an ad.
	ad_response := &twiml.Play{Text: fmt.Sprintf("https://s3.amazonaws.com/dialstar.uwaterloo.ca/%d.mp3", ad_counter), Loop: "1"}
	ad_counter++
	ad_counter %= 3

	//Marshal the response
	ad_str, ad_err := xml.Marshal(ad_response)
	if ad_err != nil {
		panic(ad_err)
	}
	//Write the reponse to the Buffer
	b.Write(ad_str)

	redirect := &twiml.Redirect{Text: "http://twilio.axyjo.com/caller/"}
	str, err := xml.Marshal(redirect)
	if err != nil {
		panic(err)
	}
	b.Write(str)

	b.WriteString(end)
	//Write the Buffer to the http.ResponseWriter
	b.WriteTo(w)
}