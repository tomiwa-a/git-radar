package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tomiwa-a/git-radar/internal/types"
)

var (
	divBorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#444"))

	divSectionTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#BD93F9"))

	divDimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272A4"))

	divHashStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFB86C"))

	divMessageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F8F8F2"))

	divAuthorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8BE9FD"))

	divAddStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50FA7B"))

	divDelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5555"))

	divWarningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFB86C"))

	divSelectedStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#44475A"))

	divIncomingTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#8BE9FD"))

	divOutgoingTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#50FA7B"))
)

type DivergenceData struct {
	TargetBranch      string
	SourceBranch      string
	MergeBase         *types.GraphCommit
	Incoming          []types.GraphCommit
	Outgoing          []types.GraphCommit
	IncomingIdx       int
	OutgoingIdx       int
	ActivePane        int
	TotalFiles        int
	TotalAdditions    int
	TotalDeletions    int
	ConflictFiles     []string
	LoadingDivergence bool
	AlertMessage      string
}

func GetDummyDivergenceData() DivergenceData {
	return DivergenceData{
		TargetBranch: "main",
		SourceBranch: "feature/logging",
		MergeBase: &types.GraphCommit{
			Hash:    "abc1234",
			Message: "initial project setup",
			Author:  "Tomiwa",
			Date:    "5 days ago",
		},
		Incoming: []types.GraphCommit{
			{Hash: "def5678", Message: "fix auth bug", Author: "Tomiwa", Date: "2 days ago", Files: []types.FileChange{{Status: "M", Path: "auth.go", Additions: 10, Deletions: 5}}},
			{Hash: "jkl3456", Message: "Merge pull request #4", Author: "Tomiwa", Date: "3 days ago", IsMerge: true},
			{Hash: "pqr1234", Message: "update dependencies", Author: "Tomiwa", Date: "4 days ago"},
		},
		Outgoing: []types.GraphCommit{
			{Hash: "ghi9012", Message: "add logging to service", Author: "Tomiwa", Date: "2 hours ago", Files: []types.FileChange{{Status: "A", Path: "logger.go", Additions: 30, Deletions: 0}, {Status: "M", Path: "service.go", Additions: 15, Deletions: 2}}},
			{Hash: "mno7890", Message: "refactor git service", Author: "Tomiwa", Date: "5 hours ago", Files: []types.FileChange{{Status: "M", Path: "service.go", Additions: 50, Deletions: 30}}},
			{Hash: "stu5678", Message: "add unit tests", Author: "Tomiwa", Date: "1 day ago"},
			{Hash: "vwx9012", Message: "cleanup unused code", Author: "Tomiwa", Date: "1 day ago"},
			{Hash: "yza3456", Message: "update documentation", Author: "Tomiwa", Date: "2 days ago"},
		},
		IncomingIdx:    0,
		OutgoingIdx:    0,
		ActivePane:     1,
		TotalFiles:     8,
		TotalAdditions: 120,
		TotalDeletions: 45,
		ConflictFiles:  []string{"service.go", "main.go"},
	}
}

func RenderDivergence(width, height int, data DivergenceData) string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF79C6"))
	title := headerStyle.Render(" Git-Radar ")
	compareLabel := divDimStyle.Render("Compare: ") + divOutgoingTitleStyle.Render(data.TargetBranch) + divDimStyle.Render(" ← ") + divIncomingTitleStyle.Render(data.SourceBranch)
	headerGap := width - lipgloss.Width(title) - lipgloss.Width(compareLabel) - 4
	if headerGap < 0 {
		headerGap = 0
	}

	if data.LoadingDivergence {
		// Change title to alert message if one exists even during loading
		if data.AlertMessage != "" {
			alertStyle := lipgloss.NewStyle().Background(lipgloss.Color("#50FA7B")).Foreground(lipgloss.Color("#282A36")).Bold(true).Padding(0, 1)
			title = alertStyle.Render(" " + data.AlertMessage + " ")
		}
		b.WriteString(title + strings.Repeat(" ", headerGap) + compareLabel + "\n")
		b.WriteString(divDimStyle.Render(strings.Repeat("─", width-2)) + "\n\n")

		loadingStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB86C")).Bold(true)
		b.WriteString("\n\n" + loadingStyle.Render("  Loading divergence data...") + "\n\n")
		help := divDimStyle.Render("y: copy hash │ esc: back │ q: quit")
		b.WriteString(help)
		return b.String()
	}

	// Change title to alert message if one exists
	if data.AlertMessage != "" {
		alertStyle := lipgloss.NewStyle().Background(lipgloss.Color("#50FA7B")).Foreground(lipgloss.Color("#282A36")).Bold(true).Padding(0, 1)
		title = alertStyle.Render(" " + data.AlertMessage + " ")
	}

	b.WriteString(title + strings.Repeat(" ", headerGap) + compareLabel + "\n")
	b.WriteString(divDimStyle.Render(strings.Repeat("─", width-2)) + "\n\n")

	b.WriteString(renderDivergedAt(width, data.MergeBase))
	b.WriteString("\n")

	b.WriteString(renderCommitPanes(width, height, data))
	b.WriteString("\n")

	// Conditionally render details if there's space
	if height > 25 {
		b.WriteString(renderSelectedCommit(width, data))
		b.WriteString("\n")
	}

	if height > 15 {
		b.WriteString(renderTotalChanges(width, height, data))
		b.WriteString("\n")
	}

	help := divDimStyle.Render("←/→: pane │ ↑/↓: nav │ enter: diff │ y: copy │ b: src │ c: target │ esc: back │ q: quit")
	b.WriteString(help)

	return b.String()
}

