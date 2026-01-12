package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tomiwa-a/git-radar/internal/types"
	"github.com/tomiwa-a/git-radar/utils"
)

func RenderDiffs(width int, commit types.Commit, fileIdx int, viewportContent string) string {
	var b strings.Builder

	file := commit.Files[fileIdx]

	backHint := utils.DetailsLabelStyle.Render("ESC: back  h/l: switch files  ↑↓: scroll")
	fileName := utils.FileNameStyle.Render(file.Path)
	stats := renderFileStats(file)

	headerLine := lipgloss.JoinHorizontal(
		lipgloss.Center,
		fileName,
		"  ",
		stats,
	)

	headerGap := width - lipgloss.Width(headerLine) - lipgloss.Width(backHint)
	if headerGap < 0 {
		headerGap = 0
	}
	header := headerLine + strings.Repeat(" ", headerGap) + backHint
	b.WriteString(header + "\n")

	commitInfo := utils.HashStyle.Render(commit.Hash) +
		" " +
		utils.DetailsTitleStyle.Render(commit.Message) +
		" " +
		utils.DetailsLabelStyle.Render("by "+commit.Author)
	b.WriteString(commitInfo + "\n")

	divider := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#44475A")).
		Render(strings.Repeat("─", width))
	b.WriteString(divider + "\n")

	b.WriteString(viewportContent)

	return b.String()
}

func renderFileStats(file types.FileChange) string {
	addStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B")).Bold(true)
	delStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555")).Bold(true)

	stats := addStyle.Render(fmt.Sprintf("+%d", file.Additions)) +
		" " +
		delStyle.Render(fmt.Sprintf("-%d", file.Deletions))

	return stats
}
