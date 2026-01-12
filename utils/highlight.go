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
	"github.com/tomiwa-a/git-radar/internal/types"
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

func RenderDiffLines(diffLines []types.DiffLine, filename string) string {
	addStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#50FA7B"))

	delStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF5555"))

	lineNumStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6272A4")).
		Width(4).
		Align(lipgloss.Right)

	dividerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#44475A"))

	var equalLines []int
	var equalCode strings.Builder
	for i, dl := range diffLines {
		if dl.Type == "equal" {
			equalLines = append(equalLines, i)
			equalCode.WriteString(dl.Content + "\n")
		}
	}

	highlightedEqual := ""
	if equalCode.Len() > 0 {
		highlightedEqual = HighlightCode(equalCode.String(), filename)
	}
	highlightedEqualLines := strings.Split(highlightedEqual, "\n")

	var result strings.Builder
	lineNum := 1
	equalIdx := 0

	for i, dl := range diffLines {
		var prefix string
		var numStr string
		var codeLine string

		switch dl.Type {
		case "add":
			prefix = addStyle.Render("+")
			numStr = fmt.Sprintf("%d", lineNum)
			codeLine = addStyle.Render(dl.Content)
			lineNum++
		case "del":
			prefix = delStyle.Render("-")
			numStr = " "
			codeLine = delStyle.Render(dl.Content)
		default:
			prefix = " "
			numStr = fmt.Sprintf("%d", lineNum)
			if equalIdx < len(highlightedEqualLines) {
				codeLine = highlightedEqualLines[equalIdx]
				equalIdx++
			} else {
				codeLine = dl.Content
			}
			lineNum++
		}

		num := lineNumStyle.Render(numStr)
		divider := dividerStyle.Render("│")
		result.WriteString(num + divider + prefix + " " + codeLine)

		if i < len(diffLines)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}
