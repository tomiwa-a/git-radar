package screens

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tomiwa-a/git-radar/internal/types"
	"github.com/tomiwa-a/git-radar/utils"
)

func RenderDiffs(width int, commit types.Commit, fileIdx int) string {
	var b strings.Builder

	backHint := utils.DetailsLabelStyle.Render("ESC: back")
	commitInfo := utils.HashStyle.Render("← " + commit.Hash + " ")
	commitMsg := utils.DetailsTitleStyle.Render(commit.Message)
	fileName := utils.FileNameStyle.Render(commit.Files[fileIdx].Path)
	headerGap := width - lipgloss.Width(backHint) - lipgloss.Width(commitInfo) - lipgloss.Width(commitMsg)
	if headerGap < 0 {
		headerGap = 0
	}
	header := commitInfo + commitMsg + fileName + strings.Repeat(" ", headerGap) + backHint
	b.WriteString(header + "\n\n")

	b.WriteString(commit.Files[fileIdx].Path + "\n\n")

	return b.String()
}
