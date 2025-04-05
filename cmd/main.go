package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/martishin/search-engine/internal/file"
	"github.com/martishin/search-engine/internal/index"
	"github.com/martishin/search-engine/internal/query"
)

func main() {
	mode := flag.String("mode", "index", "mode: index or query")
	indexPath := flag.String("index", "index.json", "index file path")
	dir := flag.String("dir", "./", "directory for search")
	flag.Parse()

	switch *mode {
	case "index":
		idx := index.NewInvertedIndex()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		err := file.LoadFiles(ctx, *dir, idx)
		if err != nil {
			fmt.Printf("Error indexing files: %v\n", err)
			os.Exit(1)
		}
		err = index.DumpIndex(*indexPath, idx)
		if err != nil {
			fmt.Printf("Error dumping index: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Indexing complete. Index dumped to %s\n", *indexPath)
	case "query":
		idx, err := index.LoadIndex(*indexPath)
		if err != nil {
			fmt.Printf("Error loading index: %v\n", err)
			os.Exit(1)
		}
		query.ConsoleClient(idx)
	default:
		fmt.Println("Invalid mode. Use 'index' or 'query'.")
		os.Exit(1)
	}
}
