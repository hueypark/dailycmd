package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	bucket := flag.String("bucket", "", "Bucket name")
	prefix := flag.String("prefix", "", "Prefix")
	downloadDir := flag.String("download-dir", "", "Download directory")
	flag.Parse()

	log.Info().
		Str("bucket", *bucket).
		Str("prefix", *prefix).
		Msg("s3 download directory")

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithSharedConfigProfile("QE-LEAD-FB"))
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("unable to load SDK config")
	}

	cli := s3.NewFromConfig(cfg)

	req := &s3.ListObjectsV2Input{
		Bucket: bucket,
		Prefix: prefix,
	}
	objectCnt := 0
	for {
		res, err := cli.ListObjectsV2(context.Background(), req)
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("unable to list objects")
		}

		for _, obj := range res.Contents {
			err = download(cli, *bucket, *obj.Key, *downloadDir)
			if err != nil {
				log.Fatal().
					Err(err).
					Msg("unable to download object")

				continue
			}

			log.Info().
				Str("objectKey", *obj.Key).
				Msg("downloaded")
		}

		objectCnt += len(res.Contents)

		if !*res.IsTruncated {
			break
		}

		req.ContinuationToken = res.NextContinuationToken
	}

	log.Info().
		Int("objectCount", objectCnt).
		Msg("done")
}

func download(cli *s3.Client, bucket, key, downloadDir string) error {
	req := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}

	res, err := cli.GetObject(context.Background(), req)
	if err != nil {
		return fmt.Errorf("unable to get object, %w", err)
	}
	defer res.Body.Close()

	p := path.Join(downloadDir, path.Base(key))
	f, err := os.Create(p)
	if err != nil {
		return fmt.Errorf("unable to create file, %w", err)
	}
	defer f.Close()

	_, err = f.ReadFrom(res.Body)
	if err != nil {
		return fmt.Errorf("unable to write to file, %w", err)
	}

	return nil
}
