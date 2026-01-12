package utils

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/charmbracelet/lipgloss"
)

func HighlightCode(code string, filename string) string {
	lexer := lexers.Match(filename)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	style := styles.Get("dracula")
	if style == nil {
		style = styles.Fallback
	}

	formatter := formatters.Get("terminal256")
	if formatter == nil {
		formatter = formatters.Fallback
	}

	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return code
	}

	var buf bytes.Buffer
	err = formatter.Format(&buf, style, iterator)
	if err != nil {
		return code
	}

	return buf.String()
}

func RenderCodeWithLineNumbers(code string, filename string, width int) string {
	highlighted := HighlightCode(code, filename)
	lines := strings.Split(highlighted, "\n")

	lineNumStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6272A4")).
		Width(4).
		Align(lipgloss.Right).
		MarginRight(1)

	dividerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#44475A"))

	var result strings.Builder
	for i, line := range lines {
		lineNum := lineNumStyle.Render(fmt.Sprintf("%d", i+1))
		divider := dividerStyle.Render("│")
		result.WriteString(lineNum + divider + " " + line)
		if i < len(lines)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}
