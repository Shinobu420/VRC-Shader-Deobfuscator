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

	defineRegex := regexp.MustCompile(`(?m)^\s*#define\s+([a-zA-Z0-9_]+)\s+(.*?)(?:\r?\n|$)`)
	matches := defineRegex.FindAllStringSubmatch(text, -1)

	defines := make(map[string]string)
	for _, match := range matches {
		defines[match[1]] = strings.TrimSpace(match[2])
	}

	cleanText := defineRegex.ReplaceAllString(text, "")
	resolvedDefines := resolveDefines(defines)

	identRegex := regexp.MustCompile(`\b[a-zA-Z_][a-zA-Z0-9_]*\b`)
	deobfuscatedText := identRegex.ReplaceAllStringFunc(cleanText, func(ident string) string {
		if val, ok := resolvedDefines[ident]; ok {
			return val
		}
		return ident
	})

	multiNewlines := regexp.MustCompile(`(?:\r?\n){3,}`)
	deobfuscatedText = multiNewlines.ReplaceAllString(deobfuscatedText, "\n\n")

	deobfuscatedText = strings.TrimSpace(deobfuscatedText)
	if deobfuscatedText != "" {
		deobfuscatedText += "\n"
	}

	deobfuscatedText = strings.ReplaceAll(deobfuscatedText, ";", ";\n")
	deobfuscatedText = strings.ReplaceAll(deobfuscatedText, "{", "{\n")
	deobfuscatedText = strings.ReplaceAll(deobfuscatedText, "}", "}\n")

	return []byte(deobfuscatedText), nil
}

func resolveDefines(defines map[string]string) map[string]string {
	resolved := make(map[string]string, len(defines))
	for k, v := range defines {
		resolved[k] = expand(v, defines, map[string]bool{k: true})
	}
	return resolved
}

func expand(val string, defines map[string]string, visited map[string]bool) string {
	identRegex := regexp.MustCompile(`\b[a-zA-Z_][a-zA-Z0-9_]*\b`)
	return identRegex.ReplaceAllStringFunc(val, func(ident string) string {
		if nextVal, ok := defines[ident]; ok {
			if visited[ident] {
				return ident
			}
			newVisited := make(map[string]bool, len(visited)+1)
			for vk, vv := range visited {
				newVisited[vk] = vv
			}
			newVisited[ident] = true
			return expand(nextVal, defines, newVisited)
		}
		return ident
	})
}
