package ui

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	StateUntouched = iota
	StateCorrect
	StateWrong
)

type Char struct {
	Expected rune
	Typed    rune
	State    int
}

type State struct {
	Err             error
	Loaded          bool
	Width           int
	Height          int
	Chars           []Char
	Cursor          int
	ShowSettings    bool
	Punctuation     bool
	SessionDuration time.Duration
	Elapsed         time.Duration
	Started         bool
	Finished        bool
	SettingsCursor  int
	WPM             int
	Accuracy        int
}

func renderStatus(state State) string {
	timerDisplay := lipgloss.NewStyle().Foreground(colorDim).Background(colorPanel).Render("●")
	remaining := int((state.SessionDuration - state.Elapsed).Seconds())
	if remaining < 0 {
		remaining = 0
	}
	if state.Started || state.Finished {
		timerDisplay = lipgloss.NewStyle().Foreground(colorSubtle).Background(colorPanel).Render(fmt.Sprintf("%d", remaining))
	}

	statusBar := timerDisplay
	if state.Finished {
		statusBar += fmt.Sprintf("   WPM: %s   Accuracy: %s",
			lipgloss.NewStyle().Foreground(colorCursor).Background(colorPanel).Bold(true).Render(fmt.Sprintf("%d", state.WPM)),
			lipgloss.NewStyle().Foreground(colorSubtle).Background(colorPanel).Render(fmt.Sprintf("%d%%", state.Accuracy)),
		)
	}

	return statusBar
}

func renderSettingsPopup(state State) string {
	const settingsPanelWidth = 66

	punctuation := "off"
	if state.Punctuation {
		punctuation = "on"
	}

	timer15 := " "
	timer30 := " "
	timer60 := " "
	switch state.SessionDuration {
	case 15 * time.Second:
		timer15 = "x"
	case time.Minute:
		timer60 = "x"
	default:
		timer30 = "x"
	}

	items := []string{
		fmt.Sprintf("punctuation: %s", punctuation),
		fmt.Sprintf("timer: [%-1s] 15 seconds", timer15),
		fmt.Sprintf("timer: [%-1s] 30 seconds", timer30),
		fmt.Sprintf("timer: [%-1s] 1 minute", timer60),
	}

	panelStyle := lipgloss.NewStyle().
		Width(settingsPanelWidth).
		Padding(1, 2).
		Background(colorPanelHi)
	settingsContentWidth := settingsPanelWidth - panelStyle.GetHorizontalFrameSize()

	lineStyle := lipgloss.NewStyle().
		Width(settingsContentWidth).
		MaxWidth(settingsContentWidth).
		Background(colorPanelHi)

	lines := make([]string, 0, len(items)+3)
	lines = append(lines, lineStyle.Render(settingsTitleStyle.Render("settings")))
	lines = append(lines, lineStyle.Render(lipgloss.NewStyle().Foreground(colorSubtle).Render(strings.Repeat("─", settingsContentWidth-2))))
	lines = append(lines, lineStyle.Render(settingsHelpStyle.Render("up/down or j/k to navigate, enter/space to apply, esc to close")))
	for i, item := range items {
		prefix := "  "
		style := settingsItemStyle
		if i == state.SettingsCursor {
			prefix = "› "
			style = settingsSelectedStyle.Background(colorPanelHi)
		}
		lines = append(lines, lineStyle.Render(style.Render(prefix+item)))
	}

	return panelStyle.Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
}

func wrappedLineRanges(chars []Char, maxWidth int) [][2]int {
	if maxWidth < 1 || len(chars) == 0 {
		return [][2]int{{0, len(chars)}}
	}

	lines := make([][2]int, 0)
	for i := 0; i < len(chars); {
		for i < len(chars) && chars[i].Expected == ' ' {
			i++
		}
		if i >= len(chars) {
			break
		}

		start := i
		width := 0
		lastSpace := -1
		wrapped := false

		for i < len(chars) {
			r := chars[i].Expected
			if r == '\n' {
				lines = append(lines, [2]int{start, i})
				i++
				wrapped = true
				break
			}

			if width == maxWidth {
				if lastSpace >= start {
					lines = append(lines, [2]int{start, lastSpace + 1})
					i = lastSpace + 1
				} else {
					lines = append(lines, [2]int{start, i})
				}
				wrapped = true
				break
			}

			if r == ' ' {
				lastSpace = i
			}

			width++
			i++
		}

		if !wrapped {
			lines = append(lines, [2]int{start, i})
		}
	}

	return lines
}

