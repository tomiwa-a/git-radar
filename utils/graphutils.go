package utils

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// CompactTime converts relative time strings to a compact form (e.g., "19 minutes ago" -> "19m").
func CompactTime(timeStr string) string {
	timeStr = strings.TrimSpace(timeStr)

	if timeStr == "just now" {
		return "now"
	}

	parts := strings.Split(timeStr, " ")
	if len(parts) >= 2 {
		num := parts[0]
		unit := parts[1]

		switch {
		case strings.HasPrefix(unit, "minute"):
			return num + "m"
		case strings.HasPrefix(unit, "hour"):
			return num + "h"
		case strings.HasPrefix(unit, "day"), unit == "yesterday":
			if timeStr == "yesterday" {
				return "1d"
			}
			return num + "d"
		case strings.HasPrefix(unit, "week"):
			return num + "w"
		case strings.HasPrefix(unit, "month"):
			return num + "mo"
		case strings.HasPrefix(unit, "year"):
			return num + "y"
		}
	}

	return timeStr
}

// WrapText wraps a string to the given width, returning a slice of lines.
func WrapText(text string, width int) []string {
	if width <= 0 {
		return []string{text}
	}

	var lines []string
	words := strings.Fields(text)
	currentLine := ""

	for _, word := range words {
		if currentLine == "" {
			currentLine = word
		} else if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

// TruncateMessage truncates a string to maxLen, adding ... if needed.
func TruncateMessage(msg string, maxLen int) string {
	if len(msg) <= maxLen {
		return msg
	}
	return msg[:maxLen-3] + "..."
}

// Max returns the maximum of two ints.
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// RenderLegend renders the legend modal for the graph screen.
func RenderLegend(width, height int, selectedDotStyle, mergeDotStyle, branchCountStyle, mainBranchStyle, localBranchStyle, remoteBranchStyle lipgloss.Style, utilsDetailsLabelStyle lipgloss.Style) string {
	modalWidth := 50

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#BD93F9")).
		Padding(1, 2)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#F8F8F2")).
		MarginBottom(1)

	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#BD93F9")).
		Bold(true).
		MarginTop(1)

	itemStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F8F8F2"))

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6272A4"))

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6272A4")).
		MarginTop(1)

	var content strings.Builder

	// Title
	title := titleStyle.Render("LEGEND")
	escHint := utilsDetailsLabelStyle.Render("ESC/?")
	titleLine := title + strings.Repeat(" ", modalWidth-lipgloss.Width(title)-lipgloss.Width(escHint)-6) + escHint
	content.WriteString(titleLine + "\n")

	// Commits section
	content.WriteString("\n" + sectionStyle.Render("COMMITS") + "\n")
	content.WriteString(itemStyle.Render("  ○") + descStyle.Render("  Regular commit") + "\n")
	content.WriteString(selectedDotStyle.Render("  ●") + descStyle.Render("  Selected commit") + "\n")
	content.WriteString(mergeDotStyle.Render("  ◆") + descStyle.Render("  Merge commit") + "\n")

	// Indicators section
	content.WriteString("\n" + sectionStyle.Render("INDICATORS") + "\n")
	content.WriteString(branchCountStyle.Render("  ⚑2") + descStyle.Render("     2 branches at this commit") + "\n")
	content.WriteString(mainBranchStyle.Render("  ★") + descStyle.Render("      main/master branch") + "\n")

	// Branches section
	content.WriteString("\n" + sectionStyle.Render("BRANCHES") + "\n")
	content.WriteString(localBranchStyle.Render("  feature/xyz") + descStyle.Render("   Local (yours)") + "\n")
	content.WriteString(remoteBranchStyle.Render("  origin/xyz") + descStyle.Render("    Remote") + "\n")

	// Keys section
	content.WriteString("\n" + sectionStyle.Render("KEYS") + "\n")
	content.WriteString(itemStyle.Render("  ↑/↓ j/k") + descStyle.Render("   Navigate commits") + "\n")
	content.WriteString(itemStyle.Render("  enter") + descStyle.Render("     View commit files") + "\n")
	content.WriteString(itemStyle.Render("  b") + descStyle.Render("         Switch branch") + "\n")
	content.WriteString(itemStyle.Render("  c") + descStyle.Render("         Compare branches") + "\n")
	content.WriteString(itemStyle.Render("  q") + descStyle.Render("         Quit") + "\n")

	content.WriteString("\n" + hintStyle.Render("Press ESC or ? to close"))

	modal := borderStyle.Render(content.String())

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, modal,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("#282A36")))
}
