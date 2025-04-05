package index

import (
	"reflect"
	"testing"
)

func TestAddTokensAndGetPositions(t *testing.T) {
	idx := NewInvertedIndex()
	tokens := []TokenizerResult{
		{Tokens: []string{"hello", "world"}, Positions: []int{1, 7}},
	}
	filePath := "test.txt"
	idx.AddTokens(filePath, tokens)

	positions := idx.GetPositions([]string{"hello"})
	if _, exists := positions[filePath]; !exists {
		t.Errorf("Expected file %s in results", filePath)
	}
	posList := positions[filePath]
	if len(posList) != 1 {
		t.Errorf("Expected one position for 'hello', got %d", len(posList))
	}
	expected := Position{
		Pos:      1,
		Row:      1,
		TokenPos: 1,
	}
	if !reflect.DeepEqual(posList[0], expected) {
		t.Errorf("Expected %v, got %v", expected, posList[0])
	}
}

func TestMultipleFiles(t *testing.T) {
	idx := NewInvertedIndex()
	file1 := "file1.txt"
	file2 := "file2.txt"
	idx.AddTokens(file1, []TokenizerResult{
		{Tokens: []string{"hello", "world"}, Positions: []int{1, 10}},
	})
	idx.AddTokens(file2, []TokenizerResult{
		{Tokens: []string{"world"}, Positions: []int{1}},
	})
	positions := idx.GetPositions([]string{"world"})
	expected := map[string][]Position{
		file1: {{Pos: 2, Row: 1, TokenPos: 10}},
		file2: {{Pos: 1, Row: 1, TokenPos: 1}},
	}
	if !reflect.DeepEqual(positions, expected) {
		t.Errorf("Expected positions %v, got %v", expected, positions)
	}
}

func TestMultipleOccurrences(t *testing.T) {
	idx := NewInvertedIndex()
	file1 := "file1.txt"
	idx.AddTokens(file1, []TokenizerResult{
		{Tokens: []string{"hello", "hello"}, Positions: []int{1, 10}},
	})
	positions := idx.GetPositions([]string{"hello"})
	expected := map[string][]Position{
		file1: {
			{Pos: 1, Row: 1, TokenPos: 1},
			{Pos: 2, Row: 1, TokenPos: 10},
		},
	}
	if !reflect.DeepEqual(positions, expected) {
		t.Errorf("Expected positions %v, got %v", expected, positions)
	}
}

func TestEmpty(t *testing.T) {
	idx := NewInvertedIndex()
	file1 := "file1.txt"
	idx.AddTokens(file1, []TokenizerResult{
		{Tokens: []string{"hello", "hello"}, Positions: []int{1, 10}},
	})
	positions := idx.GetPositions([]string{"app"})
	if len(positions) != 0 {
		t.Errorf("Expected no positions for token 'app', got %v", positions)
	}
}

func TestMultipleRows(t *testing.T) {
	idx := NewInvertedIndex()
	file1 := "file1.txt"
	idx.AddTokens(file1, []TokenizerResult{
		{Tokens: []string{"hello", "world"}, Positions: []int{1, 10}},
		{Tokens: []string{"world", "hello"}, Positions: []int{1, 10}},
	})

	expected := map[string][]Position{
		file1: {
			{Pos: 1, Row: 1, TokenPos: 1},
			{Pos: 4, Row: 2, TokenPos: 10},
		},
	}
	positions := idx.GetPositions([]string{"hello"})
	if !reflect.DeepEqual(positions, expected) {
		t.Errorf("Expected positions %v, got %v", expected, positions)
	}
}

func TestMultipleTokens(t *testing.T) {
	idx := NewInvertedIndex()
	file1 := "file1.txt"
	idx.AddTokens(file1, []TokenizerResult{
		{Tokens: []string{"hello", "world", "hello", "world"}, Positions: []int{1, 10, 100, 1000}},
	})

	expected := map[string][]Position{
		file1: {
			{Pos: 1, Row: 1, TokenPos: 1},
			{Pos: 3, Row: 1, TokenPos: 100},
		},
	}
	positions := idx.GetPositions([]string{"hello", "world"})
	if !reflect.DeepEqual(positions, expected) {
		t.Errorf("Expected positions %v, got %v", expected, positions)
	}
}

func TestMultipleTokensNoMatch(t *testing.T) {
	idx := NewInvertedIndex()
	file1 := "file1.txt"
	idx.AddTokens(file1, []TokenizerResult{
		{Tokens: []string{"hello", "hello"}, Positions: []int{1, 10}},
	})
	positions := idx.GetPositions([]string{"hello", "world"})
	if len(positions) != 0 {
		t.Errorf("Expected no match for phrase 'hello world', got %v", positions)
	}
}
