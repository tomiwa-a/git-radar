package screens

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tomiwa-a/git-radar/internal/types"
	"github.com/tomiwa-a/git-radar/utils"
)

func RenderBranchModal(width, height int, branches []types.Branch, selectedIdx int, currentBranch string) string {
	modalWidth := 40

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#BD93F9")).
		Padding(1, 2)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#F8F8F2")).
		MarginBottom(1)

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6272A4")).
		MarginTop(1)

	selectedStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#44475A")).
		Foreground(lipgloss.Color("#F8F8F2")).
		Bold(true).
		Width(modalWidth - 6)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F8F8F2")).
		Width(modalWidth - 6)

	currentMarker := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#50FA7B")).
		Bold(true)

	remoteMarker := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8BE9FD")).
		Italic(true)

	sectionHeader := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6272A4")).
		Bold(true).
		MarginTop(1)

	var content strings.Builder

	title := titleStyle.Render("Switch Branch")
	escHint := utils.DetailsLabelStyle.Render("ESC")
	titleLine := title + strings.Repeat(" ", modalWidth-lipgloss.Width(title)-lipgloss.Width(escHint)-4) + escHint
	content.WriteString(titleLine + "\n\n")

	maxVisible := height - 14
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

	// Track section headers
	inLocalSection := false
	inRemoteSection := false

	for i := startIdx; i < endIdx; i++ {
		branch := branches[i]

		// Add section headers
		if !branch.IsRemote && !inLocalSection {
			if i > startIdx {
				content.WriteString("\n")
			}
			content.WriteString(sectionHeader.Render("Local") + "\n")
			inLocalSection = true
		} else if branch.IsRemote && !inRemoteSection {
			content.WriteString("\n" + sectionHeader.Render("Remote") + "\n")
			inRemoteSection = true
		}

		suffix := ""
		displayName := branch.Name

		if branch.IsHead {
			suffix = currentMarker.Render(" ✓")
		}

		if branch.IsRemote {
			displayName = remoteMarker.Render(branch.Name)
		}

		if i == selectedIdx {
			content.WriteString(selectedStyle.Render("→ "+displayName+suffix) + "\n")
		} else {
			content.WriteString(normalStyle.Render("  "+displayName+suffix) + "\n")
		}
	}

	hints := hintStyle.Render("↑/↓: select │ enter: switch │ esc: close")
	content.WriteString("\n" + hints)

	modal := borderStyle.Render(content.String())

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, modal,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("#282A36")))
}
