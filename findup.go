// Package findup provides utilities for finding files and directories by walking up or down the directory tree.
// It's a Go port of the popular npm package find-up.
package findup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// PathType represents the type of path to search for
type PathType int

const (
	// FileType searches for files only
	FileType PathType = iota
	// DirectoryType searches for directories only
	DirectoryType
	// BothType searches for both files and directories
	BothType
)

// Options contains configuration options for find operations
type Options struct {
	// Cwd is the directory to start from (default: current working directory)
	Cwd string
	// Type specifies the type of path to match
	Type PathType
	// AllowSymlinks determines if symbolic links should be matched
	AllowSymlinks bool
	// StopAt is the directory where the search halts (only for findUp functions)
	StopAt string
	// Limit is the maximum number of matches to return (only for findUpMultiple functions)
	Limit int
	// Depth is the maximum number of directory levels to traverse (only for findDown functions)
	Depth int
	// Strategy determines the search strategy for findDown functions
	Strategy SearchStrategy
}

// SearchStrategy represents the search strategy for findDown functions
type SearchStrategy int

const (
	// BreadthFirst performs breadth-first search
	BreadthFirst SearchStrategy = iota
	// DepthFirst performs depth-first search
	DepthFirst
)

// MatcherFunc is a function that determines if a directory matches the search criteria
type MatcherFunc func(directory string) (string, bool, error)

// DefaultOptions returns default options
func DefaultOptions() *Options {
	return &Options{
		Cwd:           ".",
		Type:          FileType,
		AllowSymlinks: true,
		Limit:         -1, // -1 means no limit
		Depth:         1,
		Strategy:      BreadthFirst,
	}
}

// FindUp finds a file or directory by walking up parent directories
func FindUp(name string, options *Options) (string, error) {
	if options == nil {
		options = DefaultOptions()
	}

	opts := *options
	if opts.Cwd == "" {
		opts.Cwd = "."
	}

	// Convert to absolute path
	absCwd, err := filepath.Abs(opts.Cwd)
	if err != nil {
		return "", err
	}

	stopAt := opts.StopAt
	if stopAt != "" {
		stopAt, err = filepath.Abs(stopAt)
		if err != nil {
			return "", err
		}
	}

	return findUpInDir(absCwd, name, &opts, stopAt)
}

// FindUpMultiple finds multiple files or directories by walking up parent directories
func FindUpMultiple(name string, options *Options) ([]string, error) {
	if options == nil {
		options = DefaultOptions()
	}

	opts := *options
	if opts.Cwd == "" {
		opts.Cwd = "."
	}

	absCwd, err := filepath.Abs(opts.Cwd)
	if err != nil {
		return nil, err
	}

	stopAt := opts.StopAt
	if stopAt != "" {
		stopAt, err = filepath.Abs(stopAt)
		if err != nil {
			return nil, err
		}
	}

	var results []string
	err = findUpMultipleInDir(absCwd, name, &opts, stopAt, &results)
	return results, err
}

// FindUpWithMatcher finds a file or directory using a custom matcher function
func FindUpWithMatcher(matcher MatcherFunc, options *Options) (string, error) {
	if options == nil {
		options = DefaultOptions()
	}

	opts := *options
	if opts.Cwd == "" {
		opts.Cwd = "."
	}

	absCwd, err := filepath.Abs(opts.Cwd)
	if err != nil {
		return "", err
	}

	stopAt := opts.StopAt
	if stopAt != "" {
		stopAt, err = filepath.Abs(stopAt)
		if err != nil {
			return "", err
		}
	}

	return findUpWithMatcherInDir(absCwd, matcher, &opts, stopAt)
}

// FindDown finds a file or directory by walking down descendant directories
func FindDown(name string, options *Options) (string, error) {
	if options == nil {
		options = DefaultOptions()
	}

	opts := *options
	if opts.Cwd == "" {
		opts.Cwd = "."
	}

	absCwd, err := filepath.Abs(opts.Cwd)
	if err != nil {
		return "", err
	}

	return findDownInDir(absCwd, name, &opts, 0)
}

