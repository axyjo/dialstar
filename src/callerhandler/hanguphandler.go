package callerhandler

import (
	"fmt"
	"github.com/gorilla/schema"
	_ "io/ioutil"
	"net/http"
	"twiml"
	"utils"
	"webui"
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
	Push     *[]chan webui.PushData
}

func (c HangUpWrapper) HangUpHandler(w http.ResponseWriter, r *http.Request) {
	//Check if it's a POST call, if not return immediately
	if r.Method != "POST" {
		return
	}
	//Error checking
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	// Unmarshal Form data
	var request StatusCallbackRequest
	decoder := schema.NewDecoder()
	decoder.Decode(&request, r.Form)
	//If the user has hung up, send n new request through the channel to dequeue the user
	if request.CallStatus == "completed" {
		c.Callerid <- twiml.Thingy{request.CallSid, "", false}
		fmt.Printf("%.6s has hung up\n", request.CallSid)
	}

	userCount := utils.GetUserCount()
	for _, j := range *c.Push {
		j <- webui.PushData{
			UserCount: userCount,
			Call1Id:   request.CallSid,
		}
	}
}
