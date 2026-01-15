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
		baseView = screens.RenderGraphWithLegend(m.Width, m.Height, m.GraphCommits, m.GraphIdx, m.CurrentBranch, m.ShowLegend, viewportContent, m.LoadingCommits, m.AlertMessage)
	case CommitDetailScreen:
		displayFiles := m.SelectedCommit.Files
		if m.ShowFilter {
			displayFiles = m.FilteredFiles
		}
		baseView = screens.RenderFileList(m.Width, m.Height, m.SelectedCommit, displayFiles, m.FileIdx, m.ShowFilter, m.FilterInput.Value())
	case DiffViewScreen:
		displayFiles := m.SelectedCommit.Files
		if m.ShowFilter {
			displayFiles = m.FilteredFiles
		}
		baseView = screens.RenderDiffs(m.Width, m.SelectedCommit, displayFiles, m.FileIdx, m.Viewport.View(), m.ShowFilter)
	case DivergenceScreen:
		data := screens.DivergenceData{
			TargetBranch:      m.TargetBranch,
			SourceBranch:      m.SourceBranch,
			MergeBase:         m.MergeBase,
			Incoming:          m.Incoming,
			Outgoing:          m.Outgoing,
			IncomingIdx:       m.IncomingIdx,
			OutgoingIdx:       m.OutgoingIdx,
			ActivePane:        int(m.ActivePane),
			TotalFiles:        m.TotalFiles,
			TotalAdditions:    m.TotalAdditions,
			TotalDeletions:    m.TotalDeletions,
			ConflictFiles:     nil,
			LoadingDivergence: m.LoadingDivergence,
			AlertMessage:      m.AlertMessage,
		}
		baseView = screens.RenderDivergence(m.Width, m.Height, data)
	default:
		baseView = screens.RenderGraph(m.Width, m.GraphCommits, m.GraphIdx, m.CurrentBranch, m.AlertMessage)
	}

	if m.ShowBranchModal {
		modal := screens.RenderBranchModal(m.Width, m.Height, m.Branches, m.BranchModalIdx, m.CurrentBranch)
		return modal
	}

	if m.ShowCompareModal {
		localView := m.CompareLocalPane.View()
		remoteView := m.CompareRemotePane.View()
		modal := screens.RenderCompareModal(m.Width, m.Height, localView, remoteView, m.CompareFilterInput.Value(), int(m.ActiveComparePane))
		return modal
	}

	return baseView
}

func overlayModal(base, modal string, width, height int) string {
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, modal,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("#282A36")))
}
