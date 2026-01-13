package ui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tomiwa-a/git-radar/internal/ui/screens"
)

const linesPerCommit = 1

func (m Model) updateGraph(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit

	case "?":
		m.ShowLegend = !m.ShowLegend
		return m, nil

	case "esc":
		if m.ShowLegend {
			m.ShowLegend = false
			return m, nil
		}

	case "up", "k":
		if !m.ShowLegend && m.GraphIdx > 0 {
			m.GraphIdx--
			m = m.updateGraphViewportContent()
			m = m.scrollToGraphSelection()
			if len(m.GraphCommits) > 0 {
				m.PendingDetailsHash = m.GraphCommits[m.GraphIdx].FullHash
				return m, m.debounceDetailsCmd(m.GraphCommits[m.GraphIdx].FullHash)
			}
		}
		return m, nil

	case "down", "j":
		if !m.ShowLegend && m.GraphIdx < len(m.GraphCommits)-1 {
			m.GraphIdx++
			m = m.updateGraphViewportContent()
			m = m.scrollToGraphSelection()
			if len(m.GraphCommits) > 0 {
				m.PendingDetailsHash = m.GraphCommits[m.GraphIdx].FullHash
				return m, m.debounceDetailsCmd(m.GraphCommits[m.GraphIdx].FullHash)
			}
		}
		return m, nil

	case "enter":
		if !m.ShowLegend && len(m.GraphCommits) > 0 {
			gc := m.GraphCommits[m.GraphIdx]
			m.SelectedCommit = gc
			if m.SelectedCommit.Hash != "" {
				if len(gc.Files) == 0 && m.GitService != nil {
					parentInfos, files, _ := m.GitService.GetCommitDetails(gc.FullHash)
					m.GraphCommits[m.GraphIdx].ParentInfos = parentInfos
					m.GraphCommits[m.GraphIdx].Files = files
					m.SelectedCommit = m.GraphCommits[m.GraphIdx]
				}
				m.Screen = CommitDetailScreen
				m.FileIdx = 0
			}
		}

	case "c":
		if !m.ShowLegend {
			m.ShowCompareModal = true
			m.CompareModalIdx = 0
		}

	case "pgup", "pgdown", "home", "end":
		if !m.ShowLegend && m.GraphViewportReady {
			var cmd tea.Cmd
			m.GraphViewport, cmd = m.GraphViewport.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

func (m Model) scrollToGraphSelection() Model {
	if !m.GraphViewportReady {
		return m
	}

	selectedLine := m.GraphIdx * linesPerCommit
	viewportHeight := m.GraphViewport.Height
	currentTop := m.GraphViewport.YOffset

	if selectedLine < currentTop {
		m.GraphViewport.SetYOffset(selectedLine)
	}

	if selectedLine >= currentTop+viewportHeight-linesPerCommit {
		m.GraphViewport.SetYOffset(selectedLine - viewportHeight + linesPerCommit + 1)
	}

	return m
}

func (m Model) initGraphViewport() Model {
	headerHeight := 4
	footerHeight := 1
	leftPaneWidth := (m.Width * 60) / 100

	m.GraphViewport = viewport.New(leftPaneWidth, m.Height-headerHeight-footerHeight)
	m.GraphViewport.YPosition = headerHeight

	content := screens.RenderGraphContent(leftPaneWidth, m.GraphCommits, m.GraphIdx)
	m.GraphViewport.SetContent(content)
	m.GraphViewportReady = true

	return m
}

func (m Model) updateGraphViewportContent() Model {
	if m.GraphViewportReady {
		leftPaneWidth := (m.Width * 60) / 100
		content := screens.RenderGraphContent(leftPaneWidth, m.GraphCommits, m.GraphIdx)
		m.GraphViewport.SetContent(content)
	}
	return m
}
