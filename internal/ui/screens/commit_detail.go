package screens

import (
	"fmt"
	"strings"

	"github.com/tomiwa-a/git-radar/internal/types"
	"github.com/tomiwa-a/git-radar/utils"

	"github.com/charmbracelet/lipgloss"
)

func RenderFileList(width int, commit types.GraphCommit, fileIdx int) string {
	var b strings.Builder

	backHint := utils.DetailsLabelStyle.Render("ESC: back")
	commitInfo := utils.HashStyle.Render("← " + commit.Hash + " ")
	commitMsg := utils.DetailsTitleStyle.Render(commit.Message)
	headerGap := width - lipgloss.Width(backHint) - lipgloss.Width(commitInfo) - lipgloss.Width(commitMsg)
	if headerGap < 0 {
		headerGap = 0
	}
	header := commitInfo + commitMsg + strings.Repeat(" ", headerGap) + backHint
	b.WriteString(header + "\n\n")

	fileCount := len(commit.Files)
	listTitle := utils.DetailsTitleStyle.Render(fmt.Sprintf("FILES CHANGED (%d)", fileCount))
	b.WriteString(listTitle + "\n\n")

	fileListContent := renderFileItems(width, commit.Files, fileIdx)
	fileList := utils.PaneStyle.Width(width - 4).Render(fileListContent)
	b.WriteString(fileList + "\n")

	help := utils.HelpStyle.Render("↑/↓: navigate │ enter: view diff │ ESC: back │ q: quit")
	b.WriteString(help)

	return b.String()
}

func renderFileItems(width int, files []types.FileChange, fileIdx int) string {
	var lines []string

	for i, file := range files {
		cursor := "  "
		if i == fileIdx {
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

		pathWidth := width - 20
		if len(file.Path) > pathWidth {
			path = utils.DetailsValueStyle.Render("..." + file.Path[len(file.Path)-pathWidth+3:])
		}

		line := fmt.Sprintf("%s%s  %-*s  %s", cursor, status, pathWidth, path, stats)

		if i == fileIdx {
			line = utils.SelectedItemStyle.Render(line)
		}

		lines = append(lines, line)
	}

	for len(lines) < 5 {
		lines = append(lines, "")
	}

	return strings.Join(lines, "\n")
}
