package ui

import "charm.land/lipgloss/v2"

var (
	colorBg      = lipgloss.Color("#1a1b2e")
	colorPanel   = lipgloss.Color("#21223a")
	colorPanelHi = lipgloss.Color("#2a2b45")
	colorDim     = lipgloss.Color("#7c7f9e")
	colorSubtle  = lipgloss.Color("#9499b8")
	colorCorrect = lipgloss.Color("#f1f5f9")
	colorWrong   = lipgloss.Color("#fb7185")
	colorCursor  = lipgloss.Color("#38bdf8")
	colorTitle   = lipgloss.Color("#fbbf24")
	colorMeta    = lipgloss.Color("#475569")
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
