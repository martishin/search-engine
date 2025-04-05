package index

import (
	"fmt"
	"strings"
	"sync"
)

// Position represents a token position in a document.
type Position struct {
	Pos      int `json:"pos"`      // Sequential position in the file.
	Row      int `json:"row"`      // Line number (1-indexed).
	TokenPos int `json:"tokenPos"` // Start index in the line (1-indexed).
}

// TokenizerResult holds the tokens from a line along with their positions.
type TokenizerResult struct {
	Tokens    []string `json:"tokens"`
	Positions []int    `json:"positions"`
}

// SearchResult holds the search result for a document.
type SearchResult struct {
	Document  string     `json:"document"`
	Positions []Position `json:"positions"`
}

// String formats a search result showing only row and tokenPos.
func (sr SearchResult) String() string {
	var posStrings []string
	for _, pos := range sr.Positions {
		posStrings = append(posStrings, fmt.Sprintf("{row: %d, pos: %d}", pos.Row, pos.TokenPos))
	}
	return fmt.Sprintf("SearchResult(document=%s, positions=[%s])", sr.Document, strings.Join(posStrings, ", "))
}

// InvertedIndex stores a mapping from token to document positions.
// It is safe for concurrent use.
type InvertedIndex struct {
	index map[string]map[string][]Position
	mu    sync.RWMutex
}

// NewInvertedIndex creates an empty inverted index.
func NewInvertedIndex() *InvertedIndex {
	return &InvertedIndex{
		index: make(map[string]map[string][]Position),
	}
}

// AddTokens adds tokens from a file to the index.
func (idx *InvertedIndex) AddTokens(filePath string, tokens []TokenizerResult) {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	position := 1
	for row, tr := range tokens {
		for i, token := range tr.Tokens {
			if idx.index[token] == nil {
				idx.index[token] = make(map[string][]Position)
			}
			idx.index[token][filePath] = append(idx.index[token][filePath], Position{
				Pos:      position,
				Row:      row + 1,
				TokenPos: tr.Positions[i],
			})
			position++
		}
	}
}

func (idx *InvertedIndex) getTokenPositions(token string) map[string][]Position {
	if m, ok := idx.index[token]; ok {
		return m
	}
	return make(map[string][]Position)
}

// GetPositions returns the positions for a sequence of query tokens,
// attempting to find consecutive tokens (phrase search).
func (idx *InvertedIndex) GetPositions(queryTokens []string) map[string][]Position {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	result := make(map[string][]Position)
	if len(queryTokens) == 0 {
		return result
	}

	// Start with the first token.
	prev := idx.getTokenPositions(queryTokens[0])
	parentPositions := make(map[Position]Position)
	for _, positions := range prev {
		for _, pos := range positions {
			parentPositions[pos] = pos
		}
	}

	// For each subsequent token, look for consecutive positions.
	for i := 1; i < len(queryTokens); i++ {
		curr := idx.getTokenPositions(queryTokens[i])
		continuous := make(map[string][]Position)
		for doc, prevPositions := range prev {
			currPositions, exists := curr[doc]
			if !exists {
				continue
			}
			pp, cp := 0, 0
			for pp < len(prevPositions) && cp < len(currPositions) {
				if currPositions[cp].Pos == prevPositions[pp].Pos+1 {
					continuous[doc] = append(continuous[doc], currPositions[cp])
					parentPositions[currPositions[cp]] = parentPositions[prevPositions[pp]]
					pp++
					cp++
				} else if currPositions[cp].Pos < prevPositions[pp].Pos+1 {
					cp++
				} else {
					pp++
				}
			}
		}
		prev = continuous
	}

	// Reconstruct positions using parentPositions.
	for doc, positions := range prev {
		var finalPositions []Position
		for _, pos := range positions {
			if parent, exists := parentPositions[pos]; exists {
				finalPositions = append(finalPositions, parent)
			} else {
				finalPositions = append(finalPositions, pos)
			}
		}
		result[doc] = finalPositions
	}
	return result
}
