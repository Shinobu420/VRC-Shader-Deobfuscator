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

	deobfuscatedText = formatShader(deobfuscatedText)

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

func splitComment(line string) (code, comment string) {
	inString := false
	for i := 0; i < len(line)-1; i++ {
		if line[i] == '"' && (i == 0 || line[i-1] != '\\') {
			inString = !inString
		}
		if !inString && line[i] == '/' && line[i+1] == '/' {
			return line[:i], line[i:]
		}
	}
	return line, ""
}

func formatShader(text string) string {
	lines := strings.Split(text, "\n")
	inBlockComment := false
	var result []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if inBlockComment {
			if strings.Contains(trimmed, "*/") {
				inBlockComment = false
			}
			result = append(result, line)
			continue
		}
		if strings.Contains(trimmed, "/*") {
			if !strings.Contains(trimmed, "*/") {
				inBlockComment = true
			}
			result = append(result, line)
			continue
		}

		if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "#") {
			result = append(result, line)
			continue
		}

		code, comment := splitComment(line)
		if strings.TrimSpace(code) == "" {
			result = append(result, line)
			continue
		}

		isForLoop := strings.Contains(code, "for") && (strings.Contains(code, "(") || strings.Contains(code, " "))

		var formattedCode string
		if isForLoop {
			formattedCode = code
		} else {
			var sb strings.Builder
			for i := 0; i < len(code); i++ {
				ch := code[i]
				sb.WriteByte(ch)

				hasMoreCode := func(start int) bool {
					for j := start; j < len(code); j++ {
						c := code[j]
						if c != ' ' && c != '\t' && c != '\r' && c != '\n' {
							return true
						}
					}
					return false
				}

				if ch == ';' {
					if hasMoreCode(i + 1) {
						sb.WriteByte('\n')
					}
				} else if ch == '{' {
					if hasMoreCode(i + 1) {
						sb.WriteByte('\n')
					}
				} else if ch == '}' {
					if hasMoreCode(i + 1) {
						sb.WriteByte('\n')
					}
				}
			}
			formattedCode = sb.String()
		}

		result = append(result, formattedCode+comment)
	}

	return strings.Join(result, "\n")
}
