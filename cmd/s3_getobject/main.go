package main

import (
	"context"
	"flag"
	"io"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	bucket := flag.String("bucket", "", "Bucket name")
	key := flag.String("key", "", "Object key")
	flag.Parse()

	slog.Info("s3 read", "bucket", *bucket, "key", *key)

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		slog.Error("unable to load SDK config", "error", err)
		return
	}

	s3Cli := s3.NewFromConfig(cfg)

	obj, err := s3Cli.GetObject(
		context.Background(),
		&s3.GetObjectInput{
			Bucket: bucket,
			Key:    key,
		},
	)
	if err != nil {
		slog.Error("unable to read object", "error", err)
		return
	}

	body, err := io.ReadAll(obj.Body)
	if err != nil {
		slog.Error("unable to read object body", "error", err)
		return
	}

	slog.Info("object body", "body", string(body))
}