func renderDivergedAt(width int, mergeBase *types.GraphCommit) string {
	if mergeBase == nil {
		return ""
	}

	var content strings.Builder
	content.WriteString(" " + divSectionTitleStyle.Render("DIVERGED AT") + "\n")
	content.WriteString(" ○ " + divHashStyle.Render(mergeBase.Hash) + "  ")
	content.WriteString(divDimStyle.Render("\""+mergeBase.Message+"\"") + "  ")
	content.WriteString(divAuthorStyle.Render(mergeBase.Author) + " · " + divDimStyle.Render(mergeBase.Date))

	return divBorderStyle.Width(width-4).Render(content.String()) + "\n"
}

func renderCommitPanes(width, height int, data DivergenceData) string {
	paneWidth := (width - 6) / 2
	if paneWidth < 30 {
		paneWidth = 30
	}

	// Calculate available height for panes
	// Diverged At (~4) + Selected Commit (~8-12) + Total Changes (~6) + Help (1) + Spacing/Headers (~10)
	// This is a rough estimate, let's try to give it as much as possible but cap it
	reservedHeight := 25
	availableHeight := height - reservedHeight
	if availableHeight < 5 {
		availableHeight = 5
	}
	if availableHeight > 10 {
		availableHeight = 10 // Let's keep a reasonable cap to avoid overwhelming
	}

	leftContent := renderCommitPane("⬇ INCOMING", data.Incoming, data.IncomingIdx, data.ActivePane == 0, paneWidth, availableHeight, true)
	rightContent := renderCommitPane("⬆ OUTGOING", data.Outgoing, data.OutgoingIdx, data.ActivePane == 1, paneWidth, availableHeight, false)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftContent, "  ", rightContent)
}

func renderCommitPane(title string, commits []types.GraphCommit, selectedIdx int, isActive bool, width, height int, isIncoming bool) string {
	var b strings.Builder

	var titleStyle lipgloss.Style
	if isIncoming {
		titleStyle = divIncomingTitleStyle
	} else {
		titleStyle = divOutgoingTitleStyle
	}

	b.WriteString(" " + titleStyle.Render(fmt.Sprintf("%s (%d commits)", title, len(commits))) + "\n\n")

	maxVisible := height - 4 // Title (2) + Borders (2)
	if maxVisible < 1 {
		maxVisible = 1
	}

	startIdx := 0
	if selectedIdx >= maxVisible {
		startIdx = selectedIdx - maxVisible + 1
	}
	endIdx := startIdx + maxVisible
	if endIdx > len(commits) {
		endIdx = len(commits)
	}

	if startIdx > 0 {
		b.WriteString(divDimStyle.Render(fmt.Sprintf("  ↑ %d more above", startIdx)) + "\n")
	}

	for i := startIdx; i < endIdx; i++ {
		commit := commits[i]
		dot := "○"
		if commit.IsMerge {
			dot = "◆"
		}
		if i == selectedIdx && isActive {
			dot = "●"
		}

		msg := commit.Message
		maxMsgLen := width - 15
		if maxMsgLen < 10 {
			maxMsgLen = 10
		}
		if len(msg) > maxMsgLen {
			msg = msg[:maxMsgLen-3] + "..."
		}

		line := fmt.Sprintf(" %s %s  %s", dot, divHashStyle.Render(commit.Hash), msg)

		if i == selectedIdx && isActive {
			line = divSelectedStyle.Render(line)
		}

		b.WriteString(line + "\n")
	}

	if endIdx < len(commits) {
		b.WriteString(divDimStyle.Render(fmt.Sprintf("  ↓ %d more below", len(commits)-endIdx)) + "\n")
	}

	linesRendered := endIdx - startIdx
	if startIdx > 0 {
		linesRendered++
	}
	if endIdx < len(commits) {
		linesRendered++
	}
	for i := linesRendered; i < maxVisible+2; i++ {
		b.WriteString("\n")
	}

	paneStyle := divBorderStyle.Width(width)
	if isActive {
		paneStyle = paneStyle.BorderForeground(lipgloss.Color("#BD93F9"))
	}

	return paneStyle.Render(b.String())
}

