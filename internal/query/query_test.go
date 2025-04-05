package query

import (
	"reflect"
	"testing"

	"github.com/martishin/search-engine/internal/index"
)

func TestProcessQuery(t *testing.T) {
	idx := index.NewInvertedIndex()
	file1 := "file1.txt"
	idx.AddTokens(file1, []index.TokenizerResult{
		{Tokens: []string{"hello", "world"}, Positions: []int{1, 7}},
	})
	results := ProcessQuery(idx, "Hello")
	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}
	res := results[0]
	if res.Document != file1 {
		t.Errorf("Expected document %s, got %s", file1, res.Document)
	}
	expected := []index.Position{{Pos: 1, Row: 1, TokenPos: 1}}
	if !reflect.DeepEqual(res.Positions, expected) {
		t.Errorf("Expected positions %v, got %v", expected, res.Positions)
	}
}
