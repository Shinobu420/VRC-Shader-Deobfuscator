package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	outputFileFlag := flag.String("o", "", "Output path (defaults to <input>_deob.shader, use '-' for stdout)")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: go run deobfuscator.go [flags] <target.shader>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	inputFile := args[0]
	inputBytes, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	outputBytes, err := Deobfuscate(inputBytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if *outputFileFlag == "-" {
		fmt.Print(string(outputBytes))
		return
	}

	outputFile := *outputFileFlag
	if outputFile == "" {
		ext := filepath.Ext(inputFile)
		base := strings.TrimSuffix(inputFile, ext)
		outputFile = base + "_deob" + ext
	}

	err = os.WriteFile(outputFile, outputBytes, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Success! Saved to: %s\n", outputFile)
}

func Deobfuscate(input []byte) ([]byte, error) {
	text := string(input)

	re := regexp.MustCompile(`(?m)^\s*#define\s+([a-zA-Z0-9_]+)\s+(.*?)\r?$`)
	matches := re.FindAllStringSubmatch(text, -1)

	for _, match := range matches {
		fullMatch := match[0]
		obfuscatedName := match[1]
		realChar := strings.TrimSpace(match[2])

		text = strings.Replace(text, fullMatch+"\n", "", 1)
		text = strings.Replace(text, fullMatch, "", 1)
		text = strings.ReplaceAll(text, obfuscatedName, realChar)
	}

	cleanSpacing := regexp.MustCompile(`[ \t]{2,}`)
	text = cleanSpacing.ReplaceAllString(text, "\n")

	// Formatting is always applied
	text = strings.ReplaceAll(text, ";", ";\n")
	text = strings.ReplaceAll(text, "{", "{\n")
	text = strings.ReplaceAll(text, "}", "}\n")

	return []byte(text), nil
}
