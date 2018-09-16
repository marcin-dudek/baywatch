package main

import (
	"bytes"
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

	for _, u := range config.Request {
		go func(r Request) {
			req, err := http.NewRequest(r.Method, r.URL, bytes.NewBuffer([]byte(r.Body)))
			if err != nil {
				log.WithFields(log.Fields{"Error": err, "url": r.URL}).Warn("Configuration error")
				return
			}

			for k, v := range r.Header {
				req.Header.Set(k, v)
			}

			now := time.Now()
			client := &http.Client{}
			resp, err := client.Do(req)
			d := time.Since(now)
			if err != nil {
				log.WithFields(log.Fields{"Error": err, "Elapsed": d.String(), "url": r.URL}).Warn("Error")
			} else {
				log.WithFields(log.Fields{"StatusCode": resp.StatusCode, "Elapsed": d.String(), "url": r.URL}).Info("ok")
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
