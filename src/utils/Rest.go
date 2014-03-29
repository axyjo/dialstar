package utils

import (
	"encoding/json"
	"net/http"
)

type ActiveUsers struct {
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
