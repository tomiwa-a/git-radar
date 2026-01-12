package ui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tomiwa-a/git-radar/internal/types"
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

type Model struct {
	Incoming       []types.Commit
	Outgoing       []types.Commit
	ActivePane     Pane
	IncomingIdx    int
	OutgoingIdx    int
	Width          int
	Height         int
	TargetBranch   string
	SourceBranch   string
	Screen         Screen
	SelectedCommit types.Commit
	FileIdx        int
	Viewport       viewport.Model
	ViewportReady  bool
	GraphCommits   []types.GraphCommit
	GraphIdx       int
	CurrentBranch  string
	Branches       []string
}

func InitialModel() Model {
	return Model{
		GraphCommits:  DummyGraphCommits,
		Branches:      DummyBranches,
		CurrentBranch: "feature/user-validation",
		GraphIdx:      0,
		Incoming:      DummyIncoming,
		Outgoing:      DummyOutgoing,
		ActivePane:    IncomingPane,
		IncomingIdx:   0,
		OutgoingIdx:   0,
		Width:         80,
		Height:        24,
		TargetBranch:  "main",
		SourceBranch:  "feature/user-validation",
		Screen:        GraphScreen,
		FileIdx:       0,
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
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
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

	case "up", "k":
		if m.GraphIdx > 0 {
			m.GraphIdx--
		}

	case "down", "j":
		if m.GraphIdx < len(m.GraphCommits)-1 {
			m.GraphIdx++
		}

	case "enter":
		if len(m.GraphCommits) > 0 {
			gc := m.GraphCommits[m.GraphIdx]
			m.SelectedCommit = m.findCommitByHash(gc.Hash)
			if m.SelectedCommit.Hash != "" {
				m.Screen = CommitDetailScreen
				m.FileIdx = 0
			}
		}

	case "c":
		m.Screen = DivergenceScreen
	}
	return m, nil
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
