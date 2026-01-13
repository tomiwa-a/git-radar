package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tomiwa-a/git-radar/internal/git"
	"github.com/tomiwa-a/git-radar/internal/ui"
	"github.com/tomiwa-a/git-radar/utils"
)

func main() {
	appStart := time.Now()

	utils.InitLogger()
	defer utils.CloseLogger()

	gitService, err := git.NewService(".")
	if err != nil {
		fmt.Printf("Error: not a git repository: %v\n", err)
		os.Exit(1)
	}

	model := ui.InitialModel()
	model.GitService = gitService

	branches, err := gitService.GetBranches()
	if err != nil {
		fmt.Printf("Warning: could not load branches: %v\n", err)
	} else {
		model.Branches = branches
	}

	currentBranch, err := gitService.GetCurrentBranch()
	if err == nil {
		model.CurrentBranch = currentBranch
	}

	commits, err := gitService.GetCommits(currentBranch, 100)
	if err != nil {
		fmt.Printf("Warning: could not load commits: %v\n", err)
	} else {
		model.GraphCommits = commits
	}

	utils.LogTiming("App startup", time.Since(appStart))

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
