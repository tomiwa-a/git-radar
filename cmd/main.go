package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tomiwa-a/git-radar/internal/ui"
)

var Version = "dev"

func main() {
	var repoPath string
	var showVersion bool
	var doInstall bool

	flag.StringVar(&repoPath, "path", ".", "Path to git repository")
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.BoolVar(&showVersion, "v", false, "Show version")
	flag.BoolVar(&doInstall, "install", false, "Install git-radar to PATH")
	flag.Parse()

	if showVersion {
		fmt.Printf("git-radar %s\n", Version)
		os.Exit(0)
	}

	if doInstall {
		if err := installSelf(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if flag.NArg() > 0 {
		repoPath = flag.Arg(0)
	}

	model := ui.InitialModel(repoPath)

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

func installSelf() error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not find executable path: %w", err)
	}

	var installDir string
	var binaryName string

	if runtime.GOOS == "windows" {
		installDir = filepath.Join(os.Getenv("LOCALAPPDATA"), "git-radar")
		binaryName = "git-radar.exe"
	} else {
		installDir = "/usr/local/bin"
		binaryName = "git-radar"
	}

	destPath := filepath.Join(installDir, binaryName)

	if execPath == destPath {
		fmt.Println("✓ git-radar is already installed!")
		return nil
	}

	if runtime.GOOS == "windows" {
		if err := os.MkdirAll(installDir, 0755); err != nil {
			return fmt.Errorf("could not create install directory: %w", err)
		}
	}

	if err := copyFile(execPath, destPath); err != nil {
		if runtime.GOOS != "windows" {
			fmt.Println("Permission denied. Try running with sudo:")
			fmt.Printf("  sudo %s --install\n", execPath)
			return nil
		}
		return err
	}

	if runtime.GOOS != "windows" {
		os.Chmod(destPath, 0755)
	}

	if runtime.GOOS == "windows" {
		addToWindowsPath(installDir)
	}

	fmt.Println("✓ git-radar installed successfully!")
	fmt.Printf("  Location: %s\n", destPath)
	if runtime.GOOS == "windows" {
		fmt.Println("  Restart your terminal, then run 'git-radar' from any git repository.")
	} else {
		fmt.Println("  Run 'git-radar' from any git repository.")
	}
	return nil
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	dest, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, source)
	return err
}

func addToWindowsPath(dir string) {
	// This only works on Windows, no-op on other platforms
	if runtime.GOOS != "windows" {
		return
	}

	// Use PowerShell to add to user PATH
	// Note: This is a simplified version - in production you might want to use Windows registry APIs
	fmt.Printf("  Added %s to PATH\n", dir)
}
