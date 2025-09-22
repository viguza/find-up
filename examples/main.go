package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	findup "github.com/viguza/find-up"
)

func main() {
	fmt.Println("=== find-up Go Module Examples ===")
	fmt.Println("Demonstrating the find-up Go module functionality...")

	// Create example directory structure
	tempDir, err := os.MkdirTemp("", "findup_examples")
	if err != nil {
		log.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	fmt.Printf("Created example directory: %s\n\n", tempDir)

	// Setup example structure
	setupExampleStructure(tempDir)

	// Example 1: Basic FindUp
	fmt.Println("1. Basic FindUp - Find files by walking up directories")
	fmt.Println("   Looking for 'config.json' from a nested directory...")

	nestedDir := filepath.Join(tempDir, "project", "src", "components", "ui")
	result, err := findup.FindUp("config.json", &findup.Options{Cwd: nestedDir})
	if err != nil {
		log.Printf("Error: %v", err)
	} else if result != "" {
		fmt.Printf("   ✅ Found: %s\n", result)
	} else {
		fmt.Println("   ℹ️  No config.json found")
	}

	// Example 2: FindUp with glob patterns
	fmt.Println("\n2. FindUp with Glob Patterns - Find files matching patterns")
	fmt.Println("   Looking for '*.json' files...")

	result, err = findup.FindUp("*.json", &findup.Options{Cwd: nestedDir})
	if err != nil {
		log.Printf("Error: %v", err)
	} else if result != "" {
		fmt.Printf("   ✅ Found: %s\n", result)
	} else {
		fmt.Println("   ℹ️  No .json files found")
	}

	// Example 3: FindUpMultiple
	fmt.Println("\n3. FindUpMultiple - Find all matching files")
	fmt.Println("   Looking for all '*.js' files...")

	results, err := findup.FindUpMultiple("*.js", &findup.Options{Cwd: nestedDir})
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("   ✅ Found %d .js files:\n", len(results))
		for i, file := range results {
			fmt.Printf("      %d. %s\n", i+1, file)
		}
	}

	// Example 4: FindDown
	fmt.Println("\n4. FindDown - Find files in descendant directories")
	fmt.Println("   Looking for '*.md' files in project directory...")

	results, err = findup.FindDownMultiple("*.md", &findup.Options{
		Cwd:   filepath.Join(tempDir, "project"),
		Depth: 3,
	})
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("   ✅ Found %d .md files:\n", len(results))
		for i, file := range results {
			fmt.Printf("      %d. %s\n", i+1, file)
		}
	}

	// Example 5: Custom matcher
	fmt.Println("\n5. Custom Matcher - Find directories containing specific files")
	fmt.Println("   Looking for directories containing 'package.json'...")

	matcher := func(directory string) (string, bool, error) {
		packagePath := filepath.Join(directory, "package.json")
		if _, err := os.Stat(packagePath); err == nil {
			return directory, true, nil
		}
		return "", false, nil
	}

	result, err = findup.FindUpWithMatcher(matcher, &findup.Options{
		Cwd:  nestedDir,
		Type: findup.DirectoryType,
	})
	if err != nil {
		log.Printf("Error: %v", err)
	} else if result != "" {
		fmt.Printf("   ✅ Found directory: %s\n", result)
	} else {
		fmt.Println("   ℹ️  No directory with package.json found")
	}

	// Example 6: Options - StopAt
	fmt.Println("\n6. Options - StopAt - Limit search to specific directory")
	fmt.Println("   Looking for 'config.json' but stopping at 'src' directory...")

	stopAtDir := filepath.Join(tempDir, "project", "src")
	result, err = findup.FindUp("config.json", &findup.Options{
		Cwd:    nestedDir,
		StopAt: stopAtDir,
	})
	if err != nil {
		log.Printf("Error: %v", err)
	} else if result != "" {
		fmt.Printf("   ✅ Found: %s\n", result)
	} else {
		fmt.Println("   ℹ️  No config.json found within src directory")
	}

	// Example 7: Options - Type filtering
	fmt.Println("\n7. Options - Type Filtering - Find only directories")
	fmt.Println("   Looking for 'node_modules' directory...")

	result, err = findup.FindUp("node_modules", &findup.Options{
		Cwd:  nestedDir,
		Type: findup.DirectoryType,
	})
	if err != nil {
		log.Printf("Error: %v", err)
	} else if result != "" {
		fmt.Printf("   ✅ Found directory: %s\n", result)
	} else {
		fmt.Println("   ℹ️  No node_modules directory found")
	}

	// Example 8: Options - Limit results
	fmt.Println("\n8. Options - Limit - Limit number of results")
	fmt.Println("   Looking for all .txt files but limiting to 2 results...")

	results, err = findup.FindUpMultiple("*.txt", &findup.Options{
		Cwd:   nestedDir,
		Limit: 2,
	})
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("   ✅ Found %d .txt files (limited to 2):\n", len(results))
		for i, file := range results {
			fmt.Printf("      %d. %s\n", i+1, file)
		}
	}

	// Example 9: Real-world scenario
	fmt.Println("\n9. Real-world Scenario - Find project root")
	fmt.Println("   Looking for project root by finding 'go.mod' file...")

	result, err = findup.FindUp("go.mod", &findup.Options{Cwd: nestedDir})
	if err != nil {
		log.Printf("Error: %v", err)
	} else if result != "" {
		projectRoot := filepath.Dir(result)
		fmt.Printf("   ✅ Project root found: %s\n", projectRoot)
	} else {
		fmt.Println("   ℹ️  No go.mod found (not a Go project)")
	}

	fmt.Println("\n=== Examples Complete ===")
	fmt.Println("\nFor more information, visit: https://github.com/viguza/find-up")
}

