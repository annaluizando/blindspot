package utils

import (
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

func HighlightCode(code string, language string) string {
	code = strings.TrimSpace(code)

	var lexer chroma.Lexer
	if language == "" {
		lexer = lexers.Analyse(code)
		if lexer == nil {
			lexer = lexers.Get("go")
		}
	} else {
		lexer = lexers.Get(language)
		if lexer == nil {
			// If specified language not found, try to detect
			lexer = lexers.Analyse(code)
			if lexer == nil {
				return code
			}
		}
	}

	style := styles.Get("monokai")
	if style == nil {
		style = styles.Fallback
	}

	// Use terminal formatter
	formatter := formatters.Get("terminal")
	if formatter == nil {
		formatter = formatters.Fallback
	}

	// Tokenize the code
	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return code
	}

	// Create a buffer to hold the highlighted code
	var buf strings.Builder
	err = formatter.Format(&buf, style, iterator)
	if err != nil {
		return code
	}

	return buf.String()
}

func GetLanguageFromExtension(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) < 2 {
		return "" // No extension
	}

	ext := parts[len(parts)-1]
	switch strings.ToLower(ext) {
	case "go":
		return "go"
	case "py", "python":
		return "python"
	case "js", "javascript":
		return "javascript"
	case "php":
		return "php"
	case "java":
		return "java"
	case "c", "cpp", "cc", "cxx":
		return "c++"
	case "cs":
		return "c#"
	case "rb":
		return "ruby"
	case "rs":
		return "rust"
	case "ts":
		return "typescript"
	case "sh", "bash":
		return "bash"
	default:
		return "" // Unknown extension
	}
}

// To be used in code fix
// takes a code string and highlights a specific line or pattern
// that contains the vulnerability to make it more obvious
func HighlightVulnerability(code string, vulnerablePattern string) string {
	if vulnerablePattern == "" {
		return code // Return unchanged if no pattern provided
	}

	// Simple implementation: just highlight the vulnerable part with terminal escapes
	// This assumes we're using a terminal that supports ANSI color codes
	const redBackground = "\x1b[41m" // Red background
	const resetColor = "\x1b[0m"     // Reset colors

	// Replace the pattern with a highlighted version
	highlighted := strings.Replace(
		code,
		vulnerablePattern,
		redBackground+vulnerablePattern+resetColor,
		-1,
	)

	return highlighted
}
