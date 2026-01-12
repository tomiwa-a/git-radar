package ui

import (
	"fmt"
	"strings"

	"github.com/tomiwa-a/git-radar/internal/types"
	"github.com/tomiwa-a/git-radar/utils"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	switch m.Screen {
	case FileListScreen:
		return m.renderFileListScreen()
	default:
		return m.renderDashboardScreen()
	}
}

func (m Model) renderDashboardScreen() string {
	var b strings.Builder

	paneWidth := (m.Width - 4) / 2
	if paneWidth < 30 {
		paneWidth = 30
	}

	title := utils.TitleStyle.Render(" Git-Radar ")
	branches := utils.BranchStyle.Render(fmt.Sprintf(" %s ← %s ", m.TargetBranch, m.SourceBranch))
	headerGap := m.Width - lipgloss.Width(title) - lipgloss.Width(branches)
	if headerGap < 0 {
		headerGap = 0
	}
	header := title + strings.Repeat(" ", headerGap) + branches
	b.WriteString(header + "\n\n")

	incomingContent := m.renderCommitList(m.Incoming, m.IncomingIdx, m.ActivePane == IncomingPane)
	outgoingContent := m.renderCommitList(m.Outgoing, m.OutgoingIdx, m.ActivePane == OutgoingPane)

	incomingTitle := utils.PaneTitleIncoming.Render(fmt.Sprintf("⬇ INCOMING (Behind by %d)", len(m.Incoming)))
	outgoingTitle := utils.PaneTitleOutgoing.Render(fmt.Sprintf("⬆ OUTGOING (Ahead by %d)", len(m.Outgoing)))

	var leftPane, rightPane string
	if m.ActivePane == IncomingPane {
		leftPane = utils.ActivePaneStyle.Width(paneWidth).Render(incomingTitle + "\n\n" + incomingContent)
		rightPane = utils.PaneStyle.Width(paneWidth).Render(outgoingTitle + "\n\n" + outgoingContent)
	} else {
		leftPane = utils.PaneStyle.Width(paneWidth).Render(incomingTitle + "\n\n" + incomingContent)
		rightPane = utils.ActivePaneStyle.Width(paneWidth).Render(outgoingTitle + "\n\n" + outgoingContent)
	}

	panes := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)
	b.WriteString(panes + "\n")

	var selectedCommit types.Commit
	if m.ActivePane == IncomingPane && len(m.Incoming) > 0 {
		selectedCommit = m.Incoming[m.IncomingIdx]
	} else if m.ActivePane == OutgoingPane && len(m.Outgoing) > 0 {
		selectedCommit = m.Outgoing[m.OutgoingIdx]
	}

	detailsContent := m.renderDetails(selectedCommit)
	details := utils.DetailsStyle.Width(m.Width - 4).Render(detailsContent)
	b.WriteString(details + "\n")

	help := utils.HelpStyle.Render("←/→: switch pane │ ↑/↓: navigate │ enter: view files │ q: quit")
	b.WriteString(help)

	return b.String()
}

func (m Model) renderFileListScreen() string {
	var b strings.Builder

	// Header
	backHint := utils.DetailsLabelStyle.Render("ESC: back")
	commitInfo := utils.HashStyle.Render("← " + m.SelectedCommit.Hash + " ")
	commitMsg := utils.DetailsTitleStyle.Render(m.SelectedCommit.Message)
	headerGap := m.Width - lipgloss.Width(backHint) - lipgloss.Width(commitInfo) - lipgloss.Width(commitMsg)
	if headerGap < 0 {
		headerGap = 0
	}
	header := commitInfo + commitMsg + strings.Repeat(" ", headerGap) + backHint
	b.WriteString(header + "\n\n")

	// File list title
	fileCount := len(m.SelectedCommit.Files)
	listTitle := utils.DetailsTitleStyle.Render(fmt.Sprintf("FILES CHANGED (%d)", fileCount))
	b.WriteString(listTitle + "\n\n")

	// File list
	fileListContent := m.renderFileList()
	fileList := utils.PaneStyle.Width(m.Width - 4).Render(fileListContent)
	b.WriteString(fileList + "\n")

	// Help
	help := utils.HelpStyle.Render("↑/↓: navigate │ enter: view diff │ ESC: back │ q: quit")
	b.WriteString(help)

	return b.String()
}

