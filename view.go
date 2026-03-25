package main

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	colorBg        = lipgloss.Color("#0d0d0f")
	colorDim       = lipgloss.Color("#3a3a4a")
	colorSubtle    = lipgloss.Color("#6c6c8a")
	colorCorrect   = lipgloss.Color("#e2e2f0")
	colorWrong     = lipgloss.Color("#ff4d6d")
	colorCursor    = lipgloss.Color("#c084fc")
	colorAccent = lipgloss.Color("#f59e0b")
	colorTitle  = lipgloss.Color("#fef3c7")
	colorMeta      = lipgloss.Color("#4a4a6a")
)

var (
	untouchedStyle = lipgloss.NewStyle().Foreground(colorDim)
	correctStyle   = lipgloss.NewStyle().Foreground(colorCorrect)
	wrongStyle     = lipgloss.NewStyle().Foreground(colorWrong).Underline(true)
	cursorStyle    = lipgloss.NewStyle().Foreground(colorCursor).Underline(true).Bold(true)
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(colorTitle).
			Bold(true)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(colorSubtle).
			Italic(true)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorAccent).
			Padding(2, 4).
			Background(lipgloss.Color("#11111a"))

	hintStyle = lipgloss.NewStyle().
			Foreground(colorMeta).
			Italic(true).
			PaddingTop(1)

	statusCorrectStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#a78bfa")).
				Bold(true)

	statusWrongStyle = lipgloss.NewStyle().
				Foreground(colorWrong).
				Bold(true)
)

func (m Model) renderStatus() string {
    correct := 0
    wrong := 0
    for _, ch := range m.Chars {
        switch ch.State {
        case Correct:
            correct++
        case Wrong:
            wrong++
        }
    }

    timerDisplay := ""
    if !m.Started {
				timerDisplay = lipgloss.NewStyle().Foreground(colorDim).Render("●")
    } else {
        timerDisplay = lipgloss.NewStyle().Foreground(colorSubtle).Render(m.elapsed())
    }

    statusBar := fmt.Sprintf("%s   %s %s   %s %s",
        timerDisplay,
        statusCorrectStyle.Render(fmt.Sprintf("✓ %d", correct)),
        lipgloss.NewStyle().Foreground(colorMeta).Render("correct"),
        statusWrongStyle.Render(fmt.Sprintf("✗ %d", wrong)),
        lipgloss.NewStyle().Foreground(colorMeta).Render("wrong"),
    )

    if m.Finished {
        statusBar += fmt.Sprintf("   WPM: %s   Accuracy: %s",
            lipgloss.NewStyle().Foreground(colorCursor).Bold(true).Render(fmt.Sprintf("%d", m.wpm())),
            lipgloss.NewStyle().Foreground(colorSubtle).Render(fmt.Sprintf("%d%%", m.accuracy())),
        )
    }

    return statusBar
}

func (m Model) renderChars() string {
	var sb strings.Builder
	for i, ch := range m.Chars {
		var display string
		if ch.State == Wrong {
			display = string(ch.Typed)
		} else {
			display = string(ch.Expected)
		}

		switch ch.State {
		case Correct:
			sb.WriteString(correctStyle.Render(display))
		case Wrong:
			sb.WriteString(wrongStyle.Render(display))
		default:
			if i == m.Cursor {
				sb.WriteString(cursorStyle.Render(display))
			} else {
				sb.WriteString(untouchedStyle.Render(display))
			}
		}
	}
	return sb.String()
}

func (m Model) View() tea.View {
	var view tea.View

	if m.Err != nil {
		view.SetContent(fmt.Sprintf("error: %s\n", m.Err))
		return view
	}
	if !m.Loaded {
		view.SetContent(
			lipgloss.NewStyle().Foreground(colorSubtle).Italic(true).Render("loading..."),
		)
		return view
	}

	innerWidth := m.Width - 16 
	if innerWidth < 20 {
		innerWidth = 20
	}

	charsOutput := lipgloss.NewStyle().Width(innerWidth).Render(m.renderChars())

	box := boxStyle.Width(innerWidth + 8).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			titleStyle.Render("monkeytui"),
			subtitleStyle.Render("type the text below"),
			lipgloss.NewStyle().PaddingTop(1).Render(""),
			charsOutput,
			lipgloss.NewStyle().PaddingTop(1).Render(""),
			m.renderStatus(),
			hintStyle.Render("ctrl+c to quit"),
		),
	)

	centered := lipgloss.NewStyle().
		Width(m.Width).
		Height(m.Height).
		Align(lipgloss.Center, lipgloss.Center).
		Background(colorBg).
		Render(box)

	view.SetContent(centered)
	return view
}
