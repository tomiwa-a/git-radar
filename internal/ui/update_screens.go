package ui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tomiwa-a/git-radar/internal/types"
	"github.com/tomiwa-a/git-radar/utils"
)

func (m Model) updateDivergence(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit

	case "left", "h":
		m.ActivePane = IncomingPane

	case "right", "l":
		m.ActivePane = OutgoingPane

	case "up", "k":
		if m.ActivePane == IncomingPane && m.IncomingIdx > 0 {
			m.IncomingIdx--
		} else if m.ActivePane == OutgoingPane && m.OutgoingIdx > 0 {
			m.OutgoingIdx--
		}

	case "down", "j":
		if m.ActivePane == IncomingPane && m.IncomingIdx < len(m.Incoming)-1 {
			m.IncomingIdx++
		} else if m.ActivePane == OutgoingPane && m.OutgoingIdx < len(m.Outgoing)-1 {
			m.OutgoingIdx++
		}

	case "tab":
		if m.ActivePane == IncomingPane {
			m.ActivePane = OutgoingPane
		} else {
			m.ActivePane = IncomingPane
		}

	case "enter":
		if m.ActivePane == IncomingPane && len(m.Incoming) > 0 {
			m.SelectedCommit = m.Incoming[m.IncomingIdx]
		} else if m.ActivePane == OutgoingPane && len(m.Outgoing) > 0 {
			m.SelectedCommit = m.Outgoing[m.OutgoingIdx]
		}
		if m.SelectedCommit.Hash != "" {
			m.Screen = CommitDetailScreen
			m.FileIdx = 0
		}

	case "esc":
		m.Screen = GraphScreen
	}
	return m, nil
}

func (m Model) updateCommitDetail(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit

	case "esc":
		m.Screen = GraphScreen
		m.SelectedCommit = types.GraphCommit{}
		m.FileIdx = 0

	case "up", "k":
		if m.FileIdx > 0 {
			m.FileIdx--
		}

	case "down", "j":
		if m.FileIdx < len(m.SelectedCommit.Files)-1 {
			m.FileIdx++
		}

	case "enter":
		if len(m.SelectedCommit.Files) > 0 {
			m.Screen = DiffViewScreen
			m = m.initViewport()
		}
	}

	return m, nil
}

func (m Model) updateDiffs(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit

	case "esc":
		m.Screen = CommitDetailScreen
		m.FileIdx = 0
		m.ViewportReady = false

	case "left", "h":
		if m.FileIdx > 0 {
			m.FileIdx--
			m = m.initViewport()
		}

	case "right", "l":
		if m.FileIdx < len(m.SelectedCommit.Files)-1 {
			m.FileIdx++
			m = m.initViewport()
		}

	default:
		var cmd tea.Cmd
		m.Viewport, cmd = m.Viewport.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m Model) initViewport() Model {
	headerHeight := 3
	m.Viewport = viewport.New(m.Width, m.Height-headerHeight)
	m.Viewport.YPosition = headerHeight

	file := m.SelectedCommit.Files[m.FileIdx]

	var code string
	if m.GitService != nil {
		content, err := m.GitService.GetFileContent(m.SelectedCommit.FullHash, file.Path)
		if err != nil {
			code = "Error loading file: " + err.Error()
		} else {
			code = content
		}
	} else {
		code = "Git service not available"
	}

	content := utils.RenderCodeWithLineNumbers(code, file.Path, m.Width)
	m.Viewport.SetContent(content)
	m.ViewportReady = true

	return m
}
