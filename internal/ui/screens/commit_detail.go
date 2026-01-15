package screens

import (
	"fmt"
	"strings"

	"github.com/tomiwa-a/git-radar/internal/types"
	"github.com/tomiwa-a/git-radar/utils"

	"github.com/charmbracelet/lipgloss"
)

func RenderFileList(width, height int, commit types.GraphCommit, fileIdx int) string {
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

	// Calculate available height for the file list
	// Header (1) + spacing (2) + title (1) + spacing (2) + help (1)
	reservedHeight := 7
	availableHeight := height - reservedHeight
	if availableHeight < 5 {
		availableHeight = 5
	}

	// renderFileItems will return a fixed number of lines (availableHeight)
	fileListContent := renderFileItems(width, availableHeight, commit.Files, fileIdx)
	fileList := utils.PaneStyle.Width(width - 4).Height(availableHeight).Render(fileListContent)
	b.WriteString(fileList + "\n")

	help := utils.HelpStyle.Render("↑/↓: navigate │ enter: view diff │ ESC: back │ q: quit")
	b.WriteString(help)

	return b.String()
}

func renderFileItems(width, height int, files []types.FileChange, fileIdx int) string {
	var b strings.Builder
	leftPad := "  "

	// Each file item takes 2 lines (file info + empty line)
	maxVisibleItems := height / 2
	if maxVisibleItems < 1 {
		maxVisibleItems = 1
	}

	startIdx := 0
	if fileIdx >= maxVisibleItems {
		startIdx = fileIdx - maxVisibleItems + 1
	}
	endIdx := startIdx + maxVisibleItems
	if endIdx > len(files) {
		endIdx = len(files)
	}

	if startIdx > 0 {
		b.WriteString(utils.HelpStyle.Render(fmt.Sprintf("  ↑ %d more files above", startIdx)) + "\n")
	} else {
		b.WriteString("\n")
	}

	visibleLines := 1 // for the indicator or spacing above
	for i := startIdx; i < endIdx; i++ {
		file := files[i]
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

		rawAdd := fmt.Sprintf("+%4d", file.Additions)
		rawDel := fmt.Sprintf("-%4d", file.Deletions)
		statsRaw := fmt.Sprintf("%s    %s", rawAdd, rawDel)

		addStyle := utils.FileAddedStyle.Copy().Bold(true)
		delStyle := utils.FileDeletedStyle.Copy().Bold(true)
		statsStyled := fmt.Sprintf("%s    %s", addStyle.Render(fmt.Sprintf("+%4d", file.Additions)), delStyle.Render(fmt.Sprintf("-%4d", file.Deletions)))

		statsWidth := len(statsRaw)
		pathWidth := width - 20 - statsWidth
		if pathWidth < 10 {
			pathWidth = 10
		}

		path := file.Path
		if len(path) > pathWidth {
			path = "..." + path[len(path)-pathWidth+3:]
		}
		path = utils.DetailsValueStyle.Render(path)

		line := fmt.Sprintf("%s%s  %-*s  %s", leftPad+cursor, status, pathWidth, path, statsStyled)

		if i == fileIdx {
			line = utils.SelectedItemStyle.Render(line)
		}

		b.WriteString(line + "\n")
		b.WriteString("\n")
		visibleLines += 2
	}

	if endIdx < len(files) {
		b.WriteString(utils.HelpStyle.Render(fmt.Sprintf("  ↓ %d more files below", len(files)-endIdx)) + "\n")
		visibleLines++
	}

	// Fill remaining height with empty lines
	for i := visibleLines; i < height; i++ {
		b.WriteString("\n")
	}

	return b.String()
}
