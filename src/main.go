package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("Enter q for quit")
	go rotate(time.Now().Hour())
	p := NewPublisher()
	ch := Checker{publisher: p}
	go ch.check()

	var s string
	for s != "q" {
		fmt.Scanln(&s)
	}
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
