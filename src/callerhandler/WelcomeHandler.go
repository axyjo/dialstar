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

//CallerWrapper which holds a channel the interface Thingy

//Holds information about the user calling in

type Marshal struct {
	total interface{}
}

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
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

	//Creates a new Buffer with the initial start xml string
	b := bytes.NewBufferString(start)

	//Marshal the say_repsonse
	stats_url := `https://AC6f0fa1837933462d780f6fc1daf57d44:79ed2712d0cf06c87aa2783eee6aaa7a@api.twilio.com/2010-04-01/Accounts/AC6f0fa1837933462d780f6fc1daf57d44/Calls.json?Status=in-progress`
	resp, err := http.Get(stats_url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	active_users := &Marshal{}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&active_users.total)
	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprint(resp.Body))
	fmt.Println("[META] - Current Active Users: " + fmt.Sprint(active_users.total))
	text := "Welcome to Dial Star! There are  " + fmt.Sprint(active_users.total)
	text = text + " other users. Press star to skip a user."
	say_response := &twiml.Say{Voice: "female", Language: "en", Loop: 1, Text: text}
	str, err := xml.Marshal(say_response)
	//Error checking..
	if err != nil {
		panic(err)
	}
	//Append the Say block to the buffer
	b.Write(str)
	fmt.Println(request.CallSid + " - was welcomed")

	redirect := &twiml.Redirect{Text: "http://twilio.axyjo.com/caller/"}
	str, err = xml.Marshal(redirect)
	if err != nil {
		panic(err)
	}
	b.Write(str)
	//Write the end of the xml to the Buffer
	b.WriteString(end)
	//Write the Buffer to the http.ResponseWriter
	b.WriteTo(w)
}
