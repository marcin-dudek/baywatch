package main

import (
	"github.com/go-redis/redis"
)

// Publisher publishes health checking
type Publisher struct {
	redis *redis.Client
}

// NewPublisher creates new publisher
func NewPublisher() Publisher {
	client := redis.NewClient(&redis.Options{Addr: "redis:6379", Password: "", DB: 0})
	return Publisher{redis: client}
}

// Publish sends message
func (p Publisher) Publish(channel string, message interface{}) {
	p.redis.Publish(channel, message)
}
