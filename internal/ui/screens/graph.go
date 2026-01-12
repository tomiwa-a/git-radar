package screens

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tomiwa-a/git-radar/internal/types"
	"github.com/tomiwa-a/git-radar/utils"
)

func RenderGraph(width int, commits []types.GraphCommit, selectedIdx int, currentBranch string) string {
	var b strings.Builder

	title := utils.TitleStyle.Render("Git-Radar")
	branchLabel := utils.DetailsLabelStyle.Render("branch: ")
	branchName := utils.BranchStyle.Render(currentBranch)
	hints := utils.DetailsLabelStyle.Render("c: compare  q: quit")

	headerGap := width - lipgloss.Width(title) - lipgloss.Width(branchLabel) - lipgloss.Width(branchName) - lipgloss.Width(hints)
	if headerGap < 0 {
		headerGap = 0
	}
	header := title + strings.Repeat(" ", headerGap/2) + branchLabel + branchName + strings.Repeat(" ", headerGap/2) + hints
	b.WriteString(header + "\n")

	divider := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#44475A")).
		Render(strings.Repeat("─", width))
	b.WriteString(divider + "\n\n")

	graphStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF79C6"))
	hashStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB86C"))
	messageStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#F8F8F2"))
	branchTagStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B")).Bold(true)
	headStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#8BE9FD")).Bold(true)
	selectedStyle := lipgloss.NewStyle().Background(lipgloss.Color("#44475A"))
	authorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#6272A4"))

	for i, commit := range commits {
		graph := graphStyle.Render(commit.GraphChars)
		hash := hashStyle.Render(commit.Hash[:7])

		var branchTags string
		if len(commit.Branches) > 0 {
			var tags []string
			for _, branch := range commit.Branches {
				if branch == "HEAD" {
					tags = append(tags, headStyle.Render("HEAD"))
				} else {
					tags = append(tags, branchTagStyle.Render(branch))
				}
			}
			branchTags = " (" + strings.Join(tags, ", ") + ")"
		}

		msg := messageStyle.Render(truncateMessage(commit.Message, 50))
		author := authorStyle.Render(" <" + commit.Author + ">")

		line := graph + hash + branchTags + " " + msg + author

		if i == selectedIdx {
			line = selectedStyle.Render(line + strings.Repeat(" ", max(0, width-lipgloss.Width(line))))
		}

		b.WriteString(line + "\n")
	}

	b.WriteString("\n")
	footer := utils.DetailsLabelStyle.Render("↑/↓: navigate │ enter: view commit │ c: compare branches │ q: quit")
	b.WriteString(footer)

	return b.String()
}

func truncateMessage(msg string, maxLen int) string {
	if len(msg) <= maxLen {
		return msg
	}
	return msg[:maxLen-3] + "..."
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
