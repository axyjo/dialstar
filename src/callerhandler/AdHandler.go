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
	"webui"
	"github.com/mattbaird/elastigo/api"
	"github.com/mattbaird/elastigo/core"
)

type AdWrapper struct {
	Push      *[]chan webui.PushData
	AdsPlayed []int
}

var ad_counter = 0

func (c AdWrapper) AdHandler(w http.ResponseWriter, r *http.Request) {
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

	// Set up ES host/port.
	api.Domain = "twilio.axyjo.com"
	api.Port = "9200"


	// Store the information from the request into ElasticSearch for analytics
	bytesLine, err := json.Marshal(request)
	es_response, err2 := core.Index("hackathon", "logs", "", nil, bytesLine)
	fmt.Println(es_response)
	if (err2 != nil) {
		panic(err2)
	}

	if request.CallStatus == "completed" {
		return
	}

	//Creates a new Buffer with the initial start xml string
	b := bytes.NewBufferString(start)

	// Queue an ad.

	ad_response := &twiml.Play{Text: fmt.Sprintf("https://s3.amazonaws.com/dialstar.uwaterloo.ca/%d.mp3", ad_counter), Loop: "1"}
	c.AdsPlayed[ad_counter]++
	fmt.Println(c.AdsPlayed)
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

	pData := webui.PushData{UserCount: -1}
	if webui.UseNumbers {
		pData.Call1Id = request.From
	} else {
		pData.Call1Id = request.CallSid
		pData.Ads = c.AdsPlayed
	}

	for _, j := range *c.Push {
		j <- pData
	}
}
