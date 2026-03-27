package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	*redis.Client
}

func New(redisURL string) (*Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &Client{client}, nil
}

func (c *Client) Close() error {
	return c.Client.Close()
}
