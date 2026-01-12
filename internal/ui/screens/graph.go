package screens

import (
	"fmt"
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
	hashStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB86C")).Bold(true)
	messageStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#F8F8F2"))
	dimStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("#6272A4"))
	localBranchStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B"))
	remoteBranchStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#8BE9FD"))
	mainBranchStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B")).Bold(true)
	mergeTagStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF79C6")).Italic(true)
	selectedBgStyle   = lipgloss.NewStyle().Background(lipgloss.Color("#44475A"))
	paneBorderStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#44475A"))
	sectionTitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#BD93F9")).Bold(true)
	branchCountStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB86C"))
)

func RenderGraph(width int, commits []types.GraphCommit, selectedIdx int, currentBranch string) string {
	return RenderGraphWithLegend(width, 24, commits, selectedIdx, currentBranch, false, "")
}

func RenderGraphWithLegend(width, height int, commits []types.GraphCommit, selectedIdx int, currentBranch string, showLegend bool, viewportContent string) string {
	if showLegend {
		return utils.RenderLegend(
			width, height,
			selectedDotStyle, mergeDotStyle, branchCountStyle, mainBranchStyle, localBranchStyle, remoteBranchStyle, utils.DetailsLabelStyle,
		)
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

	// Calculate panel widths
	leftPaneWidth := (width * 60) / 100         // 60% for commits
	rightPaneWidth := width - leftPaneWidth - 3 // 3 for border

	// Panel headers
	leftHeader := sectionTitleStyle.Render("COMMITS")
	rightHeader := sectionTitleStyle.Render("DETAILS")

	headerLine := leftHeader + strings.Repeat(" ", leftPaneWidth-lipgloss.Width(leftHeader)) +
		paneBorderStyle.Render("│") + " " + rightHeader
	b.WriteString(headerLine + "\n")

	divider := paneBorderStyle.Render(strings.Repeat("─", leftPaneWidth) + "┼" + strings.Repeat("─", rightPaneWidth+1))
	b.WriteString(divider + "\n")

	// Get selected commit for details panel
	var selectedCommit *types.GraphCommit
	if selectedIdx >= 0 && selectedIdx < len(commits) {
		selectedCommit = &commits[selectedIdx]
	}

	// Viewport content on left, details on right
	leftContent := viewportContent
	rightContent := renderDetailsPanel(rightPaneWidth, selectedCommit, height-6)

	// Split viewport content into lines
	leftLines := strings.Split(leftContent, "\n")
	rightLines := strings.Split(rightContent, "\n")

	// Render side by side
	contentHeight := height - 6 // header + panel header + divider + footer
	for i := 0; i < contentHeight; i++ {
		leftLine := ""
		if i < len(leftLines) {
			leftLine = leftLines[i]
		}

		rightLine := ""
		if i < len(rightLines) {
			rightLine = rightLines[i]
		}

		// Pad left line to width
		leftLineWidth := lipgloss.Width(leftLine)
		if leftLineWidth < leftPaneWidth {
			leftLine += strings.Repeat(" ", leftPaneWidth-leftLineWidth)
		}

		b.WriteString(leftLine + paneBorderStyle.Render("│") + " " + rightLine + "\n")
	}

	// Footer
	footer := utils.DetailsLabelStyle.Render("↑/↓: navigate │ enter: view files │ b: branches │ c: compare │ ?: help │ q: quit")
	b.WriteString(footer)

	return b.String()
}

// RenderGraphContent renders compact commit lines for the viewport
func RenderGraphContent(width int, commits []types.GraphCommit, selectedIdx int) string {
	var b strings.Builder

	// Calculate available width for message
	// Format: "● abc1234 message...                    19m  ⚑2"
	// Fixed: dot(2) + hash(8) + space(1) + time(5) + space(1) + branch indicator(4) = ~21
	msgWidth := width - 25
	if msgWidth < 20 {
		msgWidth = 20
	}

	for i, commit := range commits {
		isSelected := i == selectedIdx
		line := renderCompactCommitLine(commit, isSelected, width, msgWidth)
		b.WriteString(line + "\n")
	}

	return b.String()
}

func renderCompactCommitLine(commit types.GraphCommit, isSelected bool, width, msgWidth int) string {
	// Dot
	var dot string
	if isSelected {
		dot = selectedDotStyle.Render("●")
	} else if commit.IsMerge {
		dot = mergeDotStyle.Render("◆")
	} else {
		dot = commitDotStyle.Render("○")
	}

	// Hash
	hash := hashStyle.Render(commit.Hash)

	// Message (truncated)
	msg := utils.TruncateMessage(commit.Message, msgWidth)
	msgStyled := messageStyle.Render(msg)

	// Time (compact)
	timeStr := utils.CompactTime(commit.Date)
	timeStyled := dimStyle.Render(timeStr)

	// Branch indicator
	var branchIndicator string
	if len(commit.Branches) > 0 {
		hasMain := false
		for _, br := range commit.Branches {
			if br == "main" || br == "master" {
				hasMain = true
				break
			}
		}
		if hasMain {
			branchIndicator = mainBranchStyle.Render("★")
		} else {
			branchIndicator = branchCountStyle.Render(fmt.Sprintf("⚑%d", len(commit.Branches)))
		}
	}

	// Build line
	line := " " + dot + " " + hash + " " + msgStyled

	// Calculate padding
	currentWidth := lipgloss.Width(line)
	rightPart := timeStyled
	if branchIndicator != "" {
		rightPart = timeStyled + " " + branchIndicator
	}
	rightWidth := lipgloss.Width(rightPart)

	padding := width - currentWidth - rightWidth - 2
	if padding < 1 {
		padding = 1
	}

	line = line + strings.Repeat(" ", padding) + rightPart

	if isSelected {
		line = selectedBgStyle.Render(line + strings.Repeat(" ", max(0, width-lipgloss.Width(line))))
	}

	return line
}

func renderDetailsPanel(width int, commit *types.GraphCommit, height int) string {
	if commit == nil {
		return dimStyle.Render("  No commit selected")
	}

	var b strings.Builder

	// Hash
	b.WriteString(" " + hashStyle.Render(commit.Hash) + "\n")
	b.WriteString("\n")

	// Full message (wrap if needed)
	msgLines := utils.WrapText(commit.Message, width-2)
	for _, line := range msgLines {
		b.WriteString(" " + messageStyle.Render(line) + "\n")
	}
	b.WriteString("\n")

	// Author and time
	b.WriteString(" " + dimStyle.Render(commit.Author) + "\n")
	b.WriteString(" " + dimStyle.Render(commit.Date) + "\n")
	b.WriteString("\n")

	// Divider
	b.WriteString(" " + paneBorderStyle.Render(strings.Repeat("─", width-2)) + "\n")
	b.WriteString("\n")

	// Branches section
	if len(commit.Branches) > 0 {
		b.WriteString(" " + sectionTitleStyle.Render("BRANCHES") + "\n")
		for _, branch := range commit.Branches {
			isRemote := strings.HasPrefix(branch, "origin/")
			isMain := branch == "main" || branch == "master" || branch == "origin/main" || branch == "origin/master"

			var branchLine string
			if isMain {
				branchLine = " • " + mainBranchStyle.Render(branch+" ★")
			} else if isRemote {
				branchLine = " • " + remoteBranchStyle.Render(branch) + dimStyle.Render(" (remote)")
			} else {
				branchLine = " • " + localBranchStyle.Render(branch)
			}
			b.WriteString(branchLine + "\n")
		}
		b.WriteString("\n")
	}

	// Parents section (for merge commits)
	if len(commit.Parents) > 1 {
		b.WriteString(" " + sectionTitleStyle.Render("PARENTS") + " " + mergeTagStyle.Render("(merge)") + "\n")
		for _, parent := range commit.Parents {
			short := parent
			if len(parent) > 7 {
				short = parent[:7]
			}
			b.WriteString(" " + dimStyle.Render("├ ") + hashStyle.Render(short) + "\n")
		}
		b.WriteString("\n")
	}

	return b.String()
}
