package text

import (
	"reflect"
	"testing"

	"github.com/martishin/search-engine/internal/index"
)

func TestSpaceTokenizer(t *testing.T) {
	tests := []struct {
		input          string
		expectedTokens []string
		expectedPos    []int
	}{
		{"Hello", []string{"Hello"}, []int{1}},
		{"Hello  World", []string{"Hello", "World"}, []int{1, 8}},
		{"Hello, World!", []string{"Hello", "World"}, []int{1, 8}},
		{" Hello ", []string{"Hello"}, []int{2}},
		{"test.Hello", []string{"test", "Hello"}, []int{1, 6}},
	}
	for _, tt := range tests {
		result := SpaceTokenizer(tt.input)
		if !reflect.DeepEqual(result.Tokens, tt.expectedTokens) {
			t.Errorf("SpaceTokenizer(%q): expected tokens %v, got %v", tt.input, tt.expectedTokens, result.Tokens)
		}
		if !reflect.DeepEqual(result.Positions, tt.expectedPos) {
			t.Errorf("SpaceTokenizer(%q): expected positions %v, got %v", tt.input, tt.expectedPos, result.Positions)
		}
	}
}

func TestProcessText(t *testing.T) {
	result := ProcessText("Hello  World")
	expected := index.TokenizerResult{Tokens: []string{"hello", "world"}, Positions: []int{1, 8}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ProcessText expected %v, got %v", expected, result)
	}

	result2 := ProcessText("Hello Worlds World")
	expected2 := index.TokenizerResult{Tokens: []string{"hello", "world", "world"}, Positions: []int{1, 7, 14}}
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("ProcessText modified expected %v, got %v", expected2, result2)
	}
}

func TestLowerCaseFilter(t *testing.T) {
	input := []string{"Hello", "World"}
	expected := []string{"hello", "world"}
	result := LowerCaseFilter(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("LowerCaseFilter expected %v, got %v", expected, result)
	}
}

func TestStemmingFilter(t *testing.T) {
	input := []string{"hello", "worlds", "world"}
	expected := []string{"hello", "world", "world"}
	result := StemmingFilter(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("StemmingFilter expected %v, got %v", expected, result)
	}
}
