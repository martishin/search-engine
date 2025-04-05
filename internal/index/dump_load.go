package index

import (
	"encoding/json"
	"os"
)

// DumpIndex serializes the index to a JSON file.
func DumpIndex(filePath string, idx *InvertedIndex) error {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	return encoder.Encode(idx.index)
}

// LoadIndex loads the index from a JSON file.
func LoadIndex(filePath string) (*InvertedIndex, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var data map[string]map[string][]Position
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}
	return &InvertedIndex{
		index: data,
	}, nil
}
