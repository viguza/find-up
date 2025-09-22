package findup

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindUp(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "findup_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test directory structure
	// tempDir/
	//   ├── file1.txt
	//   ├── dir1/
	//   │   ├── file2.txt
	//   │   └── dir2/
	//   │       └── file3.txt
	//   └── dir3/
	//       └── file4.txt

	dir1 := filepath.Join(tempDir, "dir1")
	dir2 := filepath.Join(tempDir, "dir1", "dir2")
	dir3 := filepath.Join(tempDir, "dir3")

	err = os.MkdirAll(dir1, 0755)
	if err != nil {
		t.Fatalf("Failed to create dir1: %v", err)
	}
	err = os.MkdirAll(dir2, 0755)
	if err != nil {
		t.Fatalf("Failed to create dir2: %v", err)
	}
	err = os.MkdirAll(dir3, 0755)
	if err != nil {
		t.Fatalf("Failed to create dir3: %v", err)
	}

	// Create test files
	files := []string{
		filepath.Join(tempDir, "file1.txt"),
		filepath.Join(tempDir, "dir1", "file2.txt"),
		filepath.Join(tempDir, "dir1", "dir2", "file3.txt"),
		filepath.Join(tempDir, "dir3", "file4.txt"),
	}

	for _, file := range files {
		err = os.WriteFile(file, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", file, err)
		}
	}

	t.Run("FindUp from nested directory", func(t *testing.T) {
		// Test finding file1.txt from dir2
		options := &Options{Cwd: dir2}
		result, err := FindUp("file1.txt", options)
		if err != nil {
			t.Fatalf("FindUp failed: %v", err)
		}
		expected := filepath.Join(tempDir, "file1.txt")
		if result != expected {
			t.Errorf("Expected %s, got %s", expected, result)
		}
	})

	t.Run("FindUp from nested directory - file not found", func(t *testing.T) {
		// Test finding non-existent file
		options := &Options{Cwd: dir2}
		result, err := FindUp("nonexistent.txt", options)
		if err != nil {
			t.Fatalf("FindUp failed: %v", err)
		}
		if result != "" {
			t.Errorf("Expected empty result, got %s", result)
		}
	})

	t.Run("FindUp with directory type", func(t *testing.T) {
		// Test finding directory
		options := &Options{Cwd: dir2, Type: DirectoryType}
		result, err := FindUp("dir1", options)
		if err != nil {
			t.Fatalf("FindUp failed: %v", err)
		}
		expected := filepath.Join(tempDir, "dir1")
		if result != expected {
			t.Errorf("Expected %s, got %s", expected, result)
		}
	})

	t.Run("FindUp with both type", func(t *testing.T) {
		// Test finding both files and directories
		options := &Options{Cwd: dir2, Type: BothType}
		result, err := FindUp("dir1", options)
		if err != nil {
			t.Fatalf("FindUp failed: %v", err)
		}
		expected := filepath.Join(tempDir, "dir1")
		if result != expected {
			t.Errorf("Expected %s, got %s", expected, result)
		}
	})

	t.Run("FindUp with stopAt", func(t *testing.T) {
		// Test with stopAt option
		options := &Options{Cwd: dir2, StopAt: dir1}
		result, err := FindUp("file1.txt", options)
		if err != nil {
			t.Fatalf("FindUp failed: %v", err)
		}
		if result != "" {
			t.Errorf("Expected empty result (file1.txt is above stopAt), got %s", result)
		}
	})
}

func TestFindUpMultiple(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "findup_multiple_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test directory structure
	// tempDir/
	//   ├── file1.txt
	//   ├── dir1/
	//   │   ├── file1.txt (duplicate name)
	//   │   └── dir2/
	//   │       └── file2.txt
	//   └── dir3/
	//       └── file1.txt (duplicate name)

	dir1 := filepath.Join(tempDir, "dir1")
	dir2 := filepath.Join(tempDir, "dir1", "dir2")
	dir3 := filepath.Join(tempDir, "dir3")

	err = os.MkdirAll(dir1, 0755)
	if err != nil {
		t.Fatalf("Failed to create dir1: %v", err)
	}
	err = os.MkdirAll(dir2, 0755)
	if err != nil {
		t.Fatalf("Failed to create dir2: %v", err)
	}
	err = os.MkdirAll(dir3, 0755)
	if err != nil {
		t.Fatalf("Failed to create dir3: %v", err)
	}

	// Create test files
	files := []string{
		filepath.Join(tempDir, "file1.txt"),
		filepath.Join(tempDir, "dir1", "file1.txt"),
		filepath.Join(tempDir, "dir1", "dir2", "file2.txt"),
		filepath.Join(tempDir, "dir3", "file1.txt"),
	}

	for _, file := range files {
		err = os.WriteFile(file, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", file, err)
		}
	}

	t.Run("FindUpMultiple from nested directory", func(t *testing.T) {
		// Test finding multiple file1.txt files
		options := &Options{Cwd: dir2}
		results, err := FindUpMultiple("file1.txt", options)
		if err != nil {
			t.Fatalf("FindUpMultiple failed: %v", err)
		}
		if len(results) != 2 {
			t.Errorf("Expected 2 results, got %d", len(results))
		}

		expected1 := filepath.Join(tempDir, "dir1", "file1.txt")
		expected2 := filepath.Join(tempDir, "file1.txt")
		found1 := false
		found2 := false
		for _, result := range results {
			if result == expected1 {
				found1 = true
			}
			if result == expected2 {
				found2 = true
			}
		}
		if !found1 || !found2 {
			t.Errorf("Expected to find both %s and %s, got %v", expected1, expected2, results)
		}
	})

	t.Run("FindUpMultiple with limit", func(t *testing.T) {
		// Test with limit option
		options := &Options{Cwd: dir2, Limit: 1}
		results, err := FindUpMultiple("file1.txt", options)
		if err != nil {
			t.Fatalf("FindUpMultiple failed: %v", err)
		}
		if len(results) != 1 {
			t.Errorf("Expected 1 result due to limit, got %d", len(results))
		}
	})
}

