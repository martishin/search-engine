package text

import (
	"strings"
	"unicode"

	"github.com/martishin/search-engine/internal/index"
)

// SpaceTokenizer splits a string by non-letter characters and records the start position.
func SpaceTokenizer(s string) index.TokenizerResult {
	var tokens []string
	var positions []int
	var builder strings.Builder
	inToken := false
	startIdx := 0

	for i, r := range s {
		if unicode.IsLetter(r) {
			if !inToken {
				startIdx = i + 1 // 1-indexed
				inToken = true
			}
			builder.WriteRune(r)
		} else {
			if inToken {
				tokens = append(tokens, builder.String())
				positions = append(positions, startIdx)
				builder.Reset()
				inToken = false
			}
		}
	}
	if inToken {
		tokens = append(tokens, builder.String())
		positions = append(positions, startIdx)
	}
	return index.TokenizerResult{Tokens: tokens, Positions: positions}
}

// LowerCaseFilter converts tokens to lower-case.
func LowerCaseFilter(tokens []string) []string {
	var res []string
	for _, token := range tokens {
		res = append(res, strings.ToLower(token))
	}
	return res
}

// StemmingFilter is a simple placeholder that removes a trailing "s" if present.
func StemmingFilter(tokens []string) []string {
	var res []string
	for _, token := range tokens {
		if len(token) > 1 && token[len(token)-1] == 's' {
			token = token[:len(token)-1]
		}
		res = append(res, token)
	}
	return res
}

func ProcessText(s string) index.TokenizerResult {
	tr := SpaceTokenizer(s)
	tr.Tokens = LowerCaseFilter(tr.Tokens)
	tr.Tokens = StemmingFilter(tr.Tokens)
	return tr
}
