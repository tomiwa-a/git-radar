package ui

import tea "github.com/charmbracelet/bubbletea"

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
		if len(comparableBranches) > 0 {
			m.TargetBranch = comparableBranches[m.CompareModalIdx]
			m.SourceBranch = m.CurrentBranch
			m.ShowCompareModal = false
			m.Screen = DivergenceScreen
		}
	}
	return m, nil
}
