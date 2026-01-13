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
	const collapseThreshold = 8
	const contextLines = 3

	collapsed := collapseDiffLines(diffLines, collapseThreshold, contextLines)

	addStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#50FA7B"))

	delStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF5555"))

	collapseStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6272A4")).
		Italic(true)

	lineNumStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6272A4")).
		Width(4).
		Align(lipgloss.Right)

	dividerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#44475A"))

	var equalCode strings.Builder
	for _, dl := range collapsed {
		if dl.Type == "equal" {
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

	for i, dl := range collapsed {
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
		case "collapse":
			prefix = " "
			numStr = " "
			codeLine = collapseStyle.Render(dl.Content)
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

		if i < len(collapsed)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

func collapseDiffLines(lines []types.DiffLine, threshold, context int) []types.DiffLine {
	var result []types.DiffLine
	i := 0

	for i < len(lines) {
		if lines[i].Type != "equal" {
			result = append(result, lines[i])
			i++
			continue
		}

		runStart := i
		for i < len(lines) && lines[i].Type == "equal" {
			i++
		}
		runEnd := i
		runLen := runEnd - runStart

		if runLen > threshold {
			for j := runStart; j < runStart+context && j < runEnd; j++ {
				result = append(result, lines[j])
			}

			hidden := runLen - (context * 2)
			if hidden > 0 {
				result = append(result, types.DiffLine{
					Type:    "collapse",
					Content: fmt.Sprintf("⋯ %d lines hidden ⋯", hidden),
				})
			}

			for j := runEnd - context; j < runEnd; j++ {
				if j >= runStart+context {
					result = append(result, lines[j])
				}
			}
		} else {
			for j := runStart; j < runEnd; j++ {
				result = append(result, lines[j])
			}
		}
	}

	return result
}
