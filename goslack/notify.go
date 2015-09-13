package goslack

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Payload struct {
	Text      string `json:"text"`
	Username  string `json:"username"`
	IconEmoji string `json:"icon_emoji"`
	IconUrl   string `json:"icon_url"`
	Channel   string `json:"channel"`
}

type Notifier struct {
	webhook_url string
	Payload
}

func New(webhook_url string) *Notifier {
	return &Notifier{webhook_url, Payload{}}
}

func (n Notifier) Send(contents string) (string, error) {
	n.Text = contents
	params, _ := json.Marshal(n.Payload)

	res, err := http.PostForm(
		n.webhook_url,
		url.Values{"payload": {string(params)}},
	)
	if err != nil {
		return "", errors.New("Unable to communicate with slack.")
	}

	body, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	return string(body), nil
}
