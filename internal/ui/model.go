package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
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

type DetailsLoadedMsg struct {
	FullHash    string
	ParentInfos []types.ParentInfo
	Files       []types.FileChange
}

type DebounceTickMsg struct {
	FullHash string
}

type DivergenceLoadedMsg struct {
	MergeBase      *types.GraphCommit
	Incoming       []types.GraphCommit
	Outgoing       []types.GraphCommit
	TotalFiles     int
	TotalAdditions int
	TotalDeletions int
}

type ClearAlertMsg struct{}

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
	PreviousScreen     Screen
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
	LoadingDetails     bool
	PendingDetailsHash string
	LoadingDivergence  bool
	MergeBase          *types.GraphCommit
	TotalFiles         int
	TotalAdditions     int
	TotalDeletions     int
	AlertMessage       string
	ShowFilter         bool
	FilterInput        textinput.Model
	FilteredFiles      []types.FileChange
}

func InitialModel() Model {
	return Model{
		GraphCommits:       []types.GraphCommit{},
		Branches:           nil,
		CurrentBranch:      "",
		GraphIdx:           0,
		Incoming:           []types.GraphCommit{},
		Outgoing:           []types.GraphCommit{},
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
		LoadingDetails:     false,
		ShowFilter:         false,
		FilterInput:        textinput.New(),
		FilteredFiles:      nil,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
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
		if len(m.GraphCommits) > 0 {
			m.PendingDetailsHash = m.GraphCommits[0].FullHash
			return m, m.debounceDetailsCmd(m.GraphCommits[0].FullHash)
		}
		return m, nil

	case DebounceTickMsg:
		if msg.FullHash == m.PendingDetailsHash && m.GitService != nil {
			return m, m.loadDetailsCmd(msg.FullHash)
		}
		return m, nil

	case DetailsLoadedMsg:
		for i := range m.GraphCommits {
			if m.GraphCommits[i].FullHash == msg.FullHash {
				m.GraphCommits[i].ParentInfos = msg.ParentInfos
				m.GraphCommits[i].Files = msg.Files
				break
			}
		}
		if m.SelectedCommit.FullHash == msg.FullHash {
			m.SelectedCommit.ParentInfos = msg.ParentInfos
			m.SelectedCommit.Files = msg.Files
		}
		m.LoadingDetails = false
		return m, nil

	case DivergenceLoadedMsg:
		m.MergeBase = msg.MergeBase
		m.Incoming = msg.Incoming
		m.Outgoing = msg.Outgoing
		m.TotalFiles = msg.TotalFiles
		m.TotalAdditions = msg.TotalAdditions
		m.TotalDeletions = msg.TotalDeletions
		m.LoadingDivergence = false
		m.IncomingIdx = 0
		m.OutgoingIdx = 0
		return m, nil

	case ClearAlertMsg:
		m.AlertMessage = ""
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

func (m Model) debounceDetailsCmd(fullHash string) tea.Cmd {
	return tea.Tick(200*time.Millisecond, func(t time.Time) tea.Msg {
		return DebounceTickMsg{FullHash: fullHash}
	})
}

func (m Model) loadDetailsCmd(fullHash string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		parentInfos, files, err := m.GitService.GetCommitDetails(fullHash)
		if err != nil {
			return DetailsLoadedMsg{FullHash: fullHash}
		}
		return DetailsLoadedMsg{
			FullHash:    fullHash,
			ParentInfos: parentInfos,
			Files:       files,
		}
	})
}

func (m Model) loadDivergenceCmd(target, source string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		mergeBase, _ := m.GitService.GetMergeBase(target, source)
		incoming, _ := m.GitService.GetIncomingCommits(target, source)
		outgoing, _ := m.GitService.GetOutgoingCommits(target, source)

		diffStats, _ := m.GitService.GetBranchDiffStats(source, target)
		totalFiles := len(diffStats)
		totalAdds, totalDels := 0, 0
		for _, f := range diffStats {
			totalAdds += f.Additions
			totalDels += f.Deletions
		}

		return DivergenceLoadedMsg{
			MergeBase:      mergeBase,
			Incoming:       incoming,
			Outgoing:       outgoing,
			TotalFiles:     totalFiles,
			TotalAdditions: totalAdds,
			TotalDeletions: totalDels,
		}
	})
}

func clearAlertCmd() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return ClearAlertMsg{}
	})
}