func renderSelectedCommit(width int, data DivergenceData) string {
	var commit types.GraphCommit
	if data.ActivePane == 0 && len(data.Incoming) > 0 {
		commit = data.Incoming[data.IncomingIdx]
	} else if data.ActivePane == 1 && len(data.Outgoing) > 0 {
		commit = data.Outgoing[data.OutgoingIdx]
	}

	if commit.Hash == "" {
		return ""
	}

	var b strings.Builder
	b.WriteString(" " + divSectionTitleStyle.Render("SELECTED COMMIT") + "\n\n")
	b.WriteString(" " + divHashStyle.Render(commit.Hash) + "  " + divMessageStyle.Render("\""+commit.Message+"\"") + "\n")
	b.WriteString(" " + divDimStyle.Render("Author:") + " " + divAuthorStyle.Render(commit.Author) + " · " + divDimStyle.Render(commit.Date) + "\n\n")

	if len(commit.Files) > 0 {
		totalAdds, totalDels := 0, 0
		for _, f := range commit.Files {
			totalAdds += f.Additions
			totalDels += f.Deletions
		}

		b.WriteString(" " + divDimStyle.Render(fmt.Sprintf("FILES (%d)", len(commit.Files))) + "           ")
		b.WriteString(divAddStyle.Render(fmt.Sprintf("+%d", totalAdds)) + "  " + divDelStyle.Render(fmt.Sprintf("-%d", totalDels)) + "\n")

		for _, file := range commit.Files {
			statusStyle := divDimStyle
			switch file.Status {
			case "A":
				statusStyle = divAddStyle
			case "M":
				statusStyle = divWarningStyle
			case "D":
				statusStyle = divDelStyle
			}
			b.WriteString(" " + statusStyle.Render(file.Status) + "  " + divDimStyle.Render(file.Path))
			b.WriteString("       " + divAddStyle.Render(fmt.Sprintf("+%d", file.Additions)) + "  " + divDelStyle.Render(fmt.Sprintf("-%d", file.Deletions)) + "\n")
		}
	} else {
		b.WriteString(" " + divDimStyle.Render("No file information available") + "\n")
	}

	return divBorderStyle.Width(width-4).Render(b.String()) + "\n"
}

func renderTotalChanges(width, height int, data DivergenceData) string {
	var b strings.Builder
	b.WriteString(" " + divSectionTitleStyle.Render(fmt.Sprintf("TOTAL CHANGES (%s → %s)", data.SourceBranch, data.TargetBranch)) + "\n\n")

	b.WriteString(" " + divDimStyle.Render(fmt.Sprintf("%d files changed", data.TotalFiles)) + "   ")
	b.WriteString(divAddStyle.Render(fmt.Sprintf("+%d", data.TotalAdditions)) + "  ")
	b.WriteString(divDelStyle.Render(fmt.Sprintf("-%d", data.TotalDeletions)) + "\n\n")

	if len(data.ConflictFiles) > 0 {
		b.WriteString(" " + divWarningStyle.Render(fmt.Sprintf("⚠ Potential conflicts: %d files modified in both branches", len(data.ConflictFiles))) + "\n")
		// Cap conflict files to 2 if height is tight, otherwise more
		maxConflicts := 3
		if height < 30 {
			maxConflicts = 1
		}
		for i, f := range data.ConflictFiles {
			if i >= maxConflicts {
				b.WriteString("   " + divDimStyle.Render(fmt.Sprintf("• and %d more...", len(data.ConflictFiles)-i)) + "\n")
				break
			}
			b.WriteString("   " + divDimStyle.Render("• "+f) + "\n")
		}
	} else {
		b.WriteString(" " + divAddStyle.Render("✓ No potential conflicts detected") + "\n")
	}

	return divBorderStyle.Width(width - 4).Render(b.String())
}
