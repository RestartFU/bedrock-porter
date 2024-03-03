package frontend

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

var (
	// ProgressBar is a progress bar.
	ProgressBar = progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(100),
	)
	// Style is the default style for the spinner and rendered text.
	Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	// Spinner is a spinner.
	Spinner = spinner.New(
		spinner.WithStyle(Style),
	)
)
