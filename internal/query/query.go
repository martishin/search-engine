package query

import (
	"bufio"
	"fmt"
	"os"

	"github.com/martishin/search-engine/internal/index"
	"github.com/martishin/search-engine/internal/text"
)

// ProcessQuery tokenizes the query string, retrieves matching positions from the index,
// and returns a slice of SearchResult.
func ProcessQuery(idx *index.InvertedIndex, queryStr string) []index.SearchResult {
	fmt.Printf("Started processing query: %s\n", queryStr)
	tr := text.ProcessText(queryStr)
	docPositions := idx.GetPositions(tr.Tokens)
	var results []index.SearchResult
	for doc, positions := range docPositions {
		results = append(results, index.SearchResult{
			Document:  doc,
			Positions: positions,
		})
	}
	fmt.Printf("Finished processing query: %s\n", queryStr)
	return results
}

// ConsoleClient reads queries from stdin and prints the search results.
func ConsoleClient(idx *index.InvertedIndex) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter your search queries (press Ctrl+C to exit):")
	for scanner.Scan() {
		queryStr := scanner.Text()
		results := ProcessQuery(idx, queryStr)
		fmt.Printf("%s: %v\n", queryStr, results)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
}
