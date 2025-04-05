package file

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/martishin/search-engine/internal/index"
)

func TestLoadFiles(t *testing.T) {
	dir := filepath.Join("testdata")
	idx := index.NewInvertedIndex()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := LoadFiles(ctx, dir, idx)
	if err != nil {
		t.Fatalf("LoadFiles error: %v", err)
	}

	positions := idx.GetPositions([]string{"hello"})
	if len(positions) == 0 {
		t.Errorf("Expected some results for token 'hello', got none")
	}
}

func TestLoadFilesCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	idx := index.NewInvertedIndex()
	err := LoadFiles(ctx, "testdata", idx)
	if err == nil {
		t.Error("Expected error from canceled context, got nil")
	}
}
