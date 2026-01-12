package utils

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/charmbracelet/lipgloss"
)

func HighlightCode(code string, filename string) string {
	lexer := lexers.Match(filename)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	style := styles.Get("dracula")
	if style == nil {
		style = styles.Fallback
	}

	formatter := formatters.Get("terminal256")
	if formatter == nil {
		formatter = formatters.Fallback
	}

	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return code
	}

	var buf bytes.Buffer
	err = formatter.Format(&buf, style, iterator)
	if err != nil {
		return code
	}

	return buf.String()
}

func RenderCodeWithLineNumbers(code string, filename string, width int) string {
	highlighted := HighlightCode(code, filename)
	lines := strings.Split(highlighted, "\n")

	lineNumStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6272A4")).
		Width(4).
		Align(lipgloss.Right).
		MarginRight(1)

	dividerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#44475A"))

	var result strings.Builder
	for i, line := range lines {
		lineNum := lineNumStyle.Render(fmt.Sprintf("%d", i+1))
		divider := dividerStyle.Render("│")
		result.WriteString(lineNum + divider + " " + line)
		if i < len(lines)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

func GetDummyCode(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".go":
		return `package processor

import (
	"fmt"
	"time"
)

type Payment struct {
	ID        string
	Amount    float64
	Currency  string
	Status    string
	CreatedAt time.Time
}

func ProcessPayment(p *Payment) error {
	if p.Amount <= 0 {
		return fmt.Errorf("invalid amount: %f", p.Amount)
	}

	p.Status = "processing"
	
	// Validate currency
	if !isValidCurrency(p.Currency) {
		return fmt.Errorf("unsupported currency: %s", p.Currency)
	}

	// Process the payment
	err := chargeCard(p)
	if err != nil {
		p.Status = "failed"
		return err
	}

	p.Status = "completed"
	return nil
}

func isValidCurrency(currency string) bool {
	validCurrencies := []string{"USD", "EUR", "GBP"}
	for _, c := range validCurrencies {
		if c == currency {
			return true
		}
	}
	return false
}`

	case ".md":
		return `# Git-Radar

A TUI tool for visualizing Git branch divergence.

## Features

- **Split-pane dashboard** - View incoming and outgoing commits
- **Syntax highlighting** - Code diffs with full color support
- **Function-aware diffs** - See which functions changed

## Installation

` + "```bash" + `
go install github.com/tomiwa-a/git-radar@latest
` + "```" + `

## Usage

Navigate to your git repository and run:

` + "```bash" + `
git-radar
` + "```"

	case ".yaml", ".yml":
		return `rate_limits:
  default:
    requests_per_minute: 60
    burst_size: 10
    
  authenticated:
    requests_per_minute: 1000
    burst_size: 100

  premium:
    requests_per_minute: 10000
    burst_size: 500

endpoints:
  - path: /api/v1/payments
    limit: authenticated
    
  - path: /api/v1/users
    limit: default
    
  - path: /api/v1/admin
    limit: premium`

	default:
		return `// Sample code file
// This is placeholder content for demonstration

function example() {
    console.log("Hello, World!");
    return true;
}

export default example;`
	}
}
