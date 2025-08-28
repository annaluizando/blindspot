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

	formatter := formatters.Get("terminal")
	if formatter == nil {
		formatter = formatters.Fallback
	}

	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return code
	}

	var buf strings.Builder
	err = formatter.Format(&buf, style, iterator)
	if err != nil {
		return code
	}

	return buf.String()
}

func GetLanguageFromChallenge(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) < 2 {
		return ""
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
		return ""
	}
}
