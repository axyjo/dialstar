package callerhandler

import (
	"fmt"
	"github.com/gorilla/schema"
	_ "io/ioutil"
	"net/http"
	"twiml"
)

type StatusCallbackRequest struct {
	CallDuration      string
	RecordingUrl      string
	RecordingSid      string
	RecordingDuration string
	CallStatus        string
	AccountSid        string
	CallSid           string
}

type HangUpWrapper struct {
	Callerid chan twiml.Thingy
}

func (c HangUpWrapper) HangUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	// Unmarshal Form data
	var request StatusCallbackRequest
	decoder := schema.NewDecoder()
	decoder.Decode(&request, r.Form)
	if request.CallStatus == "completed" {
		c.Callerid <- twiml.Thingy{request.CallSid, false}
		fmt.Println(request.CallSid + " - deQueued")
	}
}
