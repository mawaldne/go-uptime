package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/go-uptime/goslack"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	Urls        []string
	WaitSeconds int64
	Slack       SlackInfo
}

type SlackInfo struct {
	Url      string
	UserName string
}

var (
	Info   *log.Logger
	Error  *log.Logger
	config Config
	notify *goslack.Notifier
)

func main() {
	Init()
	for {
		for _, url := range config.Urls {
			Info.Printf("Url: %s\n", url)
			ping(url)
		}
		time.Sleep(time.Duration(config.WaitSeconds) * time.Second)
	}
}

func Init() {
	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	if _, err := toml.DecodeFile("uptime.toml", &config); err != nil {
		panic(err)
	}
	notify = goslack.New(config.Slack.Url)
	notify.Username = config.Slack.UserName
}

func ping(url string) {
	res, err := http.Head(url)
	if err == nil {
		if res.StatusCode == 200 {
			msg := fmt.Sprintf("%s is up...Moving on", url)
			Info.Printf(msg)
			notify_slack(url, msg)
		} else {
			msg := fmt.Sprintf("%s is down with status: %s", url, res.Status)
			Error.Printf(msg)
			notify_slack(url, msg)
		}
	} else {
		notify_slack(url, err.Error())
	}
}

func notify_slack(url string, msg string) {
	if _, err := notify.Send(msg); err != nil {
		Error.Printf(err.Error())
	}
}
