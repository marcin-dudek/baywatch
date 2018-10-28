package main

import (
	"bytes"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// Checker is service for health monitoring
type Checker struct {
	publisher Publisher
}

func (ch Checker) check() {
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
			x := log.Fields{"Elapsed": d.String(), "url": r.URL}
			if err != nil {
				x["Error"] = err
				log.WithFields(x).Warn("Error")
				ch.publisher.Publish("healthcheck", x)
			} else {
				x["StatusCode"] = resp.StatusCode
				log.WithFields(x).Info("ok")
				ch.publisher.Publish("healthcheck", x)
			}
		}(u)
	}

	time.Sleep(config.Interval.Duration)
	ch.check()
}
