package screens

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tomiwa-a/git-radar/internal/types"
	"github.com/tomiwa-a/git-radar/utils"
)

// Styles for the new graph UI
var (
	commitDotStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#BD93F9"))
	selectedDotStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B")).Bold(true)
	mergeDotStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF79C6")).Bold(true)
	lineStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#6272A4"))
	hashStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB86C")).Bold(true)
	messageStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#F8F8F2"))
	authorTimeStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#6272A4"))
	localBranchStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B"))
	remoteBranchStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#8BE9FD"))
	mainBranchStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B")).Bold(true)
	mergeTagStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF79C6")).Italic(true)
	selectedBgStyle   = lipgloss.NewStyle().Background(lipgloss.Color("#44475A"))
	treeBranchStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#6272A4"))
)

func RenderGraph(width int, commits []types.GraphCommit, selectedIdx int, currentBranch string) string {
	return RenderGraphWithLegend(width, 24, commits, selectedIdx, currentBranch, false)
}

func RenderGraphWithLegend(width, height int, commits []types.GraphCommit, selectedIdx int, currentBranch string, showLegend bool) string {
	if showLegend {
		return renderLegend(width, height)
	}

	var b strings.Builder

	// Header
	title := utils.TitleStyle.Render(" Git-Radar ")
	branchLabel := utils.DetailsLabelStyle.Render("on: ")
	branchName := utils.BranchStyle.Render(" " + currentBranch + " ")
	help := utils.DetailsLabelStyle.Render("?: help")

	headerGap := width - lipgloss.Width(title) - lipgloss.Width(branchLabel) - lipgloss.Width(branchName) - lipgloss.Width(help) - 4
	if headerGap < 0 {
		headerGap = 0
	}
	header := title + strings.Repeat(" ", headerGap/2) + branchLabel + branchName + strings.Repeat(" ", headerGap/2) + help
	b.WriteString(header + "\n")

	divider := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#44475A")).
		Render(strings.Repeat("─", width))
	b.WriteString(divider + "\n\n")

	// Render commits
	for i, commit := range commits {
		isSelected := i == selectedIdx
		renderCommitBlock(&b, commit, isSelected, width)
	}

	// Footer
	b.WriteString("\n")
	footer := utils.DetailsLabelStyle.Render("↑/↓: navigate │ enter: view commit │ b: branches │ c: compare │ ?: help │ q: quit")
	b.WriteString(footer)

	return b.String()
}

func renderCommitBlock(b *strings.Builder, commit types.GraphCommit, isSelected bool, width int) {
	line := lineStyle.Render("│")

	// Line 1: Dot + Hash + Message
	var dot string
	if isSelected {
		dot = selectedDotStyle.Render("●")
	} else if commit.IsMerge {
		dot = mergeDotStyle.Render("◆")
	} else {
		dot = commitDotStyle.Render("○")
	}

	hash := hashStyle.Render(commit.Hash)
	msg := messageStyle.Render("  " + truncateMessage(commit.Message, 60))

	line1 := "  " + dot + " " + hash + msg
	if isSelected {
		line1 = selectedBgStyle.Render(line1 + strings.Repeat(" ", max(0, width-lipgloss.Width(line1))))
	}
	b.WriteString(line1 + "\n")

	// Line 2: Continuation line + Author + Time
	authorTime := authorTimeStyle.Render(commit.Author + " • " + commit.Date)
	line2 := "  " + line + "          " + authorTime
	if isSelected {
		line2 = selectedBgStyle.Render(line2 + strings.Repeat(" ", max(0, width-lipgloss.Width(line2))))
	}
	b.WriteString(line2 + "\n")

	// Line 3: Merge tag (if merge commit)
	if commit.IsMerge {
		mergeTag := mergeTagStyle.Render("⚡ MERGE COMMIT")
		line3 := "  " + line + "          " + mergeTag
		if isSelected {
			line3 = selectedBgStyle.Render(line3 + strings.Repeat(" ", max(0, width-lipgloss.Width(line3))))
		}
		b.WriteString(line3 + "\n")
	}

	// Lines 4+: Branch labels (if any)
	if len(commit.Branches) > 0 {
		for i, branch := range commit.Branches {
			var prefix string
			if i == len(commit.Branches)-1 {
				prefix = treeBranchStyle.Render("└── ")
			} else {
				prefix = treeBranchStyle.Render("├── ")
			}

			var branchLabel string
			isRemote := strings.HasPrefix(branch, "origin/")
			isMain := branch == "main" || branch == "master" || branch == "origin/main" || branch == "origin/master"

			if isMain {
				branchLabel = mainBranchStyle.Render(branch) + mainBranchStyle.Render(" ★")
			} else if isRemote {
				branchLabel = remoteBranchStyle.Render(branch) + authorTimeStyle.Render("  (remote)")
			} else {
				branchLabel = localBranchStyle.Render(branch)
			}

			branchLine := "  " + line + "          " + prefix + branchLabel
			if isSelected {
				branchLine = selectedBgStyle.Render(branchLine + strings.Repeat(" ", max(0, width-lipgloss.Width(branchLine))))
			}
			b.WriteString(branchLine + "\n")
		}
	}

	// Empty line for spacing
	spacerLine := "  " + line
	b.WriteString(spacerLine + "\n")
}

