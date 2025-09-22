# find-up

GoDoc GitHub Actions

find-up is a simple package to find files and directories by walking up parent directories or down descendant directories. It's a Go port of the popular [find-up](https://github.com/sindresorhus/find-up) npm package.

For more detail about the library and its features, reference your local godoc once installed.

Contributions welcome!

## Installation

```bash
go get github.com/viguza/find-up
```

## Available Functions

Core Functions. Some examples below:

| Function | Description | Example |
|----------|-------------|---------|
| `FindUp` | Find a file/directory by walking up parent directories | `FindUp("go.mod", nil)` |
| `FindUpMultiple` | Find multiple files/directories by walking up | `FindUpMultiple("*.go", options)` |
| `FindUpWithMatcher` | Find using a custom matcher function | `FindUpWithMatcher(matcher, options)` |
| `FindDown` | Find a file/directory by walking down descendant directories | `FindDown("*.test.go", options)` |
| `FindDownMultiple` | Find multiple files/directories by walking down | `FindDownMultiple("*.go", options)` |

## Features

* Find files and directories by walking up parent directories
* Find files and directories by walking down descendant directories
* Custom matcher functions for advanced searching
* Multiple file/directory search capabilities
* Configurable search options (depth, limits, types)
* Symbolic link handling
* Search strategy options (breadth-first, depth-first)
* Stop-at directory support
* Result limiting

## Examples

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    "github.com/viguza/find-up"
)

func main() {
    // Find a file by walking up parent directories
    result, err := findup.FindUp("go.mod", nil)
    if err != nil {
        log.Fatal(err)
    }
    if result != "" {
        fmt.Printf("Found go.mod at: %s\n", result)
    }
}
```

### Run Comprehensive Examples

To see all features in action with a realistic project structure:

```bash
# Run the examples
make example

# Or run directly
go run examples/main.go
```

The examples demonstrate:
- Basic file searching
- Glob pattern matching (`*.js`, `*.json`)
- Multiple file results
- Directory searching
- Custom matchers
- All available options (StopAt, Limit, Type, etc.)
- Real-world scenarios

### Find Multiple Files

```go
// Find multiple files with options
options := &findup.Options{
    Cwd:   ".",
    Limit: 5,
}
results, err := findup.FindUpMultiple("*.go", options)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Found %d Go files\n", len(results))
```

### Find Using Custom Matcher

```go
// Find a directory that contains a specific file
matcher := func(directory string) (string, bool, error) {
    goModPath := filepath.Join(directory, "go.mod")
    if _, err := os.Stat(goModPath); err == nil {
        return directory, true, nil
    }
    return "", false, nil
}

result, err := findup.FindUpWithMatcher(matcher, &findup.Options{
    Type: findup.DirectoryType,
})
```

### Find Down Descendant Directories

```go
// Find test files in descendant directories
options := &findup.Options{
    Cwd:   ".",
    Depth: 3,
}
result, err := findup.FindDown("*.test.go", options)
if err != nil {
    log.Fatal(err)
}
if result != "" {
    fmt.Printf("Found test file at: %s\n", result)
}
```

### Update Search Options

```go
// Create custom options
options := &findup.Options{
    Cwd:           "/path/to/start",
    Type:          findup.BothType, // Search for both files and directories
    AllowSymlinks: false,           // Don't follow symbolic links
    StopAt:        "/path/to/stop", // Stop searching at this directory
    Limit:         10,              // Limit results to 10
    Depth:         5,               // Limit search depth
    Strategy:      findup.BreadthFirst, // Use breadth-first search
}

result, err := findup.FindUp("config.json", options)
```

### Set Search Type

```go
// Search for directories only
options := &findup.Options{
    Type: findup.DirectoryType,
}
result, err := findup.FindUp("src", options)

// Search for files only (default)
options := &findup.Options{
    Type: findup.FileType,
}
result, err := findup.FindUp("*.go", options)

// Search for both files and directories
options := &findup.Options{
    Type: findup.BothType,
}
result, err := findup.FindUp("config", options)
```

### Use Different Search Strategies

```go
// Breadth-first search (default)
options := &findup.Options{
    Strategy: findup.BreadthFirst,
}

// Depth-first search
options := &findup.Options{
    Strategy: findup.DepthFirst,
}
```

### Chain and Pipe Output

```go
// Find multiple files and process them
results, err := findup.FindDownMultiple("*.go", &findup.Options{
    Cwd:   ".",
    Limit: 10,
})
if err != nil {
    log.Fatal(err)
}

for _, file := range results {
    fmt.Printf("Found Go file: %s\n", file)
}
```

### Get Search Status

```go
// Check if search found anything
result, err := findup.FindUp("config.json", nil)
if err != nil {
    log.Fatal(err)
}
if result != "" {
    fmt.Printf("Configuration found at: %s\n", result)
} else {
    fmt.Println("No configuration file found")
}
```

## Options

The `Options` struct provides configuration for all search functions:

```go
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
```

## Path Types

```go
const (
    FileType      PathType = iota // Search for files only
    DirectoryType                  // Search for directories only
    BothType                       // Search for both files and directories
)
```

## Search Strategies

```go
const (
    BreadthFirst SearchStrategy = iota // Breadth-first search
    DepthFirst                         // Depth-first search
)
```

## Matcher Functions

```go
type MatcherFunc func(directory string) (string, bool, error)
```

The matcher function receives the current directory path and should return:
- `string`: The path to return if this directory matches
- `bool`: Whether to stop searching (true) or continue (false)
- `error`: Any error that occurred

## Default Options

```go
func DefaultOptions() *Options
```

Returns default options with the following values:
- `Cwd`: "."
- `Type`: `FileType`
- `AllowSymlinks`: `true`
- `Limit`: -1 (no limit)
- `Depth`: 1
- `Strategy`: `BreadthFirst`

## Performance

The package is designed for efficiency:
- Uses `os.Stat` for file existence checks
- Supports early termination with `StopAt` option
- Configurable depth limits for `FindDown` functions
- Breadth-first and depth-first search strategies

## Error Handling

All functions return errors that should be checked:
- File system errors (permission denied, etc.)
- Invalid options
- Path resolution errors

## License

MIT

## Related

- [find-up](https://github.com/sindresorhus/find-up) - Original npm package
- [filepath](https://pkg.go.dev/path/filepath) - Go standard library for file paths
- [os](https://pkg.go.dev/os) - Go standard library for operating system interface