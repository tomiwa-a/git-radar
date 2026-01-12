package ui

import (
	"github.com/tomiwa-a/git-radar/internal/ui/screens"
)

func (m Model) View() string {
	switch m.Screen {
	case FileListScreen:
		return screens.RenderFileList(m.Width, m.SelectedCommit, m.FileIdx)
	case DiffViewScreen:
		return screens.RenderDiffs(m.Width, m.SelectedCommit, m.FileIdx, m.Viewport.View())
	default:
		return screens.RenderDashboard(
			m.Width,
			m.TargetBranch, m.SourceBranch,
			m.Incoming, m.Outgoing,
			m.IncomingIdx, m.OutgoingIdx,
			int(m.ActivePane),
		)
	}
}
