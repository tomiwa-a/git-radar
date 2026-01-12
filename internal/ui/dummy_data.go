package ui

import "github.com/tomiwa-a/git-radar/internal/types"

var DummyGraphCommits = []types.GraphCommit{
	{Hash: "e5f6g7h", Message: "Add user input validation", Author: "You", Date: "2026-01-11 14:30", Branches: []string{"HEAD", "feature/user-validation"}, Parents: []string{"a1b2c3d"}, GraphChars: "* "},
	{Hash: "a1b2c3d", Message: "Fix payment processing bug", Author: "Alice Chen", Date: "2026-01-10 09:15", Branches: []string{"origin/main"}, Parents: []string{"d4e5f6g"}, GraphChars: "* "},
	{Hash: "m2n3o4p", Message: "Add rate limiting middleware", Author: "You", Date: "2026-01-10 17:45", Branches: []string{"feature/rate-limit"}, Parents: []string{"d4e5f6g"}, GraphChars: "│ * "},
	{Hash: "d4e5f6g", Message: "Update README with new API docs", Author: "Bob Smith", Date: "2026-01-09 16:42", Branches: []string{}, Parents: []string{"h7i8j9k"}, GraphChars: "├─┘ "},
	{Hash: "i9j0k1l", Message: "Refactor auth module", Author: "You", Date: "2026-01-11 11:15", Branches: []string{"feature/auth"}, Parents: []string{"h7i8j9k"}, GraphChars: "│ * "},
	{Hash: "h7i8j9k", Message: "Bump dependencies to latest", Author: "Alice Chen", Date: "2026-01-09 11:20", Branches: []string{"main"}, Parents: []string{"q5r6s7t"}, GraphChars: "├─┘ "},
	{Hash: "q5r6s7t", Message: "Fix typo in config loader", Author: "You", Date: "2026-01-10 10:22", Branches: []string{}, Parents: []string{"abc1234"}, GraphChars: "* "},
	{Hash: "abc1234", Message: "Initial commit", Author: "You", Date: "2026-01-01 10:00", Branches: []string{}, Parents: []string{}, GraphChars: "* "},
}

var DummyBranches = []string{"main", "feature/user-validation", "feature/rate-limit", "feature/auth"}

var DummyIncoming = []types.Commit{
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
}

var DummyOutgoing = []types.Commit{
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
}
