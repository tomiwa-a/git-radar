package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tomiwa-a/git-radar/internal/types"
)

type dummyDivergence struct {
	Incoming []types.GraphCommit
	Outgoing []types.GraphCommit
}

func getDummyDivergenceCommits() dummyDivergence {
	return dummyDivergence{
		Incoming: []types.GraphCommit{
			{Hash: "def5678", FullHash: "def5678abc123", Message: "fix auth bug", Author: "Tomiwa", Date: "2 days ago", Files: []types.FileChange{{Status: "M", Path: "auth.go", Additions: 10, Deletions: 5}}},
			{Hash: "jkl3456", FullHash: "jkl3456def456", Message: "Merge pull request #4", Author: "Tomiwa", Date: "3 days ago", IsMerge: true},
			{Hash: "pqr1234", FullHash: "pqr1234ghi789", Message: "update dependencies", Author: "Tomiwa", Date: "4 days ago"},
		},
		Outgoing: []types.GraphCommit{
			{Hash: "ghi9012", FullHash: "ghi9012jkl012", Message: "add logging to service", Author: "Tomiwa", Date: "2 hours ago", Files: []types.FileChange{{Status: "A", Path: "logger.go", Additions: 30, Deletions: 0}, {Status: "M", Path: "service.go", Additions: 15, Deletions: 2}}},
			{Hash: "mno7890", FullHash: "mno7890mno345", Message: "refactor git service", Author: "Tomiwa", Date: "5 hours ago", Files: []types.FileChange{{Status: "M", Path: "service.go", Additions: 50, Deletions: 30}}},
			{Hash: "stu5678", FullHash: "stu5678pqr678", Message: "add unit tests", Author: "Tomiwa", Date: "1 day ago"},
			{Hash: "vwx9012", FullHash: "vwx9012stu901", Message: "cleanup unused code", Author: "Tomiwa", Date: "1 day ago"},
			{Hash: "yza3456", FullHash: "yza3456vwx234", Message: "update documentation", Author: "Tomiwa", Date: "2 days ago"},
		},
	}
}

func (m Model) updateBranchModal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "b":
		m.ShowBranchModal = false

	case "up", "k":
		if m.BranchModalIdx > 0 {
			m.BranchModalIdx--
		}

	case "down", "j":
		if m.BranchModalIdx < len(m.Branches)-1 {
			m.BranchModalIdx++
		}

	case "enter":
		if len(m.Branches) > 0 {
			m.CurrentBranch = m.Branches[m.BranchModalIdx].Name
			m.ShowBranchModal = false
			if m.Screen == DivergenceScreen {
				m.SourceBranch = m.CurrentBranch
				m.LoadingDivergence = true
				m.Incoming = nil
				m.Outgoing = nil
				m.MergeBase = nil
				return m, tea.Batch(
					m.loadCommitsCmd(m.CurrentBranch, 100),
					m.loadDivergenceCmd(m.TargetBranch, m.SourceBranch),
				)
			}
			m.LoadingCommits = true
			return m, m.loadCommitsCmd(m.CurrentBranch, 100)
		}
	}
	return m, nil
}

func (m Model) updateCompareModal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	comparableBranches := m.getComparableBranches()

	switch msg.String() {
	case "esc", "c":
		m.ShowCompareModal = false

	case "up", "k":
		if m.CompareModalIdx > 0 {
			m.CompareModalIdx--
		}

	case "down", "j":
		if m.CompareModalIdx < len(comparableBranches)-1 {
			m.CompareModalIdx++
		}

	case "enter":
		if len(comparableBranches) > 0 && m.GitService != nil {
			m.TargetBranch = comparableBranches[m.CompareModalIdx]
			m.SourceBranch = m.CurrentBranch
			m.ShowCompareModal = false
			m.Screen = DivergenceScreen
			m.IncomingIdx = 0
			m.OutgoingIdx = 0
			m.ActivePane = OutgoingPane
			m.LoadingDivergence = true
			m.Incoming = nil
			m.Outgoing = nil
			m.MergeBase = nil
			return m, m.loadDivergenceCmd(m.TargetBranch, m.SourceBranch)
		}
	}
	return m, nil
}
