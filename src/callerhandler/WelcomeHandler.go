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

type ActiveUsers struct {
	Total int `json:"total"`
}

type ActiveConferences struct {
	Total int `json:"total"`
}

func GetUserCount() int {
	//Marshal the say_repsonse
	stats_url := `https://AC6f0fa1837933462d780f6fc1daf57d44:79ed2712d0cf06c87aa2783eee6aaa7a@api.twilio.com/2010-04-01/Accounts/AC6f0fa1837933462d780f6fc1daf57d44/Calls.json?Status=in-progress`
	resp, err := http.Get(stats_url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	active_users := &ActiveUsers{}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&active_users)
	if err != nil {
		panic(err)
	}
	return active_users.Total
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

	userCount := GetUserCount()

	fmt.Printf("%.5s has called in\n", request.CallSid)
	fmt.Printf("There are %d users\n", userCount+1)
	text := fmt.Sprintf("Welcome to Dial Star! There are %d other users. Press star to skip a user.", userCount)
	say_response := &twiml.Say{Voice: "female", Language: "en", Loop: 1, Text: text}
	str, err := xml.Marshal(say_response)
	//Error checking..
	if err != nil {
		panic(err)
	}
	//Append the Say block to the buffer
	b.Write(str)

	redirect := &twiml.Redirect{Text: "http://twilio.axyjo.com/ad/"}
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