func TestFindUpWithMatcher(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "findup_matcher_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test directory structure
	// tempDir/
	//   ├── file1.txt
	//   ├── dir1/
	//   │   ├── file2.txt
	//   │   └── dir2/
	//   │       └── file3.txt
	//   └── dir3/
	//       └── file4.txt

	dir1 := filepath.Join(tempDir, "dir1")
	dir2 := filepath.Join(tempDir, "dir1", "dir2")
	dir3 := filepath.Join(tempDir, "dir3")

	err = os.MkdirAll(dir1, 0755)
	if err != nil {
		t.Fatalf("Failed to create dir1: %v", err)
	}
	err = os.MkdirAll(dir2, 0755)
	if err != nil {
		t.Fatalf("Failed to create dir2: %v", err)
	}
	err = os.MkdirAll(dir3, 0755)
	if err != nil {
		t.Fatalf("Failed to create dir3: %v", err)
	}

	// Create test files
	files := []string{
		filepath.Join(tempDir, "file1.txt"),
		filepath.Join(tempDir, "dir1", "file2.txt"),
		filepath.Join(tempDir, "dir1", "dir2", "file3.txt"),
		filepath.Join(tempDir, "dir3", "file4.txt"),
	}

	for _, file := range files {
		err = os.WriteFile(file, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", file, err)
		}
	}

	t.Run("FindUpWithMatcher - find directory with specific file", func(t *testing.T) {
		// Matcher function that looks for a directory containing file1.txt
		matcher := func(directory string) (string, bool, error) {
			file1Path := filepath.Join(directory, "file1.txt")
			if _, err := os.Stat(file1Path); err == nil {
				return directory, true, nil
			}
			return "", false, nil
		}

		options := &Options{Cwd: dir2}
		result, err := FindUpWithMatcher(matcher, options)
		if err != nil {
			t.Fatalf("FindUpWithMatcher failed: %v", err)
		}
		if result != tempDir {
			t.Errorf("Expected %s, got %s", tempDir, result)
		}
	})

	t.Run("FindUpWithMatcher - no match found", func(t *testing.T) {
		// Matcher function that looks for a non-existent file
		matcher := func(directory string) (string, bool, error) {
			filePath := filepath.Join(directory, "nonexistent.txt")
			if _, err := os.Stat(filePath); err == nil {
				return directory, true, nil
			}
			return "", false, nil
		}

		options := &Options{Cwd: dir2}
		result, err := FindUpWithMatcher(matcher, options)
		if err != nil {
			t.Fatalf("FindUpWithMatcher failed: %v", err)
		}
		if result != "" {
			t.Errorf("Expected empty result, got %s", result)
		}
	})
}

