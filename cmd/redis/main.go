package main

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// airbridge-api-private-ro.rtid8t.ng.0001.apne1.cache.amazonaws.com:6379

	r := redis.NewClient(&redis.Options{
		Addr: "airbridge-checkpoint-private.rtid8t.clustercfg.apne1.cache.amazonaws.com:6379",
	})

	res, err := r.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("redis ping failed")
	}

	log.Info().
		Str("response", res).
		Msg("ping success")
}
