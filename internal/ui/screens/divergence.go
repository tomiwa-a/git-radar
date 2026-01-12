package screens

import (
	"fmt"
	"strings"

	"github.com/tomiwa-a/git-radar/internal/types"
	"github.com/tomiwa-a/git-radar/utils"

	"github.com/charmbracelet/lipgloss"
)

func RenderDashboard(
	width int,
	targetBranch, sourceBranch string,
	incoming, outgoing []types.Commit,
	incomingIdx, outgoingIdx int,
	activePane int,
) string {
	var b strings.Builder

	paneWidth := (width - 4) / 2
	if paneWidth < 30 {
		paneWidth = 30
	}

	title := utils.TitleStyle.Render(" Git-Radar ")
	branches := utils.BranchStyle.Render(fmt.Sprintf(" %s ← %s ", targetBranch, sourceBranch))
	headerGap := width - lipgloss.Width(title) - lipgloss.Width(branches)
	if headerGap < 0 {
		headerGap = 0
	}
	header := title + strings.Repeat(" ", headerGap) + branches
	b.WriteString(header + "\n\n")

	incomingContent := renderCommitList(incoming, incomingIdx, activePane == 0)
	outgoingContent := renderCommitList(outgoing, outgoingIdx, activePane == 1)

	incomingTitle := utils.PaneTitleIncoming.Render(fmt.Sprintf("⬇ INCOMING (Behind by %d)", len(incoming)))
	outgoingTitle := utils.PaneTitleOutgoing.Render(fmt.Sprintf("⬆ OUTGOING (Ahead by %d)", len(outgoing)))

	var leftPane, rightPane string
	if activePane == 0 {
		leftPane = utils.ActivePaneStyle.Width(paneWidth).Render(incomingTitle + "\n\n" + incomingContent)
		rightPane = utils.PaneStyle.Width(paneWidth).Render(outgoingTitle + "\n\n" + outgoingContent)
	} else {
		leftPane = utils.PaneStyle.Width(paneWidth).Render(incomingTitle + "\n\n" + incomingContent)
		rightPane = utils.ActivePaneStyle.Width(paneWidth).Render(outgoingTitle + "\n\n" + outgoingContent)
	}

	panes := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)
	b.WriteString(panes + "\n")

	var selectedCommit types.Commit
	if activePane == 0 && len(incoming) > 0 {
		selectedCommit = incoming[incomingIdx]
	} else if activePane == 1 && len(outgoing) > 0 {
		selectedCommit = outgoing[outgoingIdx]
	}

	detailsContent := renderDetails(selectedCommit)
	details := utils.DetailsStyle.Width(width - 4).Render(detailsContent)
	b.WriteString(details + "\n")

	help := utils.HelpStyle.Render("←/→: switch pane │ ↑/↓: navigate │ enter: view files │ q: quit")
	b.WriteString(help)

	return b.String()
}

func renderCommitList(commits []types.Commit, selectedIdx int, isActive bool) string {
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

func renderDetails(commit types.Commit) string {
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