func TestFindDown(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "finddown_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test directory structure
	// tempDir/
	//   ├── file1.txt
	//   ├── dir1/
	//   │   ├── file2.txt
	//   │   └── dir2/
	//   │       └── file3.txt
	//   └── dir3/
	//       └── file4.txt

	dir1 := filepath.Join(tempDir, "dir1")
	dir2 := filepath.Join(tempDir, "dir1", "dir2")
	dir3 := filepath.Join(tempDir, "dir3")

	err = os.MkdirAll(dir1, 0755)
	if err != nil {
		t.Fatalf("Failed to create dir1: %v", err)
	}
	err = os.MkdirAll(dir2, 0755)
	if err != nil {
		t.Fatalf("Failed to create dir2: %v", err)
	}
	err = os.MkdirAll(dir3, 0755)
	if err != nil {
		t.Fatalf("Failed to create dir3: %v", err)
	}

	// Create test files
	files := []string{
		filepath.Join(tempDir, "file1.txt"),
		filepath.Join(tempDir, "dir1", "file2.txt"),
		filepath.Join(tempDir, "dir1", "dir2", "file3.txt"),
		filepath.Join(tempDir, "dir3", "file4.txt"),
	}

	for _, file := range files {
		err = os.WriteFile(file, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", file, err)
		}
	}

	t.Run("FindDown from root directory", func(t *testing.T) {
		// Test finding file3.txt from tempDir
		options := &Options{Cwd: tempDir}
		result, err := FindDown("file3.txt", options)
		if err != nil {
			t.Fatalf("FindDown failed: %v", err)
		}
		expected := filepath.Join(tempDir, "dir1", "dir2", "file3.txt")
		if result != expected {
			t.Errorf("Expected %s, got %s", expected, result)
		}
	})

	t.Run("FindDown with depth limit", func(t *testing.T) {
		// Test with depth limit
		options := &Options{Cwd: tempDir, Depth: 1}
		result, err := FindDown("file3.txt", options)
		if err != nil {
			t.Fatalf("FindDown failed: %v", err)
		}
		if result != "" {
			t.Errorf("Expected empty result (file3.txt is deeper than depth 1), got %s", result)
		}
	})

	t.Run("FindDown with directory type", func(t *testing.T) {
		// Test finding directory
		options := &Options{Cwd: tempDir, Type: DirectoryType}
		result, err := FindDown("dir1", options)
		if err != nil {
			t.Fatalf("FindDown failed: %v", err)
		}
		expected := filepath.Join(tempDir, "dir1")
		if result != expected {
			t.Errorf("Expected %s, got %s", expected, result)
		}
	})
}

func TestFindDownMultiple(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "finddown_multiple_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test directory structure
	// tempDir/
	//   ├── file1.txt
	//   ├── dir1/
	//   │   ├── file1.txt (duplicate name)
	//   │   └── dir2/
	//   │       └── file1.txt (duplicate name)
	//   └── dir3/
	//       └── file1.txt (duplicate name)

	dir1 := filepath.Join(tempDir, "dir1")
	dir2 := filepath.Join(tempDir, "dir1", "dir2")
	dir3 := filepath.Join(tempDir, "dir3")

	err = os.MkdirAll(dir1, 0755)
	if err != nil {
		t.Fatalf("Failed to create dir1: %v", err)
	}
	err = os.MkdirAll(dir2, 0755)
	if err != nil {
		t.Fatalf("Failed to create dir2: %v", err)
	}
	err = os.MkdirAll(dir3, 0755)
	if err != nil {
		t.Fatalf("Failed to create dir3: %v", err)
	}

	// Create test files
	files := []string{
		filepath.Join(tempDir, "file1.txt"),
		filepath.Join(tempDir, "dir1", "file1.txt"),
		filepath.Join(tempDir, "dir1", "dir2", "file1.txt"),
		filepath.Join(tempDir, "dir3", "file1.txt"),
	}

	for _, file := range files {
		err = os.WriteFile(file, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", file, err)
		}
	}

	t.Run("FindDownMultiple from root directory", func(t *testing.T) {
		// Test finding multiple file1.txt files
		options := &Options{Cwd: tempDir}
		results, err := FindDownMultiple("file1.txt", options)
		if err != nil {
			t.Fatalf("FindDownMultiple failed: %v", err)
		}
		if len(results) != 4 {
			t.Errorf("Expected 4 results, got %d", len(results))
		}
	})

	t.Run("FindDownMultiple with limit", func(t *testing.T) {
		// Test with limit option
		options := &Options{Cwd: tempDir, Limit: 2}
		results, err := FindDownMultiple("file1.txt", options)
		if err != nil {
			t.Fatalf("FindDownMultiple failed: %v", err)
		}
		if len(results) != 2 {
			t.Errorf("Expected 2 results due to limit, got %d", len(results))
		}
	})
}

func TestDefaultOptions(t *testing.T) {
	options := DefaultOptions()
	if options.Cwd != "." {
		t.Errorf("Expected Cwd to be '.', got %s", options.Cwd)
	}
	if options.Type != FileType {
		t.Errorf("Expected Type to be FileType, got %v", options.Type)
	}
	if !options.AllowSymlinks {
		t.Error("Expected AllowSymlinks to be true")
	}
	if options.Limit != -1 {
		t.Errorf("Expected Limit to be -1, got %d", options.Limit)
	}
	if options.Depth != 1 {
		t.Errorf("Expected Depth to be 1, got %d", options.Depth)
	}
	if options.Strategy != BreadthFirst {
		t.Errorf("Expected Strategy to be BreadthFirst, got %v", options.Strategy)
	}
}

func TestPathType(t *testing.T) {
	if FileType != 0 {
		t.Error("Expected FileType to be 0")
	}
	if DirectoryType != 1 {
		t.Error("Expected DirectoryType to be 1")
	}
	if BothType != 2 {
		t.Error("Expected BothType to be 2")
	}
}

func TestSearchStrategy(t *testing.T) {
	if BreadthFirst != 0 {
		t.Error("Expected BreadthFirst to be 0")
	}
	if DepthFirst != 1 {
		t.Error("Expected DepthFirst to be 1")
	}
}
