# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of find-up Go module
- Core functionality for finding files and directories by walking up parent directories
- Core functionality for finding files and directories by walking down descendant directories
- Support for custom matcher functions
- Multiple file/directory search capabilities
- Synchronous and asynchronous versions of all functions
- Comprehensive test coverage
- Performance benchmarks
- Full API documentation

### Features
- `FindUp` - Find files/directories by walking up parent directories
- `FindUpMultiple` - Find multiple files/directories by walking up parent directories
- `FindUpWithMatcher` - Find using custom matcher functions
- `FindDown` - Find files/directories by walking down descendant directories
- `FindDownMultiple` - Find multiple files/directories by walking down descendant directories
- Synchronous versions of all functions (`*Sync` variants)
- Configurable options for search behavior
- Support for different path types (file, directory, both)
- Symbolic link handling
- Search depth limits
- Search strategy options (breadth-first, depth-first)
- Result limiting
- Stop-at directory support

## [1.0.0] - 2024-01-XX

### Added
- Initial release
