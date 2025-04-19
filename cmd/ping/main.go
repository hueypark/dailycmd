package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
		},
	)

	url := flag.String("url", "", "URL")
	interval := flag.Duration("interval", time.Second, "Interval")
	flag.Parse()

	log.Info().
		Str("url", *url).
		Dur("interval", *interval).
		Msg("ping started")

	for {
		go func() {
			start := time.Now()
			resp, err := http.Get(*url)
			if err != nil {
				log.Info().
					Err(err).
					Msg("ping failed")
			}
			resp.Body.Close()

			// body, err := io.ReadAll(resp.Body)
			// if err != nil {
			// 	log.Info().
			// 		Err(err).
			// 		Msg("read body failed")
			// }

			if resp.StatusCode == http.StatusOK {
				log.Info().
					Time("start", start).
					Int("status", resp.StatusCode).
					Msg("ping success")
			} else {
				log.Warn().
					Time("start", start).
					Int("status", resp.StatusCode).
					Msg("ping failed")
			}
		}()

		time.Sleep(*interval)
	}
}