func (m Model) renderFileList() string {
	var lines []string

	for i, file := range m.SelectedCommit.Files {
		cursor := "  "
		if i == m.FileIdx {
			cursor = "→ "
		}

		var statusStyle lipgloss.Style
		switch file.Status {
		case "A":
			statusStyle = utils.FileAddedStyle
		case "M":
			statusStyle = utils.FileModifiedStyle
		case "D":
			statusStyle = utils.FileDeletedStyle
		default:
			statusStyle = utils.NormalItemStyle
		}

		status := statusStyle.Render(file.Status)
		path := utils.DetailsValueStyle.Render(file.Path)

		additions := utils.FileAddedStyle.Render(fmt.Sprintf("+%d", file.Additions))
		deletions := utils.FileDeletedStyle.Render(fmt.Sprintf("-%d", file.Deletions))
		stats := fmt.Sprintf("%s  %s", additions, deletions)

		pathWidth := m.Width - 20
		if len(file.Path) > pathWidth {
			path = utils.DetailsValueStyle.Render("..." + file.Path[len(file.Path)-pathWidth+3:])
		}

		line := fmt.Sprintf("%s%s  %-*s  %s", cursor, status, pathWidth, path, stats)

		if i == m.FileIdx {
			line = utils.SelectedItemStyle.Render(line)
		}

		lines = append(lines, line)
	}

	for len(lines) < 5 {
		lines = append(lines, "")
	}

	return strings.Join(lines, "\n")
}

func (m Model) renderCommitList(commits []types.Commit, selectedIdx int, isActive bool) string {
	var lines []string

	for i, commit := range commits {
		cursor := "  "
		if i == selectedIdx && isActive {
			cursor = "→ "
		}

		hash := utils.HashStyle.Render(commit.Hash)
		msg := commit.Message
		if len(msg) > 28 {
			msg = msg[:25] + "..."
		}

		line := fmt.Sprintf("%s%s %s", cursor, hash, msg)

		if i == selectedIdx && isActive {
			line = utils.SelectedItemStyle.Render(line)
		} else {
			line = utils.NormalItemStyle.Render(line)
		}

		lines = append(lines, line)
	}

	for len(lines) < 5 {
		lines = append(lines, "")
	}

	return strings.Join(lines, "\n")
}

func (m Model) renderDetails(commit types.Commit) string {
	if commit.Hash == "" {
		return utils.DetailsTitleStyle.Render("COMMIT DETAILS") + "\n\n" +
			utils.DetailsValueStyle.Render("No commit selected")
	}

	var b strings.Builder

	b.WriteString(utils.DetailsTitleStyle.Render("COMMIT DETAILS") + "\n\n")

	b.WriteString(utils.HashStyle.Render(commit.Hash) + " ")
	b.WriteString(utils.DetailsValueStyle.Render(commit.Message) + "\n")
	b.WriteString(utils.DetailsLabelStyle.Render("Author: "))
	b.WriteString(utils.DetailsValueStyle.Render(fmt.Sprintf("%s <%s>", commit.Author, commit.Email)) + "\n")
	b.WriteString(utils.DetailsLabelStyle.Render("Date:   "))
	b.WriteString(utils.DetailsValueStyle.Render(commit.Date) + "\n\n")

	b.WriteString(utils.DetailsLabelStyle.Render("Files changed:") + "\n")
	for _, file := range commit.Files {
		var statusStyle lipgloss.Style
		switch file.Status {
		case "A":
			statusStyle = utils.FileAddedStyle
		case "M":
			statusStyle = utils.FileModifiedStyle
		case "D":
			statusStyle = utils.FileDeletedStyle
		default:
			statusStyle = utils.NormalItemStyle
		}
		b.WriteString(fmt.Sprintf("  %s  %s\n", statusStyle.Render(file.Status), utils.DetailsValueStyle.Render(file.Path)))
	}

	return b.String()
}
