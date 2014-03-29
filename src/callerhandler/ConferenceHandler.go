package callerhandler

import (
	"bytes"
	"encoding/xml"
	"fmt"
	_ "io/ioutil"
	"net/http"
	"twiml"
)

func ConferenceHandler(w http.ResponseWriter, r *http.Request) {
	//Error checking
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	//Debugging statements
	Conf_id := r.Form["ConferenceId"]
	fmt.Println("[META] - Twilio request to connect " + Conf_id[0])
	//Create a new Buffer and writes to it. Similar to callerhandler
	b := bytes.NewBufferString(start)
	Say_name := r.Form["OtherCity"]
	say_response := &twiml.Say{Voice: "female", Language: "en", Loop: 1, Text: "Connecting to user from " + Say_name[0]}

	str, err := xml.Marshal(say_response)
	if err != nil {
		panic(err)
	}
	b.Write(str)
	//End the conference on exit
	response := &twiml.Conference{Text: Conf_id[0], EndConferenceOnExit: "true"}
	dial := &twiml.Dial{Conference: *response, HangupOnStar: "true"}
	str, err = xml.Marshal(dial)
	if err != nil {
		panic(err)
	}
	b.Write(str)
	redirect := &twiml.Redirect{Text: "http://twilio.axyjo.com/ad/"}
	str, err = xml.Marshal(redirect)
	if err != nil {
		panic(err)
	}
	b.Write(str)
	b.WriteString(end)
	b.WriteTo(w)
}
