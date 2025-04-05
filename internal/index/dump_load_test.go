package index

import (
	"os"
	"testing"
)

func TestDumpAndLoad(t *testing.T) {
	idx := NewInvertedIndex()
	filePath := "dummy_index.json"

	idx.AddTokens("file1.txt", []TokenizerResult{
		{Tokens: []string{"Hello", "world"}, Positions: []int{1, 10}},
	})

	if err := DumpIndex(filePath, idx); err != nil {
		t.Fatalf("DumpIndex failed: %v", err)
	}
	defer os.Remove(filePath)

	loadedIdx, err := LoadIndex(filePath)
	if err != nil {
		t.Fatalf("LoadIndex failed: %v", err)
	}

	positions := loadedIdx.GetPositions([]string{"Hello"})
	if len(positions) != 1 {
		t.Errorf("Expected 1 document for token 'Hello', got %d", len(positions))
	}
}
