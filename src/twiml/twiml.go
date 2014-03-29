package twiml

type Queue struct {
	Text   string `xml:",chardata"`
	Url    string `xml:"url,attr,omitempty"`
	Method string `xml:"method,attr,omitempty"`
}

type Dial struct {
	Text         string     `xml:",chardata"`
	Action       string     `xml:"action,attr,omitempty"`
	Method       string     `xml:"method,attr,omitempty"`
	Timeout      string     `xml:"timeout,attr,omitempty"`
	HangupOnStar string     `xml:"hangupOnStar,attr,omitempty"`
	TimeLimit    string     `xml:"timeLimit,attr,omitempty"`
	CallerId     string     `xml:"callerId,attr,omitempty"`
	Record       string     `xml:"record,attr,omitempty"`
	Numbers      []Number   `xml:"Number"`
	Clients      []Client   `xml:"Client"`
	Conference   Conference `xml:"Conference,omitempty"`
}

type Thingy struct {
	CallSid string
	City    string
	Add     bool
}

type Conference struct {
	Text                   string `xml:",chardata"`
	Muted                  string `xml:"muted,attr,omitempty"`
	Beep                   string `xml:"beep,attr,omitempty"`
	StartConferenceOnEnter string `xml:"startConferenceOnEnter,attr,omitempty"`
	EndConferenceOnExit    string `xml:"endConferenceOnExit,attr,omitempty"`
	WaitUrl                string `xml:"waitUrl,attr,omitempty"`
	WaitMethod             string `xml:"waitMethod,attr,omitempty"`
	MaxParticipants        string `xml:"maxParticipants,attr,omitempty"`
}

type Client struct {
	Text   string `xml:",chardata"`
	Url    string `xml:"url,attr,omitempty"`
	Method string `xml:"method,attr,omitempty"`
}

type Number struct {
	Text       string `xml:",chardata"`
	SendDigits string `xml:"sendDigits,attr,omitempty"`
	Url        string `xml:"url,attr,omitempty"`
	Method     string `xml:"method,attr,omitempty"`
}

type Play struct {
	Text   string `xml:",chardata"`
	Loop   string `xml:"loop,attr,omitempty"`
	Digits string `xml:"digits,attr,omitempty"`
}
