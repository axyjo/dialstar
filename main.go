package main

import (
	"callerhandler"
	"container/list"
	"fmt"
	"net/http"
	"net/url"
	"twiml"
)

var user_queue *list.List

//This is just all the handlers and shit
func main() {
	callers_waiting := make(chan twiml.Thingy, 10)
	Conf_waiters := callerhandler.CallerWrapper{Callerid: callers_waiting}
	Conf_dequeue := callerhandler.HangUpWrapper{Callerid: callers_waiting}
	go PollWaiters(callers_waiting)
	http.HandleFunc("/caller/", Conf_waiters.CallerHandler)
	http.HandleFunc("/conference/", callerhandler.ConferenceHandler)
	http.HandleFunc("/hangup/", Conf_dequeue.HangUpHandler)
	http.ListenAndServe(":3000", nil)
}

func PollWaiters(c chan twiml.Thingy) {
	user_queue = list.New()

	for element := range c {
		if element.Add {
			_ = user_queue.PushBack(element)
			if user_queue.Len() >= 2 {
				fmt.Println("[META] - Got two or more people.")
				first := user_queue.Front()
				user_queue.Remove(first)
				f := first.Value.(twiml.Thingy).CallSid
				second := user_queue.Front()
				user_queue.Remove(second)
				s := second.Value.(twiml.Thingy).CallSid
				fmt.Println("[META] Paired " + f + " with " + s)

				ConferenceId := f + s
				ConfURLBase := "http://twilio.axyjo.com/conference/?ConferenceId=" + ConferenceId + "&OtherCity="
				authToken := "79ed2712d0cf06c87aa2783eee6aaa7a"
				accountId := "AC6f0fa1837933462d780f6fc1daf57d44"

				values := make(url.Values)
				values.Set("Url", ConfURLBase+url.QueryEscape(second.Value.(twiml.Thingy).City))
				fmt.Println(values)
				_, err := http.PostForm("https://"+accountId+":"+authToken+"@api.twilio.com/2010-04-01/Accounts/"+accountId+"/Calls/"+f, values)
				if err != nil {
					panic(err)
				}
				values.Set("Url", ConfURLBase+url.QueryEscape(first.Value.(twiml.Thingy).City))
				fmt.Println(values)
				_, err = http.PostForm("https://"+accountId+":"+authToken+"@api.twilio.com/2010-04-01/Accounts/"+accountId+"/Calls/"+s, values)
				if err != nil {
					panic(err)
				}
			}
		} else {
			for i := user_queue.Front(); i != nil; i = i.Next() {
				if i.Value.(twiml.Thingy).CallSid == element.CallSid {
					user_queue.Remove(i)
					break
				}
			}
		}
	}
}
