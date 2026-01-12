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

type Model struct {
	Incoming     []types.Commit
	Outgoing     []types.Commit
	ActivePane   Pane
	IncomingIdx  int
	OutgoingIdx  int
	Width        int
	Height       int
	TargetBranch string
	SourceBranch string
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
					{Status: "M", Path: "src/payments/processor.go"},
					{Status: "M", Path: "src/payments/validator.go"},
				},
			},
			{
				Hash:    "d4e5f6g",
				Message: "Update README with new API docs",
				Author:  "Bob Smith",
				Email:   "bob@example.com",
				Date:    "2026-01-09 16:42",
				Files: []types.FileChange{
					{Status: "M", Path: "README.md"},
				},
			},
			{
				Hash:    "h7i8j9k",
				Message: "Bump dependencies to latest",
				Author:  "Alice Chen",
				Email:   "alice@example.com",
				Date:    "2026-01-09 11:20",
				Files: []types.FileChange{
					{Status: "M", Path: "go.mod"},
					{Status: "M", Path: "go.sum"},
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
					{Status: "A", Path: "src/validation/rules.go"},
					{Status: "M", Path: "src/handlers/user.go"},
				},
			},
			{
				Hash:    "i9j0k1l",
				Message: "Refactor auth module",
				Author:  "You",
				Email:   "you@example.com",
				Date:    "2026-01-11 11:15",
				Files: []types.FileChange{
					{Status: "M", Path: "src/auth/jwt.go"},
					{Status: "M", Path: "src/auth/middleware.go"},
					{Status: "D", Path: "src/auth/legacy.go"},
				},
			},
			{
				Hash:    "m2n3o4p",
				Message: "Add rate limiting middleware",
				Author:  "You",
				Email:   "you@example.com",
				Date:    "2026-01-10 17:45",
				Files: []types.FileChange{
					{Status: "A", Path: "src/middleware/rate_limit.go"},
					{Status: "A", Path: "src/config/limits.yaml"},
					{Status: "M", Path: "src/server.go"},
				},
			},
			{
				Hash:    "q5r6s7t",
				Message: "Fix typo in config loader",
				Author:  "You",
				Email:   "you@example.com",
				Date:    "2026-01-10 10:22",
				Files: []types.FileChange{
					{Status: "M", Path: "src/config/loader.go"},
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
		switch msg.String() {
		case "ctrl+c", "q":
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
		}
	}
	return m, nil
}
