package ui

import (
	"github.com/tomiwa-a/git-radar/internal/ui/screens"
)

func (m Model) View() string {
	switch m.Screen {
	case GraphScreen:
		return screens.RenderGraph(m.Width, m.GraphCommits, m.GraphIdx, m.CurrentBranch)
	case CommitDetailScreen:
		return screens.RenderFileList(m.Width, m.SelectedCommit, m.FileIdx)
	case DiffViewScreen:
		return screens.RenderDiffs(m.Width, m.SelectedCommit, m.FileIdx, m.Viewport.View())
	case DivergenceScreen:
		return screens.RenderDashboard(
			m.Width,
			m.TargetBranch, m.SourceBranch,
			m.Incoming, m.Outgoing,
			m.IncomingIdx, m.OutgoingIdx,
			int(m.ActivePane),
		)
	default:
		return screens.RenderGraph(m.Width, m.GraphCommits, m.GraphIdx, m.CurrentBranch)
	}
}
