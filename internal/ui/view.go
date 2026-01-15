package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/tomiwa-a/git-radar/internal/ui/screens"
)

func (m Model) View() string {
	var baseView string

	switch m.Screen {
	case GraphScreen:
		viewportContent := ""
		if m.GraphViewportReady {
			viewportContent = m.GraphViewport.View()
		}
		baseView = screens.RenderGraphWithLegend(m.Width, m.Height, m.GraphCommits, m.GraphIdx, m.CurrentBranch, m.ShowLegend, viewportContent, m.LoadingCommits)
	case CommitDetailScreen:
		baseView = screens.RenderFileList(m.Width, m.SelectedCommit, m.FileIdx)
	case DiffViewScreen:
		baseView = screens.RenderDiffs(m.Width, m.SelectedCommit, m.FileIdx, m.Viewport.View())
	case DivergenceScreen:
		data := screens.DivergenceData{
			TargetBranch:   m.TargetBranch,
			SourceBranch:   m.SourceBranch,
			MergeBase:      nil,
			Incoming:       m.Incoming,
			Outgoing:       m.Outgoing,
			IncomingIdx:    m.IncomingIdx,
			OutgoingIdx:    m.OutgoingIdx,
			ActivePane:     int(m.ActivePane),
			TotalFiles:     0,
			TotalAdditions: 0,
			TotalDeletions: 0,
			ConflictFiles:  nil,
		}
		baseView = screens.RenderDivergence(m.Width, m.Height, data)
	default:
		baseView = screens.RenderGraph(m.Width, m.GraphCommits, m.GraphIdx, m.CurrentBranch)
	}

	if m.ShowBranchModal {
		modal := screens.RenderBranchModal(m.Width, m.Height, m.Branches, m.BranchModalIdx, m.CurrentBranch)
		return modal
	}

	if m.ShowCompareModal {
		comparableBranches := m.getComparableBranches()
		modal := screens.RenderCompareModal(m.Width, m.Height, comparableBranches, m.CompareModalIdx, m.CurrentBranch)
		return modal
	}

	return baseView
}

func overlayModal(base, modal string, width, height int) string {
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, modal,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("#282A36")))
}
