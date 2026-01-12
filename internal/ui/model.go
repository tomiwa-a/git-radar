package ui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tomiwa-a/git-radar/internal/git"
	"github.com/tomiwa-a/git-radar/internal/types"
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
	Incoming           []types.GraphCommit
	Outgoing           []types.GraphCommit
	ActivePane         Pane
	IncomingIdx        int
	OutgoingIdx        int
	Width              int
	Height             int
	TargetBranch       string
	SourceBranch       string
	Screen             Screen
	SelectedCommit     types.GraphCommit
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
		if m.Screen == GraphScreen {
			m = m.initGraphViewport()
		}
		return m, nil

	case CommitsLoadedMsg:
		m.GraphCommits = msg.Commits
		m.LoadingCommits = false
		m.GraphIdx = 0
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

func (m Model) getComparableBranches() []string {
	var branches []string
	for _, branch := range m.Branches {
		if branch.Name != m.CurrentBranch {
			branches = append(branches, branch.Name)
		}
	}
	return branches
}

func (m Model) findCommitByHash(hash string) types.GraphCommit {
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
	return types.GraphCommit{Hash: hash, Message: "Commit details", Author: "Unknown", Files: []types.FileChange{}}
}

func (m Model) loadCommitsCmd(branch string, limit int) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		commits, err := m.GitService.GetCommits(branch, limit)
		if err != nil {
			return CommitsLoadedMsg{Commits: []types.GraphCommit{}}
		}
		return CommitsLoadedMsg{Commits: commits}
	})
}
