package ui

import (
	"github.com/tomiwa-a/git-radar/internal/types"

	tea "github.com/charmbracelet/bubbletea"
)

type Pane int

const (
	IncomingPane Pane = iota
	OutgoingPane
)

type Screen int

const (
	DashboardScreen Screen = iota
	FileListScreen
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
}

func InitialModel() Model {
	return Model{
		Incoming: []types.Commit{
			{
				Hash:    "a1b2c3d",
				Message: "Fix payment processing bug",
				Author:  "Alice Chen",
				Email:   "alice@example.com",
				Date:    "2026-01-10 09:15",
				Files: []types.FileChange{
					{Status: "M", Path: "src/payments/processor.go", Additions: 23, Deletions: 8},
					{Status: "M", Path: "src/payments/validator.go", Additions: 15, Deletions: 3},
				},
			},
			{
				Hash:    "d4e5f6g",
				Message: "Update README with new API docs",
				Author:  "Bob Smith",
				Email:   "bob@example.com",
				Date:    "2026-01-09 16:42",
				Files: []types.FileChange{
					{Status: "M", Path: "README.md", Additions: 45, Deletions: 12},
				},
			},
			{
				Hash:    "h7i8j9k",
				Message: "Bump dependencies to latest",
				Author:  "Alice Chen",
				Email:   "alice@example.com",
				Date:    "2026-01-09 11:20",
				Files: []types.FileChange{
					{Status: "M", Path: "go.mod", Additions: 5, Deletions: 5},
					{Status: "M", Path: "go.sum", Additions: 120, Deletions: 98},
				},
			},
		},
		Outgoing: []types.Commit{
			{
				Hash:    "e5f6g7h",
				Message: "Add user input validation",
				Author:  "You",
				Email:   "you@example.com",
				Date:    "2026-01-11 14:30",
				Files: []types.FileChange{
					{Status: "A", Path: "src/validation/rules.go", Additions: 87, Deletions: 0},
					{Status: "M", Path: "src/handlers/user.go", Additions: 34, Deletions: 12},
				},
			},
			{
				Hash:    "i9j0k1l",
				Message: "Refactor auth module",
				Author:  "You",
				Email:   "you@example.com",
				Date:    "2026-01-11 11:15",
				Files: []types.FileChange{
					{Status: "M", Path: "src/auth/jwt.go", Additions: 56, Deletions: 23},
					{Status: "M", Path: "src/auth/middleware.go", Additions: 28, Deletions: 15},
					{Status: "D", Path: "src/auth/legacy.go", Additions: 0, Deletions: 145},
				},
			},
			{
				Hash:    "m2n3o4p",
				Message: "Add rate limiting middleware",
				Author:  "You",
				Email:   "you@example.com",
				Date:    "2026-01-10 17:45",
				Files: []types.FileChange{
					{Status: "A", Path: "src/middleware/rate_limit.go", Additions: 142, Deletions: 0},
					{Status: "A", Path: "src/config/limits.yaml", Additions: 28, Deletions: 0},
					{Status: "M", Path: "src/server.go", Additions: 12, Deletions: 3},
				},
			},
			{
				Hash:    "q5r6s7t",
				Message: "Fix typo in config loader",
				Author:  "You",
				Email:   "you@example.com",
				Date:    "2026-01-10 10:22",
				Files: []types.FileChange{
					{Status: "M", Path: "src/config/loader.go", Additions: 1, Deletions: 1},
				},
			},
		},
		ActivePane:   IncomingPane,
		IncomingIdx:  0,
		OutgoingIdx:  0,
		Width:        80,
		Height:       24,
		TargetBranch: "main",
		SourceBranch: "feature/user-validation",
		Screen:       DashboardScreen,
		FileIdx:      0,
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
		case DashboardScreen:
			return m.updateDashboard(msg)
		case FileListScreen:
			return m.updateFileList(msg)
		case DiffViewScreen:
			return m.updateDiffs(msg)
		}
	}
	return m, nil
}

func (m Model) updateDashboard(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
			m.Screen = FileListScreen
			m.FileIdx = 0
		}
	}
	return m, nil
}

func (m Model) updateFileList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit

	case "esc":
		m.Screen = DashboardScreen
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
		}

	}

	return m, nil
}

func (m Model) updateDiffs(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit

	case "esc":
		m.Screen = FileListScreen
		// m.SelectedCommit = types.Commit{}
		m.FileIdx = 0

	case "up", "k":
		if m.FileIdx > 0 {
			m.FileIdx--
		}

	case "down", "j":
		if m.FileIdx < len(m.SelectedCommit.Files)-1 {
			m.FileIdx++
		}
	}
	return m, nil
}