func currentLineIndex(chars []Char, lines [][2]int, cursor int) int {
	if len(lines) == 0 {
		return 0
	}

	if cursor >= len(chars) && len(chars) > 0 {
		cursor = len(chars) - 1
	}

	for i, line := range lines {
		if cursor >= line[0] && cursor < line[1] {
			return i
		}
	}

	return len(lines) - 1
}

func renderLine(chars []Char, cursor, start, end int) string {
	var sb strings.Builder
	for i := start; i < end; i++ {
		ch := chars[i]
		display := string(ch.Expected)
		if ch.State == StateWrong {
			display = string(ch.Typed)
		}

		switch ch.State {
		case StateCorrect:
			sb.WriteString(correctStyle.Render(display))
		case StateWrong:
			sb.WriteString(wrongStyle.Render(display))
		default:
			if i == cursor {
				sb.WriteString(cursorStyle.Render(display))
			} else {
				sb.WriteString(untouchedStyle.Render(display))
			}
		}
	}
	return sb.String()
}

func renderChars(chars []Char, cursor, maxWidth int) string {
	lines := wrappedLineRanges(chars, maxWidth)
	currentLine := currentLineIndex(chars, lines, cursor)

	startLine := 0
	if currentLine >= 2 {
		startLine = currentLine - 1
	}

	endLine := startLine + 3
	if endLine > len(lines) {
		endLine = len(lines)
	}

	output := make([]string, 0, 3)
	for i := startLine; i < endLine; i++ {
		output = append(output, renderLine(chars, cursor, lines[i][0], lines[i][1]))
	}
	for len(output) < 3 {
		output = append(output, "")
	}

	return strings.Join(output, "\n")
}

func Render(state State) tea.View {
	var view tea.View
	bgColor := colorBg

	if state.Err != nil {
		view.SetContent(fmt.Sprintf("error: %s\n", state.Err))
		return view
	}
	if !state.Loaded {
		view.SetContent(
			lipgloss.NewStyle().Foreground(colorSubtle).Italic(true).Render("loading...\n"),
		)
		return view
	}

	innerWidth := state.Width - boxStyle.GetHorizontalFrameSize() - 4
	if innerWidth < 20 {
		innerWidth = 20
	}

	charsOutput := renderChars(state.Chars, state.Cursor, innerWidth)
	panelLine := lipgloss.NewStyle().Width(innerWidth).Background(colorPanel)
	charLines := strings.Split(charsOutput, "\n")
	for i, line := range charLines {
		charLines[i] = panelLine.Render(line)
	}
	charsOutput = strings.Join(charLines, "\n")
	statusOutput := panelLine.Render(renderStatus(state))
	hintOutput := panelLine.Render(hintStyle.Render("ctrl+c to quit   tab+enter to restart   / settings"))

	box := boxStyle.Width(innerWidth + boxStyle.GetHorizontalFrameSize()).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			panelLine.Render(titleStyle.Render("monkeytui")),
			panelLine.Render(""),
			charsOutput,
			panelLine.Render(""),
			statusOutput,
			hintOutput,
		),
	)

	if state.ShowSettings {
		popup := renderSettingsPopup(state)
		popupRow := lipgloss.NewStyle().
			Width(lipgloss.Width(box)).
			Align(lipgloss.Center, lipgloss.Top).
			Background(bgColor).
			Render(popup)
		box = lipgloss.JoinVertical(lipgloss.Left, box, popupRow)
	}

	centered := lipgloss.NewStyle().
		Width(state.Width).
		Height(state.Height).
		Align(lipgloss.Center, lipgloss.Center).
		Background(bgColor).
		Render(box)

	view.SetContent(centered)
	view.AltScreen = true
	return view
}
