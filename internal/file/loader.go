package file

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/martishin/search-engine/internal/index"
)

// LoadFiles walks the base directory, reads files concurrently, tokenizes them,
// and adds tokens to the index. It supports cancellation via the provided context.
func LoadFiles(ctx context.Context, baseDir string, idx *index.InvertedIndex) error {
	files, err := ListFiles(baseDir)
	if err != nil {
		return err
	}

	remainingFiles := int32(len(files))
	var wg sync.WaitGroup

	for _, filePath := range files {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()
			fmt.Printf("Started parsing file %s\n", filePath)
			tokens, err := ReadFile(filePath)
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", filePath, err)
				return
			}
			idx.AddTokens(filePath, tokens)
			atomic.AddInt32(&remainingFiles, -1)
			fmt.Printf("Finished parsing file %s. Remaining: %d\n", filePath, atomic.LoadInt32(&remainingFiles))
		}(filePath)
	}

	wg.Wait()
	return nil
}
