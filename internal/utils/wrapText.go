package utils

import "strings"

func WrapText(text string, width int) string {
	if width <= 0 {
		return text
	}

	width = width - 4 // leaving some margin

	var result strings.Builder
	lines := strings.Split(text, "\n")

	for i, line := range lines {
		if len(line) <= width { // line fits
			result.WriteString(line)
		} else { // line need wrapping
			words := strings.Fields(line)
			if len(words) == 0 {
				continue
			}

			lineLength := 0
			for _, word := range words {
				// Check if adding this word would exceed the width
				if lineLength+len(word)+1 > width && lineLength > 0 {
					// Start a new line
					result.WriteString("\n")
					lineLength = 0
				}

				if lineLength > 0 {
					// Add a space before the word
					result.WriteString(" ")
					lineLength++
				}

				// Add the word
				result.WriteString(word)
				lineLength += len(word)
			}
		}

		// Add newline if not the last line
		if i < len(lines)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}
