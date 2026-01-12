package utils

import "github.com/charmbracelet/lipgloss"

var (
	PrimaryColor   = lipgloss.Color("#FF5FAF")
	SecondaryColor = lipgloss.Color("#585858")
	IncomingColor  = lipgloss.Color("#00FF87")
	OutgoingColor  = lipgloss.Color("#87CEFA")
	CursorColor    = lipgloss.Color("#FF87D7")
	DimColor       = lipgloss.Color("#626262")

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(PrimaryColor).
			Padding(0, 1)

	BranchStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#5F5FFF")).
			Padding(0, 1)

	PaneStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(SecondaryColor).
			Padding(0, 1)

	ActivePaneStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(PrimaryColor).
			Padding(0, 1)

	PaneTitleIncoming = lipgloss.NewStyle().
				Bold(true).
				Foreground(IncomingColor)

	PaneTitleOutgoing = lipgloss.NewStyle().
				Bold(true).
				Foreground(OutgoingColor)

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(CursorColor).
				Bold(true)

	NormalItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D0D0D0"))

	HashStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFAF00"))

	DetailsStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(SecondaryColor).
			Padding(0, 1)

	DetailsTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(PrimaryColor)

	DetailsLabelStyle = lipgloss.NewStyle().
				Foreground(DimColor)

	DetailsValueStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#D0D0D0"))

	FileAddedStyle = lipgloss.NewStyle().
			Foreground(IncomingColor)

	FileModifiedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFAF00"))

	FileDeletedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF0000"))

	HelpStyle = lipgloss.NewStyle().
			Foreground(DimColor).
			Padding(0, 1)
)
