package main

import (
	"github.com/BurntSushi/toml"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	Urls        []string
	WaitSeconds int64
}

var (
	Info   *log.Logger
	Error  *log.Logger
	config Config
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

	if _, err := toml.DecodeFile("urls.toml", &config); err != nil {
		panic(err)
	}
}

func ping(url string) {
	res, err := http.Head(url)
	if err == nil {
		if res.StatusCode == 200 {
			Info.Printf("%s is up. Moving on...", url)
		} else {
			Error.Printf("%s is down with status: %s", url, res.Status)
		}
	} else {
		Error.Printf("%s is down with status: %s", url, err)
	}
}
