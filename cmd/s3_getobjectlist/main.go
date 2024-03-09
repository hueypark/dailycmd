package main

import (
	"context"
	"flag"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func main() {
	bucket := flag.String("bucket", "", "Bucket name")
	prefix := flag.String("prefix", "", "Prefix")
	filteredStorageClass := (*types.ObjectStorageClass)(flag.String("filteredStorageClass", "", "Filtered storage class"))
	flag.Parse()

	slog.Info("s3 read", "bucket", *bucket, "prefix", *prefix)

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		slog.Error("unable to load SDK config", "error", err)
		return
	}

	s3Cli := s3.NewFromConfig(cfg)

	req := &s3.ListObjectsV2Input{
		Bucket: bucket,
		Prefix: prefix,
	}
	counter := 0
	for {
		res, err := s3Cli.ListObjectsV2(context.Background(), req)
		if err != nil {
			slog.Error("unable to read object", "error", err)
			return
		}

		for _, obj := range res.Contents {
			counter++

			if counter%100_000 == 0 {
				slog.Info("Progressing...", "counter", counter)
			}

			if obj.StorageClass == *filteredStorageClass {
				continue
			}

			slog.Info("Object", "key", string(*obj.Key), "storageClass", string(obj.StorageClass))
		}

		if !*res.IsTruncated {
			break
		}

		req.ContinuationToken = res.NextContinuationToken
	}

	slog.Info("Done", "counter", counter)
}
