package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tomiwa-a/git-radar/internal/git"
	"github.com/tomiwa-a/git-radar/internal/ui"
)

func main() {
	// Initialize git service
	gitService, err := git.NewService(".")
	if err != nil {
		fmt.Printf("Error: not a git repository: %v\n", err)
		os.Exit(1)
	}

	// Create model with git data
	model := ui.InitialModel()

	// Load branches
	branches, err := gitService.GetBranches()
	if err != nil {
		fmt.Printf("Warning: could not load branches: %v\n", err)
	} else {
		model.Branches = branches
	}

	// Get current branch
	currentBranch, err := gitService.GetCurrentBranch()
	if err == nil {
		model.CurrentBranch = currentBranch
	}

	// Load commits for graph
	commits, err := gitService.GetCommits(100)
	if err != nil {
		fmt.Printf("Warning: could not load commits: %v\n", err)
	} else {
		model.GraphCommits = commits
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
