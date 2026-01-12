# git-radar

A terminal-based Git visualization tool built with Go and the [Charm](https://charm.sh/) ecosystem.

![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)

## Why?

Because `git log --graph --oneline --all` and `git diff` aren't simple enough.

Navigate commits, switch branches, view diffs – all in one interactive terminal UI.

## Features

- **Commit Graph** – Browse commit history with branch labels and merge indicators
- **Branch Switching** – Quick branch navigation with `b` key
- **Branch Comparison** – Compare divergence between branches
- **Commit Details** – View file changes, additions, and deletions per commit
- **Diff Viewer** – Syntax-highlighted code diffs with line numbers

## Installation

```bash
go install github.com/tomiwa-a/git-radar@latest
```

Or build from source:

```bash
git clone https://github.com/tomiwa-a/git-radar.git
cd git-radar
go build -o git-radar ./cmd
```

## Usage

Navigate to any Git repository and run:

```bash
git-radar
```

## Keybindings

### Global

| Key      | Action               |
| -------- | -------------------- |
| `q`      | Quit                 |
| `Ctrl+C` | Force quit           |
| `b`      | Open branch switcher |

### Graph View

| Key         | Action                      |
| ----------- | --------------------------- |
| `j` / `↓`   | Move down                   |
| `k` / `↑`   | Move up                     |
| `Enter`     | View commit details         |
| `c`         | Compare with another branch |
| `?`         | Toggle legend               |
| `PgUp/PgDn` | Scroll viewport             |

### Commit Detail View

| Key       | Action               |
| --------- | -------------------- |
| `j` / `↓` | Select next file     |
| `k` / `↑` | Select previous file |
| `Enter`   | View file diff       |
| `Esc`     | Back to graph        |

### Diff View

| Key       | Action                 |
| --------- | ---------------------- |
| `h` / `←` | Previous file          |
| `l` / `→` | Next file              |
| `j/k`     | Scroll diff            |
| `Esc`     | Back to commit details |

### Divergence View

| Key       | Action                           |
| --------- | -------------------------------- |
| `Tab`     | Switch between incoming/outgoing |
| `h` / `←` | Select incoming pane             |
| `l` / `→` | Select outgoing pane             |
| `Enter`   | View commit details              |
| `Esc`     | Back to graph                    |

## Project Structure

```
git-radar/
├── cmd/                    # Application entry point
├── internal/
│   ├── git/                # Git operations (go-git wrapper)
│   ├── types/              # Domain types
│   └── ui/                 # TUI components
│       ├── model.go        # Core model and update loop
│       ├── view.go         # Main view renderer
│       ├── update_graph.go # Graph screen handlers
│       ├── update_modals.go# Modal handlers
│       ├── update_screens.go# Screen-specific handlers
│       └── screens/        # Screen renderers
└── utils/                  # Utilities (time formatting, code rendering)
```

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) – TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) – Styling
- [go-git](https://github.com/go-git/go-git) – Pure Go git implementation
- [Chroma](https://github.com/alecthomas/chroma) – Syntax highlighting

## License

MIT
