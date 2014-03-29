package callerhandler

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/gorilla/schema"
	_ "io/ioutil"
	"net/http"
	"twiml"
	"utils"
	"webui"
)

type WelcomeWrapper struct {
	Push *[]chan webui.PushData
}

//start and end of the xml sent to Twilio
//CallerWrapper which holds a channel the interface Thingy

func (c WelcomeWrapper) WelcomeHandler(w http.ResponseWriter, r *http.Request) {
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

	userCount := utils.GetUserCount()

	fmt.Printf("%.6s has called in\n", request.CallSid)
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

	pData := webui.PushData{UserCount: userCount + 1}

	if webui.UseNumbers {
		pData.Call1Id = request.From
	} else {
		pData.Call1Id = request.CallSid
	}

	for _, j := range *c.Push {
		j <- pData
	}

}
