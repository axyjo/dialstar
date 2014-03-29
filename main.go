package main

import (
	"callerhandler"
	"container/list"
	"fmt"
	"net/http"
	"net/url"
	"twiml"
	"utils"
	"webui"
)

var user_queue *list.List

//This is just all the handlers and shit
func main() {
	//Create a new channel of size 10 (shouldn't get much larger than this)
	callers_waiting := make(chan twiml.Thingy, 10)
	push := make([]chan webui.PushData, 0)

	//Create a new CallerHandler with a CallerWrapper/HangupWrapper with the shared channel callers_waiting
	Conf_waiters := callerhandler.CallerWrapper{Callerid: callers_waiting}
	Conf_dequeue := callerhandler.HangUpWrapper{Callerid: callers_waiting, Push: &push}
	Conf_newUser := callerhandler.WelcomeWrapper{Push: &push}
	Conf_push := webui.WebSocketWrapper{Push: &push}
	Conf_adremove := callerhandler.AdWrapper{Push: &push}

	//Have a function that polls users and queues and dequeues users as necessary
	go PollWaiters(callers_waiting, &push)

	//Register the Handle functions for the given patters and appropriate handlers
	http.HandleFunc("/caller/", Conf_waiters.CallerHandler)
	http.HandleFunc("/conference/", callerhandler.ConferenceHandler)
	http.HandleFunc("/hangup/", Conf_dequeue.HangUpHandler)
	http.HandleFunc("/welcome/", Conf_newUser.WelcomeHandler)
	http.HandleFunc("/ad/", Conf_adremove.AdHandler)
	http.HandleFunc("/webui/websocket", Conf_push.WebSocketHandler)
	http.HandleFunc("/webui/", webui.WebHandler)
	//Starts the HTTP server at the address Localhost:3000

	port := 3000
	fmt.Printf("HTTP server listening on port %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func PollWaiters(c chan twiml.Thingy, p *[]chan webui.PushData) {
	//Creates an empty queue to put users in
	user_queue = list.New()
	//Iterates over each element in the channel
	for element := range c {
		//If the user is to be added to the queue..
		if element.Add {
			//Push the user onto the queue
			_ = user_queue.PushBack(element)
			fmt.Printf("%.6s was added into the queue\n", element.CallSid)
			len := user_queue.Len()
			fmt.Printf("There are %d users in the queue\n", len)
			//If there are 2 or more users in the queue
			if len >= 3 {
				//Get a pointer to the first element of the queue
				first := user_queue.Front()
				f := first.Value.(twiml.Thingy).CallSid
				//Get a pointer to the second user who's in the third position in the queue to prevent call repetition
				second := user_queue.Front().Next().Next()
				s := second.Value.(twiml.Thingy).CallSid
				//remove the first and second user from the queue
				user_queue.Remove(first)
				user_queue.Remove(second)
				fmt.Printf("Created a conference for %.6s and %.6s\n", f, s)
				fmt.Printf("There are %d users in the queue\n", user_queue.Len())
				//Concatenate the first and second user's CallSid to be used in the ConfURL
				ConferenceId := f + s
				ConfURLBase := "http://twilio.axyjo.com/conference/?ConferenceId=" + ConferenceId + "&OtherCity="
				//Authentication stuff
				authToken := "79ed2712d0cf06c87aa2783eee6aaa7a"
				accountId := "AC6f0fa1837933462d780f6fc1daf57d44"
				//Map of strings to arry of strings
				values := make(url.Values)

				//Set the URLs for both the first and second user and put them into a conference together
				values.Set("Url", ConfURLBase+url.QueryEscape(second.Value.(twiml.Thingy).City))
				_, err := http.PostForm("https://"+accountId+":"+authToken+"@api.twilio.com/2010-04-01/Accounts/"+accountId+"/Calls/"+f, values)
				if err != nil {
					panic(err)
				}
				values.Set("Url", ConfURLBase+url.QueryEscape(first.Value.(twiml.Thingy).City))
				_, err = http.PostForm("https://"+accountId+":"+authToken+"@api.twilio.com/2010-04-01/Accounts/"+accountId+"/Calls/"+s, values)
				if err != nil {
					panic(err)
				}
				for _, j := range *p {
					j <- webui.PushData{
						UserCount: utils.GetUserCount(),
						Call1Id:   f,
						Call2Id:   s,
					}
				}
			}
		} else {
			//Otherwise the user is to be dequeued.
			//Iterate through the list and find a matching CallSid, remove the user from the queue and break out of the loop
			for i := user_queue.Front(); i != nil; i = i.Next() {
				if i.Value.(twiml.Thingy).CallSid == element.CallSid {
					user_queue.Remove(i)

					for _, j := range *p {
						j <- webui.PushData{UserCount: utils.GetUserCount()}
					}
					break
				}
			}
		}
	}
}
