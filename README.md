# Search engine

On-disk full-text search engine written in Go. It consists of two main components:

- **Index Builder**: Recursively reads files from a specified directory, tokenizes their content, and builds an inverted index.
- **Search Query Executor**: Loads the generated index and executes search queries to find matching document positions.

The engine supports concurrent file processing and allows cancellation of the indexing process.   
It is designed with modularity in mind, splitting functionality into separate internal packages for file operations, text processing, indexing, and query handling.

## Features

- **Recursive File Indexing**: Walk through a directory tree and process every file (excluding directories).
- **Tokenization and Filtering**: Tokenizes file content using a simple space-based tokenizer, converts tokens to lower-case, and applies a simple stemming filter.
- **Inverted Index**: Builds an inverted index mapping tokens to their occurrences (including document name, row number, and token position within the line).
- **Concurrent Processing**: Uses goroutines and synchronization primitives to load files in parallel.
- **Cancellable Indexing**: The indexing process can be cancelled via a context (e.g., after a timeout).
- **Interactive Query Mode**: Run queries interactively from the command line.

## Usage

The project comes with two main modes: **index** and **query**. You can run the application using the Go command.

### Indexing Mode

This mode walks through a given directory, builds the index, and dumps it to a JSON file.

```bash
go run cmd/main.go --mode=index --index=index.json --dir=./data
```

### Query Mode

After building the index, run the query mode to search through the indexed documents:

```bash
go run cmd/main.go --mode=query --index=index.json
```

### Testing

Run the tests:
```bash
go test ./... -v
```

### Command-line Arguments

```
Usage:
  --mode         Mode to run: "index" or "query"
  --index        Path to the index file (e.g., index.json)
  --dir          Directory for search (used in indexing mode)
  -h, --help     Show this help message and exit
```
