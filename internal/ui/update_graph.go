package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tomiwa-a/git-radar/internal/types"
	"github.com/tomiwa-a/git-radar/internal/ui/screens"
	"golang.design/x/clipboard"
)

func copyToClipboard(text string) {
	clipboard.Init()
	clipboard.Write(clipboard.FmtText, []byte(text))
}

const linesPerCommit = 1

func (m Model) updateGraph(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit

	case "?":
		m.ShowLegend = !m.ShowLegend
		return m, nil

	case "esc":
		if m.ShowGraphSearch {
			m.ShowGraphSearch = false
			m.GraphSearchInput.SetValue("")
			m.FilteredGraphCommits = nil
			m.GraphIdx = 0
			m = m.updateGraphViewportContent()
			return m, nil
		}
		if m.ShowLegend {
			m.ShowLegend = false
			return m, nil
		}

	case "/":
		if !m.ShowLegend && !m.ShowGraphSearch {
			m.ShowGraphSearch = true
			m.GraphSearchInput.Focus()
			m.GraphSearchInput.SetValue("")
			m.FilteredGraphCommits = nil
			return m, nil
		}

	case "up", "k":
		if m.ShowGraphSearch {
			return m, nil
		}
		if !m.ShowLegend {
			commits := m.getDisplayCommits()
			if m.GraphIdx > 0 {
				m.GraphIdx--
				m = m.updateGraphViewportContent()
				m = m.scrollToGraphSelection()
				if len(commits) > 0 && m.GraphIdx < len(commits) {
					m.PendingDetailsHash = commits[m.GraphIdx].FullHash
					return m, m.debounceDetailsCmd(commits[m.GraphIdx].FullHash)
				}
			}
		}
		return m, nil

	case "down", "j":
		if m.ShowGraphSearch {
			return m, nil
		}
		if !m.ShowLegend {
			commits := m.getDisplayCommits()
			if m.GraphIdx < len(commits)-1 {
				m.GraphIdx++
				m = m.updateGraphViewportContent()
				m = m.scrollToGraphSelection()
				if len(commits) > 0 && m.GraphIdx < len(commits) {
					m.PendingDetailsHash = commits[m.GraphIdx].FullHash
					return m, m.debounceDetailsCmd(commits[m.GraphIdx].FullHash)
				}
			}
		}
		return m, nil

	case "enter":
		if m.ShowGraphSearch {
			m.ShowGraphSearch = false
			return m, nil
		}
		if !m.ShowLegend {
			commits := m.getDisplayCommits()
			if len(commits) > 0 && m.GraphIdx < len(commits) {
				gc := commits[m.GraphIdx]
				m.SelectedCommit = gc
				if m.SelectedCommit.Hash != "" {
					if len(gc.Files) == 0 && m.GitService != nil {
						parentInfos, files, _ := m.GitService.GetCommitDetails(gc.FullHash)
						for i := range m.GraphCommits {
							if m.GraphCommits[i].FullHash == gc.FullHash {
								m.GraphCommits[i].ParentInfos = parentInfos
								m.GraphCommits[i].Files = files
								m.SelectedCommit = m.GraphCommits[i]
								break
							}
						}
					}
					m.PreviousScreen = m.Screen
					m.Screen = CommitDetailScreen
					m.ShowFilter = false
					m.FilterInput.SetValue("")
					m.FilteredFiles = nil
					m.FileIdx = 0
					m.ShowGraphSearch = false
					m.GraphSearchInput.SetValue("")
					m.FilteredGraphCommits = nil
				}
			}
		}

	case "c":
		if !m.ShowLegend {
			m.ShowCompareModal = true
			m.CompareModalIdx = 0
			m.ActiveComparePane = LocalComparePane
			m.CompareFilterInput.SetValue("")
			m.CompareFilterInput.Focus()

			// Split branches
			m.LocalBranches = nil
			m.RemoteBranches = nil
			for _, b := range m.Branches {
				if b.Name == m.CurrentBranch {
					continue
				}
				if b.IsRemote {
					m.RemoteBranches = append(m.RemoteBranches, b)
				} else {
					m.LocalBranches = append(m.LocalBranches, b)
				}
			}
			m.FilteredLocal = m.LocalBranches
			m.FilteredRemote = m.RemoteBranches

			// Initialize viewports
			m = m.initCompareViewports()
		}

	case "y":
		if !m.ShowLegend && len(m.GraphCommits) > 0 {
			hash := m.GraphCommits[m.GraphIdx].FullHash
			copyToClipboard(hash)
			m.AlertMessage = "Hash copied!"
			return m, clearAlertCmd()
		}

	case "pgup", "pgdown", "home", "end":
		if !m.ShowLegend && !m.ShowGraphSearch && m.GraphViewportReady {
			var cmd tea.Cmd
			m.GraphViewport, cmd = m.GraphViewport.Update(msg)
			return m, cmd
		}

	default:
		if m.ShowGraphSearch {
			var cmd tea.Cmd
			m.GraphSearchInput, cmd = m.GraphSearchInput.Update(msg)
			m = m.filterGraphCommits()
			m.GraphIdx = 0
			m = m.updateGraphViewportContent()
			return m, cmd
		}
	}
	return m, nil
}

func (m Model) getDisplayCommits() []types.GraphCommit {
	if m.ShowGraphSearch && len(m.FilteredGraphCommits) > 0 {
		return m.FilteredGraphCommits
	}
	if m.ShowGraphSearch && m.GraphSearchInput.Value() != "" {
		return m.FilteredGraphCommits
	}
	return m.GraphCommits
}