// FindDownMultiple finds multiple files or directories by walking down descendant directories
func FindDownMultiple(name string, options *Options) ([]string, error) {
	if options == nil {
		options = DefaultOptions()
	}

	opts := *options
	if opts.Cwd == "" {
		opts.Cwd = "."
	}

	absCwd, err := filepath.Abs(opts.Cwd)
	if err != nil {
		return nil, err
	}

	var results []string
	err = findDownMultipleInDir(absCwd, name, &opts, 0, &results)
	return results, err
}

// Helper functions

// isGlobPattern checks if the name contains glob patterns
func isGlobPattern(name string) bool {
	return strings.Contains(name, "*") || strings.Contains(name, "?") || strings.Contains(name, "[")
}

// matchesGlob checks if a file matches a glob pattern
func matchesGlob(filename, pattern string) (bool, error) {
	matched, err := filepath.Match(pattern, filename)
	return matched, err
}

func findUpInDir(dir, name string, options *Options, stopAt string) (string, error) {
	current := dir

	for {
		// Check if we should stop at this directory
		if stopAt != "" && current == stopAt {
			break
		}

		// Check if the target exists in current directory
		if isGlobPattern(name) {
			// Handle glob patterns by listing directory contents
			entries, err := os.ReadDir(current)
			if err == nil {
				for _, entry := range entries {
					entryName := entry.Name()
					if matched, err := matchesGlob(entryName, name); err == nil && matched {
						target := filepath.Join(current, entryName)
						if matches, err := pathMatches(target, options); err == nil && matches {
							return target, nil
						}
					}
				}
			}
		} else {
			// Handle exact filename match
			target := filepath.Join(current, name)
			if matches, err := pathMatches(target, options); err == nil && matches {
				return target, nil
			}
		}

		// Move to parent directory
		parent := filepath.Dir(current)
		if parent == current {
			// Reached root directory
			break
		}
		current = parent
	}

	return "", nil
}

func findUpMultipleInDir(dir, name string, options *Options, stopAt string, results *[]string) error {
	current := dir

	for {
		// Check if we should stop at this directory
		if stopAt != "" && current == stopAt {
			break
		}

		// Check if the target exists in current directory
		if isGlobPattern(name) {
			// Handle glob patterns by listing directory contents
			entries, err := os.ReadDir(current)
			if err == nil {
				for _, entry := range entries {
					entryName := entry.Name()
					if matched, err := matchesGlob(entryName, name); err == nil && matched {
						target := filepath.Join(current, entryName)
						if matches, err := pathMatches(target, options); err == nil && matches {
							*results = append(*results, target)

							// Check if we've reached the limit
							if options.Limit > 0 && len(*results) >= options.Limit {
								return nil
							}
						}
					}
				}
			}
		} else {
			// Handle exact filename match
			target := filepath.Join(current, name)
			if matches, err := pathMatches(target, options); err == nil && matches {
				*results = append(*results, target)

				// Check if we've reached the limit
				if options.Limit > 0 && len(*results) >= options.Limit {
					return nil
				}
			}
		}

		// Move to parent directory
		parent := filepath.Dir(current)
		if parent == current {
			// Reached root directory
			break
		}
		current = parent
	}

	return nil
}

func findUpWithMatcherInDir(dir string, matcher MatcherFunc, options *Options, stopAt string) (string, error) {
	current := dir

	for {
		// Check if we should stop at this directory
		if stopAt != "" && current == stopAt {
			break
		}

		// Call the matcher function
		result, shouldStop, err := matcher(current)
		if err != nil {
			return "", err
		}

		if shouldStop {
			return result, nil
		}

		// Move to parent directory
		parent := filepath.Dir(current)
		if parent == current {
			// Reached root directory
			break
		}
		current = parent
	}

	return "", nil
}

