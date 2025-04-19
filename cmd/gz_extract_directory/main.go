package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	srcDir := flag.String("src-dir", "", "Source directory")
	dstDir := flag.String("dst-dir", "", "Destination directory")
	flag.Parse()

	log.Info().
		Str("src-dir", *srcDir).
		Str("dst-dir", *dstDir).
		Msg("gz extract directory")

	err := extractGzInDir(*srcDir, *dstDir)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to extract gz in directory")
	}

	log.Info().
		Msg("done")
}

func extractGzInDir(srcDir, dstDir string) error {
	_, err := os.Stat(dstDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dstDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dstDir, err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check directory %s: %v", dstDir, err)
	}

	return filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(d.Name()) != ".gz" {
			return nil
		}

		dest := strings.Replace(filepath.Join(dstDir, d.Name()), ".gz", "", 1)

		err = extractGz(path, dest)
		if err != nil {
			return fmt.Errorf("failed to extract %s: %v", path, err)
		}

		return nil

	})
}

func extractGz(src, dest string) error {
	log.Info().
		Str("src", src).
		Str("dest", dest).
		Msg("extract gz")

	file, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", src, err)
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %v", err)
	}
	defer gzReader.Close()

	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(f, gzReader)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %v", err)
	}

	return nil
}
