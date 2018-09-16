package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("Enter q for quit")
	go rotate(time.Now().Hour())
	go check()

	var s string
	for s != "q" {
		fmt.Scanln(&s)
	}
}

func check() {
	config := readConfiguration()

	for _, u := range config.URL {
		go func(url string) {
			now := time.Now()
			r, err := http.Get(url)
			d := time.Since(now)
			if err != nil {
				log.WithFields(log.Fields{"Error": err, "Elapsed": d.String(), "url": url}).Warn("Error")
			} else {
				log.WithFields(log.Fields{"StatusCode": r.StatusCode, "Elapsed": d.String(), "url": url}).Info("ok")
			}
		}(u)
	}

	time.Sleep(config.Interval.Duration)
	check()
}

func rotate(hour int) {
	time.Sleep(time.Minute)

	now := time.Now()
	if now.Hour() < hour {
		// change only if when switched from 23->0
		name, _ := filepath.Abs("./logs/log-" + now.Format("20060102") + ".log")
		file, _ := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		mw := io.MultiWriter(os.Stdout, file)
		log.SetOutput(mw)
	}

	rotate(now.Hour())
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	os.Mkdir("logs", os.ModeDir)
	name, _ := filepath.Abs("./logs/log-" + time.Now().Format("20060102") + ".log")
	file, _ := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)
}