func findDownInDir(dir, name string, options *Options, currentDepth int) (string, error) {
	// Check if we've exceeded the depth limit
	if options.Depth > 0 && currentDepth > options.Depth {
		return "", nil
	}

	// Check if the target exists in current directory
	if isGlobPattern(name) {
		// Handle glob patterns by listing directory contents
		entries, err := os.ReadDir(dir)
		if err == nil {
			for _, entry := range entries {
				entryName := entry.Name()
				if matched, err := matchesGlob(entryName, name); err == nil && matched {
					target := filepath.Join(dir, entryName)
					if matches, err := pathMatches(target, options); err == nil && matches {
						return target, nil
					}
				}
			}
		}
	} else {
		// Handle exact filename match
		target := filepath.Join(dir, name)
		if matches, err := pathMatches(target, options); err == nil && matches {
			return target, nil
		}
	}

	// Read directory contents
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	// Collect subdirectories
	var subdirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			subdirs = append(subdirs, filepath.Join(dir, entry.Name()))
		}
	}

	// Search subdirectories based on strategy
	if options.Strategy == BreadthFirst {
		// Breadth-first: search all subdirectories at current level first
		for _, subdir := range subdirs {
			if result, err := findDownInDir(subdir, name, options, currentDepth+1); err == nil && result != "" {
				return result, nil
			}
		}
	} else {
		// Depth-first: search each subdirectory completely before moving to next
		for _, subdir := range subdirs {
			if result, err := findDownInDir(subdir, name, options, currentDepth+1); err == nil && result != "" {
				return result, nil
			}
		}
	}

	return "", nil
}

func findDownMultipleInDir(dir, name string, options *Options, currentDepth int, results *[]string) error {
	// Check if we've exceeded the depth limit
	if options.Depth > 0 && currentDepth > options.Depth {
		return nil
	}

	// Check if the target exists in current directory
	if isGlobPattern(name) {
		// Handle glob patterns by listing directory contents
		entries, err := os.ReadDir(dir)
		if err == nil {
			for _, entry := range entries {
				entryName := entry.Name()
				if matched, err := matchesGlob(entryName, name); err == nil && matched {
					target := filepath.Join(dir, entryName)
					if matches, err := pathMatches(target, options); err == nil && matches {
						*results = append(*results, target)

						// Check if we've reached the limit
						if options.Limit > 0 && len(*results) >= options.Limit {
							return nil
						}
					}
				}
			}
		}
	} else {
		// Handle exact filename match
		target := filepath.Join(dir, name)
		if matches, err := pathMatches(target, options); err == nil && matches {
			*results = append(*results, target)

			// Check if we've reached the limit
			if options.Limit > 0 && len(*results) >= options.Limit {
				return nil
			}
		}
	}

	// Read directory contents
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	// Collect subdirectories
	var subdirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			subdirs = append(subdirs, filepath.Join(dir, entry.Name()))
		}
	}

	// Search subdirectories
	for _, subdir := range subdirs {
		if err := findDownMultipleInDir(subdir, name, options, currentDepth+1, results); err != nil {
			return err
		}

		// Check if we've reached the limit
		if options.Limit > 0 && len(*results) >= options.Limit {
			return nil
		}
	}

	return nil
}

func pathMatches(path string, options *Options) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	// Check if it's a symlink
	if info.Mode()&os.ModeSymlink != 0 {
		if !options.AllowSymlinks {
			return false, nil
		}

		// Resolve the symlink
		resolved, err := os.Readlink(path)
		if err != nil {
			return false, err
		}

		// Make path absolute if it's relative
		if !filepath.IsAbs(resolved) {
			resolved = filepath.Join(filepath.Dir(path), resolved)
		}

		// Check the resolved path
		resolvedInfo, err := os.Stat(resolved)
		if err != nil {
			return false, err
		}
		info = resolvedInfo
	}

	// Check the type
	switch options.Type {
	case FileType:
		return !info.IsDir(), nil
	case DirectoryType:
		return info.IsDir(), nil
	case BothType:
		return true, nil
	default:
		return false, fmt.Errorf("invalid path type: %v", options.Type)
	}
}
