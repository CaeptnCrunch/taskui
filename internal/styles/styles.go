package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	FocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	BlurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("246"))
	SubtleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("242"))
	NoStyle      = lipgloss.NewStyle()
	Prompt       = "❯ "
	AntiPrompt   = " ❮"
	Dot          = " • "
)
