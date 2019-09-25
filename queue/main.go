package queue

import (
	"github.com/go-redis/redis"
)

func Init(redisClient *redis.Client) *Config {
	return &Config{
		SMS:      channel{name: "sms", client: redisClient},
		Whatsapp: channel{name: "whatsapp", client: redisClient},
	}
}

type channel struct {
	name   string
	client *redis.Client
}

type Config struct {
	SMS      channel
	Whatsapp channel
}

func (c channel) Publish(message string) error {
	return c.client.Publish(c.name, message).Err()
}

func (c channel) Subscribe() *redis.PubSub {
	return c.client.Subscribe(c.name)
}

func (c channel) Close() error {
	return c.client.Subscribe(c.name).Close()
}
