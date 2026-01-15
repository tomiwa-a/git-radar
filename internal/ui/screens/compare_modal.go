package screens

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tomiwa-a/git-radar/internal/types"
)

func RenderCompareModal(width, height int, localView, remoteView string, filterValue string, activePane int) string {
	modalWidth := int(float64(width) * 0.8)
	modalHeight := int(float64(height) * 0.7)

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#BD93F9")).
		Padding(0, 1)

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#F8F8F2")).
		Padding(0, 1)

	localHeader := headerStyle.Render("LOCAL")
	if activePane == 0 {
		localHeader = headerStyle.Copy().
			Background(lipgloss.Color("#BD93F9")).
			Foreground(lipgloss.Color("#282A36")).
			Render(" LOCAL ")
	}

	remoteHeader := headerStyle.Render("REMOTE")
	if activePane == 1 {
		remoteHeader = headerStyle.Copy().
			Background(lipgloss.Color("#BD93F9")).
			Foreground(lipgloss.Color("#282A36")).
			Render(" REMOTE ")
	}

	// Filter bar
	filterStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6272A4")).
		Padding(0, 1)
	filterBar := filterStyle.Render("Filter: ") + lipgloss.NewStyle().Foreground(lipgloss.Color("#F8F8F2")).Render(filterValue)

	divider := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#44475A")).
		Render("│")

	paneWidth := (modalWidth - 3) / 2

	titles := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().Width(paneWidth).Align(lipgloss.Center).Render(localHeader),
		" ",
		lipgloss.NewStyle().Width(paneWidth).Align(lipgloss.Center).Render(remoteHeader),
	)

	body := lipgloss.JoinHorizontal(
		lipgloss.Top,
		localView,
		" "+divider+" ",
		remoteView,
	)

	footer := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6272A4")).
		MarginTop(1).
		Render("tab: switch pane • ↑/↓: navigate • enter: compare • esc: close")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		titles,
		"\n",
		body,
		"\n",
		filterBar,
		footer,
	)

	modal := borderStyle.Width(modalWidth).Height(modalHeight).Render(content)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, modal,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("#282A36")))
}

func RenderBranchListContent(width int, branches []types.Branch, selectedIdx int, isActive bool) string {
	var b strings.Builder

	selectedStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#44475A")).
		Foreground(lipgloss.Color("#50FA7B")).
		Bold(true).
		Width(width)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F8F8F2")).
		Width(width)

	dimmedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6272A4")).
		Width(width)

	if len(branches) == 0 {
		return "\n  No branches found"
	}

	for i, branch := range branches {
		if i == selectedIdx && isActive {
			b.WriteString(selectedStyle.Render("→ "+branch.Name) + "\n")
		} else if i == selectedIdx && !isActive {
			b.WriteString(lipgloss.NewStyle().
				Background(lipgloss.Color("#282A36")).
				Foreground(lipgloss.Color("#BD93F9")).
				Render("  "+branch.Name) + "\n")
		} else if !isActive {
			b.WriteString(dimmedStyle.Render("  "+branch.Name) + "\n")
		} else {
			b.WriteString(normalStyle.Render("  "+branch.Name) + "\n")
		}
	}

	return b.String()
}