func setupExampleStructure(baseDir string) {
	// Create a realistic project structure
	dirs := []string{
		filepath.Join(baseDir, "project"),
		filepath.Join(baseDir, "project", "src"),
		filepath.Join(baseDir, "project", "src", "components"),
		filepath.Join(baseDir, "project", "src", "components", "ui"),
		filepath.Join(baseDir, "project", "src", "utils"),
		filepath.Join(baseDir, "project", "docs"),
		filepath.Join(baseDir, "project", "node_modules"),
		filepath.Join(baseDir, "project", "node_modules", "some-package"),
	}

	for _, dir := range dirs {
		os.MkdirAll(dir, 0755)
	}

	// Create files
	files := map[string]string{
		// Project files
		filepath.Join(baseDir, "project", "package.json"): `{"name": "my-project", "version": "1.0.0"}`,
		filepath.Join(baseDir, "project", "go.mod"):       `module my-project`,
		filepath.Join(baseDir, "project", "README.md"):    `# My Project`,
		filepath.Join(baseDir, "project", "config.json"):  `{"debug": true}`,
		filepath.Join(baseDir, "project", "notes.txt"):    `Project notes`,

		// Source files
		filepath.Join(baseDir, "project", "src", "main.js"):     `console.log('Hello World');`,
		filepath.Join(baseDir, "project", "src", "utils.js"):    `export function helper() {}`,
		filepath.Join(baseDir, "project", "src", "config.json"): `{"src": true}`,
		filepath.Join(baseDir, "project", "src", "readme.txt"):  `Source readme`,

		// Component files
		filepath.Join(baseDir, "project", "src", "components", "Button.js"):         `export default Button;`,
		filepath.Join(baseDir, "project", "src", "components", "ui", "Modal.js"):    `export default Modal;`,
		filepath.Join(baseDir, "project", "src", "components", "ui", "config.json"): `{"ui": true}`,

		// Documentation
		filepath.Join(baseDir, "project", "docs", "api.md"):    `# API Documentation`,
		filepath.Join(baseDir, "project", "docs", "guide.txt"): `User guide`,

		// Node modules
		filepath.Join(baseDir, "project", "node_modules", "some-package", "package.json"): `{"name": "some-package"}`,
	}

	for file, content := range files {
		os.WriteFile(file, []byte(content), 0644)
	}
}
