package ui

import (
	"strings"

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

	case "c":
		m.ShowCompareModal = true
		m.CompareModalIdx = 0

	case "enter":
		var commit types.GraphCommit
		if m.ActivePane == IncomingPane && len(m.Incoming) > 0 {
			commit = m.Incoming[m.IncomingIdx]
		} else if m.ActivePane == OutgoingPane && len(m.Outgoing) > 0 {
			commit = m.Outgoing[m.OutgoingIdx]
		}
		if commit.Hash != "" {
			m.SelectedCommit = commit
			m.PreviousScreen = m.Screen
			m.Screen = CommitDetailScreen
			m.ShowFilter = false
			m.FilterInput.SetValue("")
			m.FilteredFiles = nil
			m.FileIdx = 0
			if len(commit.Files) == 0 && m.GitService != nil {
				m.LoadingDetails = true
				return m, m.loadDetailsCmd(commit.FullHash)
			}
		}

	case "y":
		var hash string
		if m.ActivePane == IncomingPane && len(m.Incoming) > 0 {
			hash = m.Incoming[m.IncomingIdx].FullHash
		} else if m.ActivePane == OutgoingPane && len(m.Outgoing) > 0 {
			hash = m.Outgoing[m.OutgoingIdx].FullHash
		}
		if hash != "" {
			copyToClipboard(hash)
			m.AlertMessage = "Hash copied!"
			return m, clearAlertCmd()
		}

	case "esc":
		m.Screen = GraphScreen
	}
	return m, nil
}

func (m Model) updateCommitDetail(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.ShowFilter {
		switch msg.String() {
		case "esc":
			m.ShowFilter = false
			m.FilterInput.SetValue("")
			m.FilteredFiles = nil
			m.FileIdx = 0
			return m, nil
		case "/":
			m.ShowFilter = false
			return m, nil
		case "up", "k":
			if m.FileIdx > 0 {
				m.FileIdx--
			}
			return m, nil
		case "down", "j":
			if len(m.FilteredFiles) > 0 && m.FileIdx < len(m.FilteredFiles)-1 {
				m.FileIdx++
			}
			return m, nil
		case "enter":
			if len(m.FilteredFiles) > 0 {
				m = m.initViewport()
			}
			return m, nil
		}

		var cmd tea.Cmd
		m.FilterInput, cmd = m.FilterInput.Update(msg)

		// Update filtering
		query := strings.ToLower(m.FilterInput.Value())
		if query == "" {
			m.FilteredFiles = m.SelectedCommit.Files
		} else {
			var filtered []types.FileChange
			isExt := strings.HasPrefix(query, ".")
			for _, f := range m.SelectedCommit.Files {
				path := strings.ToLower(f.Path)
				if isExt {
					if strings.HasSuffix(path, query) {
						filtered = append(filtered, f)
					}
				} else {
					if strings.Contains(path, query) {
						filtered = append(filtered, f)
					}
				}
			}
			m.FilteredFiles = filtered
		}

		// Ensure FileIdx is valid
		if m.FileIdx >= len(m.FilteredFiles) {
			m.FileIdx = 0
			if len(m.FilteredFiles) > 0 {
				m.FileIdx = len(m.FilteredFiles) - 1
			}
		}
		if len(m.FilteredFiles) == 0 {
			m.FileIdx = 0
		}

		return m, cmd
	}

	switch msg.String() {
	case "q":
		return m, tea.Quit

	case "esc":
		m.Screen = m.PreviousScreen
		m.SelectedCommit = types.GraphCommit{}
		m.FileIdx = 0

	case "/":
		m.ShowFilter = true
		m.FilterInput.Focus()
		m.FilteredFiles = m.SelectedCommit.Files
		m.FileIdx = 0
		return m, nil

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
		displayFiles := m.SelectedCommit.Files
		if m.ShowFilter {
			displayFiles = m.FilteredFiles
		}
		if m.FileIdx < len(displayFiles)-1 {
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

	displayFiles := m.SelectedCommit.Files
	if m.ShowFilter {
		displayFiles = m.FilteredFiles
	}
	file := displayFiles[m.FileIdx]

	var content string
	if m.GitService != nil {
		diffLines, err := m.GitService.GetFileDiff(m.SelectedCommit.FullHash, file.Path)
		if err != nil {
			content = "Error loading diff: " + err.Error()
		} else if diffLines == nil {
			content = "No changes in this file"
		} else {
			content = utils.RenderDiffLines(diffLines, file.Path)
		}
	} else {
		content = "Git service not available"
	}

	m.Viewport.SetContent(content)
	m.ViewportReady = true

	return m
}
