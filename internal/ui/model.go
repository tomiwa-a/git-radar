package ui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tomiwa-a/git-radar/internal/git"
	"github.com/tomiwa-a/git-radar/internal/types"
	"github.com/tomiwa-a/git-radar/internal/ui/screens"
	"github.com/tomiwa-a/git-radar/utils"
)

type Pane int

const (
	IncomingPane Pane = iota
	OutgoingPane
)

type Screen int

const (
	GraphScreen Screen = iota
	DivergenceScreen
	CommitDetailScreen
	DiffViewScreen
)

type CommitsLoadedMsg struct {
	Commits []types.GraphCommit
}

type Model struct {
	Incoming           []types.Commit
	Outgoing           []types.Commit
	ActivePane         Pane
	IncomingIdx        int
	OutgoingIdx        int
	Width              int
	Height             int
	TargetBranch       string
	SourceBranch       string
	Screen             Screen
	SelectedCommit     types.Commit
	FileIdx            int
	Viewport           viewport.Model
	ViewportReady      bool
	GraphCommits       []types.GraphCommit
	GraphIdx           int
	CurrentBranch      string
	Branches           []types.Branch
	ShowBranchModal    bool
	BranchModalIdx     int
	ShowCompareModal   bool
	CompareModalIdx    int
	ShowLegend         bool
	GraphViewport      viewport.Model
	GraphViewportReady bool
	GitService         *git.Service
	LoadingCommits     bool
}

func InitialModel() Model {
	return Model{
		GraphCommits:       DummyGraphCommits,
		Branches:           nil,
		CurrentBranch:      "",
		GraphIdx:           0,
		Incoming:           DummyIncoming,
		Outgoing:           DummyOutgoing,
		ActivePane:         IncomingPane,
		IncomingIdx:        0,
		OutgoingIdx:        0,
		Width:              80,
		Height:             24,
		TargetBranch:       "",
		SourceBranch:       "",
		Screen:             GraphScreen,
		FileIdx:            0,
		ShowBranchModal:    false,
		BranchModalIdx:     0,
		ShowCompareModal:   false,
		CompareModalIdx:    0,
		ShowLegend:         false,
		GraphViewportReady: false,
		GitService:         nil,
		LoadingCommits:     false,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		// Initialize or resize graph viewport
		if m.Screen == GraphScreen {
			m = m.initGraphViewport()
		}
		return m, nil

	case CommitsLoadedMsg:
		m.GraphCommits = msg.Commits
		m.LoadingCommits = false
		m.GraphIdx = 0
		// Reinitialize viewport with new content
		if m.Screen == GraphScreen {
			m = m.initGraphViewport()
		}
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if m.ShowBranchModal {
			return m.updateBranchModal(msg)
		}

		if m.ShowCompareModal {
			return m.updateCompareModal(msg)
		}

		if msg.String() == "b" {
			m.ShowBranchModal = true
			m.BranchModalIdx = 0
			for i, branch := range m.Branches {
				if branch.Name == m.CurrentBranch {
					m.BranchModalIdx = i
					break
				}
			}
			return m, nil
		}

		switch m.Screen {
		case GraphScreen:
			return m.updateGraph(msg)
		case DivergenceScreen:
			return m.updateDivergence(msg)
		case CommitDetailScreen:
			return m.updateCommitDetail(msg)
		case DiffViewScreen:
			return m.updateDiffs(msg)
		}
	}
	return m, nil
}

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
		}
		return m, nil

	case "down", "j":
		if !m.ShowLegend && m.GraphIdx < len(m.GraphCommits)-1 {
			m.GraphIdx++
			m = m.updateGraphViewportContent()
			m = m.scrollToGraphSelection()
		}
		return m, nil

	case "enter":
		if !m.ShowLegend && len(m.GraphCommits) > 0 {
			gc := m.GraphCommits[m.GraphIdx]
			m.SelectedCommit = m.findCommitByHash(gc.Hash)
			if m.SelectedCommit.Hash != "" {
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
		// Only page keys scroll the viewport directly
		if !m.ShowLegend && m.GraphViewportReady {
			var cmd tea.Cmd
			m.GraphViewport, cmd = m.GraphViewport.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

// linesPerCommit is now 1 for compact single-line format
const linesPerCommit = 1

func (m Model) scrollToGraphSelection() Model {
	if !m.GraphViewportReady {
		return m
	}

	// Calculate where the selected commit is in the content
	selectedLine := m.GraphIdx * linesPerCommit
	viewportHeight := m.GraphViewport.Height
	currentTop := m.GraphViewport.YOffset

	// If selection is above viewport, scroll up
	if selectedLine < currentTop {
		m.GraphViewport.SetYOffset(selectedLine)
	}

	// If selection is below viewport, scroll down
	if selectedLine >= currentTop+viewportHeight-linesPerCommit {
		m.GraphViewport.SetYOffset(selectedLine - viewportHeight + linesPerCommit + 1)
	}

	return m
}

func (m Model) initGraphViewport() Model {
	headerHeight := 4                     // header + panel header + divider
	footerHeight := 1                     // footer
	leftPaneWidth := (m.Width * 60) / 100 // 60% for commits list

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

func (m Model) getComparableBranches() []string {
	var branches []string
	for _, branch := range m.Branches {
		if branch.Name != m.CurrentBranch {
			branches = append(branches, branch.Name)
		}
	}
	return branches
}

func (m Model) findCommitByHash(hash string) types.Commit {
	for _, c := range m.Incoming {
		if c.Hash == hash {
			return c
		}
	}
	for _, c := range m.Outgoing {
		if c.Hash == hash {
			return c
		}
	}
	return types.Commit{Hash: hash, Message: "Commit details", Author: "Unknown", Files: []types.FileChange{}}
}

func (m Model) loadCommitsCmd(branch string, limit int) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		commits, err := m.GitService.GetCommits(branch, limit)
		if err != nil {
			// For simplicity, return empty on error
			return CommitsLoadedMsg{Commits: []types.GraphCommit{}}
		}
		return CommitsLoadedMsg{Commits: commits}
	})
}

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
		m.SelectedCommit = types.Commit{}
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
	code := utils.GetDummyCode(file.Path)
	content := utils.RenderCodeWithLineNumbers(code, file.Path, m.Width)
	m.Viewport.SetContent(content)
	m.ViewportReady = true

	return m
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
