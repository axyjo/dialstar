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
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	Conf_id := r.Form["ConferenceId"]
	fmt.Println("[META] - Twilio request to connect " + Conf_id[0])
	b := bytes.NewBufferString(start)
	Say_name := r.Form["OtherCity"]
	say_response := &Say{Voice: "female", Language: "en", Loop: 1, Text: "Connecting to user from " + Say_name[0]}

	str, err := xml.Marshal(say_response)
	if err != nil {
		panic(err)
	}
	b.Write(str)
	response := &twiml.Conference{Text: Conf_id[0], EndConferenceOnExit: "true"}
	dial := &twiml.Dial{Conference: *response}
	str, err = xml.Marshal(dial)
	if err != nil {
		panic(err)
	}
	b.Write(str)
	b.WriteString(end)
	b.WriteTo(w)
}
