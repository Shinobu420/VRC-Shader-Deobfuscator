package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

/*
	Tool to deobfuscate unity HLSL shaders that use #define to obfuscate variable names.
*/
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run deobfuscator.go <target.shader>")
		return
	}

	inputFile := os.Args[1]
	bytes, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	text := string(bytes)

	// Regex to match: "#define neon_zL10649"
	re := regexp.MustCompile(`(?m)^\s*#define\s+([a-zA-Z0-9_]+)\s+(.*?)\r?$`)
	matches := re.FindAllStringSubmatch(text, -1)

	for _, match := range matches {
		fullMatch := match[0]
		obfuscatedName := match[1]
		realChar := strings.TrimSpace(match[2])

		// 1. Remove the #define line itself to clean up the code
		text = strings.Replace(text, fullMatch+"\n", "", 1)
		text = strings.Replace(text, fullMatch, "", 1) // Fallback for last line

		// 2. Replace all remaining instances of the obfuscated variable
		text = strings.ReplaceAll(text, obfuscatedName, realChar)
	}

	// Format output filename (e.g., file.shader -> file_deob.shader)
	ext := filepath.Ext(inputFile)
	base := strings.TrimSuffix(inputFile, ext)
	outputFile := base + "_deob" + ext

	// Clean up the dead space left by the removed #defines
	cleanSpacing := regexp.MustCompile(`[ \t]{2,}`)
	text = cleanSpacing.ReplaceAllString(text, "\n")

	// (Optional) Ensure newlines after semicolons and braces
	text = strings.ReplaceAll(text, ";", ";\n")
	text = strings.ReplaceAll(text, "{", "{\n")
	text = strings.ReplaceAll(text, "}", "}\n")

	// Write the cleaned data
	err = os.WriteFile(outputFile, []byte(text), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Printf("Success! Saved to: %s\n", outputFile)
}