func (m Model) filterGraphCommits() Model {
	query := strings.ToLower(m.GraphSearchInput.Value())
	if query == "" {
		m.FilteredGraphCommits = nil
		return m
	}

	var filtered []types.GraphCommit
	for _, c := range m.GraphCommits {
		if strings.Contains(strings.ToLower(c.Message), query) ||
			strings.Contains(strings.ToLower(c.Hash), query) ||
			strings.Contains(strings.ToLower(c.Author), query) {
			filtered = append(filtered, c)
		}
	}
	m.FilteredGraphCommits = filtered
	return m
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
	if m.ShowGraphSearch {
		headerHeight = 5
	}
	leftPaneWidth := (m.Width * 60) / 100

	m.GraphViewport = viewport.New(leftPaneWidth, m.Height-headerHeight-footerHeight)
	m.GraphViewport.YPosition = headerHeight

	commits := m.getDisplayCommits()
	content := screens.RenderGraphContent(leftPaneWidth, commits, m.GraphIdx)
	m.GraphViewport.SetContent(content)
	m.GraphViewportReady = true

	return m
}

func (m Model) updateGraphViewportContent() Model {
	if m.GraphViewportReady {
		leftPaneWidth := (m.Width * 60) / 100
		commits := m.getDisplayCommits()
		content := screens.RenderGraphContent(leftPaneWidth, commits, m.GraphIdx)
		m.GraphViewport.SetContent(content)
	}
	return m
}

func (m Model) initCompareViewports() Model {
	modalWidth := int(float64(m.Width) * 0.8)
	modalHeight := int(float64(m.Height) * 0.7)
	paneWidth := (modalWidth - 6) / 2 // Borders and divider
	paneHeight := modalHeight - 7     // Headers, filter, footer

	m.CompareLocalPane = viewport.New(paneWidth, paneHeight)
	m.CompareRemotePane = viewport.New(paneWidth, paneHeight)

	return m.updateCompareViewportContent()
}

func (m Model) updateCompareViewportContent() Model {
	modalWidth := int(float64(m.Width) * 0.8)
	paneWidth := (modalWidth - 6) / 2

	m = m.scrollToCompareSelection()

	localContent := screens.RenderBranchListContent(paneWidth, m.FilteredLocal, m.CompareModalIdx, m.ActiveComparePane == LocalComparePane)
	m.CompareLocalPane.SetContent(localContent)

	remoteContent := screens.RenderBranchListContent(paneWidth, m.FilteredRemote, m.CompareModalIdx, m.ActiveComparePane == RemoteComparePane)
	m.CompareRemotePane.SetContent(remoteContent)

	return m
}

func (m Model) scrollToCompareSelection() Model {
	if m.ActiveComparePane == LocalComparePane {
		selectedLine := m.CompareModalIdx
		viewportHeight := m.CompareLocalPane.Height
		currentTop := m.CompareLocalPane.YOffset

		if selectedLine < currentTop {
			m.CompareLocalPane.SetYOffset(selectedLine)
		} else if selectedLine >= currentTop+viewportHeight {
			m.CompareLocalPane.SetYOffset(selectedLine - viewportHeight + 1)
		}
	} else {
		selectedLine := m.CompareModalIdx
		viewportHeight := m.CompareRemotePane.Height
		currentTop := m.CompareRemotePane.YOffset

		if selectedLine < currentTop {
			m.CompareRemotePane.SetYOffset(selectedLine)
		} else if selectedLine >= currentTop+viewportHeight {
			m.CompareRemotePane.SetYOffset(selectedLine - viewportHeight + 1)
		}
	}
	return m
}

func (m Model) initBranchViewports() Model {
	modalWidth := int(float64(m.Width) * 0.8)
	modalHeight := int(float64(m.Height) * 0.7)
	paneWidth := (modalWidth - 6) / 2
	paneHeight := modalHeight - 7

	m.BranchLocalPane = viewport.New(paneWidth, paneHeight)
	m.BranchRemotePane = viewport.New(paneWidth, paneHeight)

	return m.updateBranchViewportContent()
}

func (m Model) updateBranchViewportContent() Model {
	modalWidth := int(float64(m.Width) * 0.8)
	paneWidth := (modalWidth - 6) / 2

	m = m.scrollToBranchSelection()

	localContent := screens.RenderBranchListContent(paneWidth, m.BranchFilteredLocal, m.BranchModalIdx, m.ActiveBranchPane == LocalComparePane)
	m.BranchLocalPane.SetContent(localContent)

	remoteContent := screens.RenderBranchListContent(paneWidth, m.BranchFilteredRemote, m.BranchModalIdx, m.ActiveBranchPane == RemoteComparePane)
	m.BranchRemotePane.SetContent(remoteContent)

	return m
}

func (m Model) scrollToBranchSelection() Model {
	if m.ActiveBranchPane == LocalComparePane {
		selectedLine := m.BranchModalIdx
		viewportHeight := m.BranchLocalPane.Height
		currentTop := m.BranchLocalPane.YOffset

		if selectedLine < currentTop {
			m.BranchLocalPane.SetYOffset(selectedLine)
		} else if selectedLine >= currentTop+viewportHeight {
			m.BranchLocalPane.SetYOffset(selectedLine - viewportHeight + 1)
		}
	} else {
		selectedLine := m.BranchModalIdx
		viewportHeight := m.BranchRemotePane.Height
		currentTop := m.BranchRemotePane.YOffset

		if selectedLine < currentTop {
			m.BranchRemotePane.SetYOffset(selectedLine)
		} else if selectedLine >= currentTop+viewportHeight {
			m.BranchRemotePane.SetYOffset(selectedLine - viewportHeight + 1)
		}
	}
	return m
}