func renderLegend(width, height int) string {
	modalWidth := 50

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#BD93F9")).
		Padding(1, 2)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#F8F8F2")).
		MarginBottom(1)

	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#BD93F9")).
		Bold(true).
		MarginTop(1)

	itemStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F8F8F2"))

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6272A4"))

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6272A4")).
		MarginTop(1)

	var content strings.Builder

	// Title
	title := titleStyle.Render("LEGEND")
	escHint := utils.DetailsLabelStyle.Render("ESC/?")
	titleLine := title + strings.Repeat(" ", modalWidth-lipgloss.Width(title)-lipgloss.Width(escHint)-6) + escHint
	content.WriteString(titleLine + "\n")

	// Commits section
	content.WriteString("\n" + sectionStyle.Render("COMMITS") + "\n")
	content.WriteString(itemStyle.Render("  ○") + descStyle.Render("  Regular commit") + "\n")
	content.WriteString(selectedDotStyle.Render("  ●") + descStyle.Render("  Selected commit") + "\n")
	content.WriteString(mergeDotStyle.Render("  ◆") + descStyle.Render("  Merge commit") + "\n")

	// Branches section
	content.WriteString("\n" + sectionStyle.Render("BRANCHES") + "\n")
	content.WriteString(localBranchStyle.Render("  feature/xyz") + descStyle.Render("      Local branch (yours)") + "\n")
	content.WriteString(remoteBranchStyle.Render("  origin/xyz") + descStyle.Render("       Remote branch") + "\n")
	content.WriteString(mainBranchStyle.Render("  main ★") + descStyle.Render("            Default branch") + "\n")

	// Graph section
	content.WriteString("\n" + sectionStyle.Render("GRAPH") + "\n")
	content.WriteString(lineStyle.Render("  │") + descStyle.Render("   History continues") + "\n")
	content.WriteString(mergeTagStyle.Render("  ⚡") + descStyle.Render("  Merge point") + "\n")

	// Keys section
	content.WriteString("\n" + sectionStyle.Render("KEYS") + "\n")
	content.WriteString(itemStyle.Render("  ↑/↓ j/k") + descStyle.Render("   Navigate commits") + "\n")
	content.WriteString(itemStyle.Render("  enter") + descStyle.Render("     View commit details") + "\n")
	content.WriteString(itemStyle.Render("  b") + descStyle.Render("         Switch branch") + "\n")
	content.WriteString(itemStyle.Render("  c") + descStyle.Render("         Compare branches") + "\n")
	content.WriteString(itemStyle.Render("  q") + descStyle.Render("         Quit") + "\n")

	content.WriteString("\n" + hintStyle.Render("Press ESC or ? to close"))

	modal := borderStyle.Render(content.String())

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, modal,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("#282A36")))
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
