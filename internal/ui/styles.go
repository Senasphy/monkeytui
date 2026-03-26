package ui

import "charm.land/lipgloss/v2"

var (
	colorBg      = lipgloss.Color("#0d0d0f")
	colorPanel   = lipgloss.Color("#11111a")
	colorPanelHi = lipgloss.Color("#171726")
	colorDim     = lipgloss.Color("#3a3a4a")
	colorSubtle  = lipgloss.Color("#6c6c8a")
	colorCorrect = lipgloss.Color("#e2e2f0")
	colorWrong   = lipgloss.Color("#ff4d6d")
	colorCursor  = lipgloss.Color("#c084fc")
	colorTitle   = lipgloss.Color("#fef3c7")
	colorMeta    = lipgloss.Color("#4a4a6a")
)

var (
	untouchedStyle = lipgloss.NewStyle().Foreground(colorDim).Background(colorPanel)
	correctStyle   = lipgloss.NewStyle().Foreground(colorCorrect).Background(colorPanel)
	wrongStyle     = lipgloss.NewStyle().Foreground(colorWrong).Underline(true).Background(colorPanel)
	cursorStyle    = lipgloss.NewStyle().Foreground(colorCursor).Underline(true).Bold(true).Background(colorPanel)
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(colorTitle).
			Background(colorPanel).
			Bold(true)

	boxStyle = lipgloss.NewStyle().
			Padding(2, 4).
			Background(colorPanel)

	hintStyle = lipgloss.NewStyle().
			Foreground(colorMeta).
			Background(colorPanel).
			Italic(true).
			PaddingTop(1)

	settingsTitleStyle = lipgloss.NewStyle().
				Foreground(colorTitle).
				Bold(true)

	settingsSelectedStyle = lipgloss.NewStyle().
				Foreground(colorCursor).
				Bold(true)

	settingsItemStyle = lipgloss.NewStyle().
				Foreground(colorCorrect)

	settingsHelpStyle = lipgloss.NewStyle().
				Foreground(colorMeta).
				Italic(true)
)
