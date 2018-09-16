package main

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

// Config contains configuration information
type Config struct {
	URL      []string `json:"url"`
	Interval Duration `json:"interval"`
}

func readConfiguration() Config {
	configFile, err := os.Open("config.json")
	if err != nil {
		log.WithFields(log.Fields{"Error": err.Error()}).Error("opening config file")
	}

	var config Config

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&config); err != nil {
		log.WithFields(log.Fields{"Error": err.Error()}).Error("parsing config file")
	}

	return config
}

// Duration copy
type Duration struct {
	time.Duration
}

// MarshalJSON convert duration to string
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// UnmarshalJSON converts back
func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}
