package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// Redis encapsule le client Redis
type Redis struct {
	client *redis.Client
}

func NewRedis(url string) *Redis {
	opt, err := redis.ParseURL(url)
	if err != nil {
		log.Fatal().Err(err).Msg("redis parse URL")
	}
	client := redis.NewClient(opt)

	// Ping
	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Warn().Err(err).Msg("Redis ping failed — continuera en dégradé")
	} else {
		log.Info().Msg("Redis connecté")
	}
	return &Redis{client: client}
}

func (r *Redis) Client() *redis.Client { return r.client }

func (r *Redis) Close() { r.client.Close() }
