package file

import (
	"os"
	"path/filepath"
	"testing"
)

func testdataDir(t *testing.T) string {
	dir := "testdata"
	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		t.Fatalf("testdata directory %s does not exist", dir)
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		t.Fatalf("cannot determine absolute path: %v", err)
	}
	return abs
}

func TestListFiles(t *testing.T) {
	dir := testdataDir(t)
	files, err := ListFiles(dir)
	if err != nil {
		t.Fatalf("Error listing files: %v", err)
	}
	if len(files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(files))
	}
	expected := []string{
		filepath.Join(dir, "test_file1.txt"),
		filepath.Join(dir, "subfolder", "test_file2.txt"),
	}
	for _, exp := range expected {
		found := false
		for _, f := range files {
			if f == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected file %s not found in list", exp)
		}
	}
}

func TestReadFile(t *testing.T) {
	fixturePath := filepath.Join(testdataDir(t), "test_file1.txt")
	results, err := ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}
	if len(results) != 4 {
		t.Errorf("Expected 4 lines, got %d", len(results))
	}

	if !equalStringSlices(results[0].Tokens, []string{"hello", "world"}) {
		t.Errorf("Line 1 tokens expected %v, got %v", []string{"hello", "world"}, results[0].Tokens)
	}
	if !equalStringSlices(results[1].Tokens, []string{"hello", "world"}) {
		t.Errorf("Line 2 tokens expected %v, got %v", []string{"hello", "world"}, results[1].Tokens)
	}
	if len(results[2].Tokens) != 0 {
		t.Errorf("Line 3 expected no tokens, got %v", results[2].Tokens)
	}
	if !equalStringSlices(results[3].Tokens, []string{"hello"}) {
		t.Errorf("Line 4 tokens expected %v, got %v", []string{"hello"}, results[3].Tokens)
	}

	if !equalIntSlices(results[0].Positions, []int{1, 7}) {
		t.Errorf("Line 1 positions expected %v, got %v", []int{1, 7}, results[0].Positions)
	}
	if !equalIntSlices(results[1].Positions, []int{1, 8}) {
		t.Errorf("Line 2 positions expected %v, got %v", []int{1, 8}, results[1].Positions)
	}
	if !equalIntSlices(results[3].Positions, []int{2}) {
		t.Errorf("Line 4 positions expected %v, got %v", []int{2}, results[3].Positions)
	}
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, s := range a {
		if s != b[i] {
			return false
		}
	}
	return true
}

func equalIntSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, n := range a {
		if n != b[i] {
			return false
		}
	}
	return true
}
