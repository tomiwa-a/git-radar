package screens

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tomiwa-a/git-radar/utils"
)

func RenderCompareModal(width, height int, branches []string, selectedIdx int, currentBranch string) string {
	modalWidth := 44

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#50FA7B")).
		Padding(1, 2)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#F8F8F2")).
		MarginBottom(1)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6272A4")).
		Italic(true).
		MarginBottom(1)

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6272A4")).
		MarginTop(1)

	selectedStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#44475A")).
		Foreground(lipgloss.Color("#50FA7B")).
		Bold(true).
		Width(modalWidth - 6)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F8F8F2")).
		Width(modalWidth - 6)

	var content strings.Builder

	title := titleStyle.Render("Compare Branches")
	escHint := utils.DetailsLabelStyle.Render("ESC")
	titleLine := title + strings.Repeat(" ", modalWidth-lipgloss.Width(title)-lipgloss.Width(escHint)-4) + escHint
	content.WriteString(titleLine + "\n")

	subtitle := subtitleStyle.Render("Compare " + currentBranch + " against:")
	content.WriteString(subtitle + "\n\n")

	if len(branches) == 0 {
		content.WriteString(normalStyle.Render("  No other branches available") + "\n")
	} else {
		maxVisible := height - 12
		if maxVisible < 3 {
			maxVisible = 3
		}

		startIdx := 0
		if len(branches) > maxVisible && selectedIdx >= maxVisible {
			startIdx = selectedIdx - maxVisible + 1
		}
		endIdx := startIdx + maxVisible
		if endIdx > len(branches) {
			endIdx = len(branches)
		}

		for i := startIdx; i < endIdx; i++ {
			branch := branches[i]
			if i == selectedIdx {
				content.WriteString(selectedStyle.Render("→ "+branch) + "\n")
			} else {
				content.WriteString(normalStyle.Render("  "+branch) + "\n")
			}
		}
	}

	hints := hintStyle.Render("↑/↓: select │ enter: compare │ esc: close")
	content.WriteString("\n" + hints)

	modal := borderStyle.Render(content.String())

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, modal,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("#282A36")))
}
