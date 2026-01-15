package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tomiwa-a/git-radar/internal/types"
	"github.com/tomiwa-a/git-radar/utils"
)

func RenderDiffs(width int, commit types.GraphCommit, files []types.FileChange, fileIdx int, viewportContent string, showFilter bool) string {
	var b strings.Builder

	file := files[fileIdx]

	backHint := utils.DetailsLabelStyle.Render("ESC: back  h/l: switch files  ↑↓: scroll")

	filterIndicator := ""
	if showFilter {
		filterIndicator = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF79C6")).
			Bold(true).
			Render("[FILTERED] ")
	}

	indexIndicator := utils.DetailsLabelStyle.Render(fmt.Sprintf("%d of %d", fileIdx+1, len(files)))
	fileName := utils.FileNameStyle.Render(file.Path)
	stats := renderFileStats(file)

	headerLine := lipgloss.JoinHorizontal(
		lipgloss.Center,
		filterIndicator,
		fileName,
		"  ",
		stats,
		"  ",
		indexIndicator,
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
