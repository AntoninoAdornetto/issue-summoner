package ui

import "github.com/charmbracelet/lipgloss"

var (
	BackgroundStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#2A2A2E")) // A dark, almost black background

	PrimaryTextStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#E5E5E6")). // Light grey for primary text
				Bold(true)

	SecondaryTextStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFA658")). // Orange (Carbon Fox's orange)
				Italic(true)

	AccentTextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3B82F6")). // (Carbon Fox's blue)
			Bold(true)

	SuccessTextStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#10B981")).Bold(true) // Green

	ErrorTextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EF4444")).Bold(true) // Red

	NoteTextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EAB308")). // Yellow
			Italic(true)

	DimTextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")). // Muted grey
			Italic(true)
)
