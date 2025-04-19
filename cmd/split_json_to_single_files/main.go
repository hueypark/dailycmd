package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/pretty"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	srcDir := flag.String("src-dir", "", "Source directory")
	dstDir := flag.String("dst-dir", "", "Destination directory")
	dstPrefix := flag.String("dst-prefix", "", "Destination file prefix")
	flag.Parse()

	log.Info().
		Str("src-dir", *srcDir).
		Str("dst-dir", *dstDir).
		Msg("split json to single files")

	err := splitJsonToSingleFilesWithDir(*srcDir, *dstDir, *dstPrefix)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to split json to single files")
	}

	log.Info().
		Msg("done")
}

func splitJsonToSingleFilesWithDir(srcDir, dstDir string, dstPrefix string) error {
	_, err := os.Stat(dstDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dstDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dstDir, err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check directory %s: %v", dstDir, err)
	}

	idx := 0

	return filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		idx, err = splitJsonToSingleFiles(path, dstDir, dstPrefix, idx)
		if err != nil {
			return fmt.Errorf("failed to extract %s: %v", path, err)
		}

		return nil

	})
}

func splitJsonToSingleFiles(path string, dstDir string, dstPrefix string, idx int) (int, error) {
	sf, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("failed to open file %s: %v", path, err)
	}
	defer sf.Close()

	scanner := bufio.NewScanner(sf)
	scanner.Buffer(nil, 1025*1024)
	for scanner.Scan() {
		dstPath := filepath.Join(dstDir, fmt.Sprintf("%s%d.json", dstPrefix, idx))

		text := scanner.Text()

		err := createFile(dstPath, []byte(text))
		if err != nil {
			return 0, fmt.Errorf("failed to create file %s: %v", dstPath, err)
		}

		idx++
	}

	err = scanner.Err()
	if err != nil {
		return 0, fmt.Errorf("failed to scan file content: %v", err)
	}

	return idx, nil
}

func createFile(path string, content []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", path, err)
	}
	defer f.Close()

	_, err = f.Write(pretty.Pretty(content))
	if err != nil {
		return fmt.Errorf("failed to write content to file %s: %v", path, err)
	}

	return nil
}
