package file

import (
	"bufio"
	"os"
	"path/filepath"

	"github.com/martishin/search-engine/internal/index"
	"github.com/martishin/search-engine/internal/text"
)

// ListFiles recursively lists all files (not directories) under baseDirectory.
func ListFiles(baseDirectory string) ([]string, error) {
	var files []string
	err := filepath.Walk(baseDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// ReadFile reads a file line by line, tokenizes each line, and returns the tokens.
func ReadFile(filePath string) ([]index.TokenizerResult, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var results []index.TokenizerResult
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		results = append(results, text.ProcessText(line))
	}
	if err := scanner.Err(); err != nil {
		return results, err
	}
	return results, nil
}
