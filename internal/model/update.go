package model

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"monkeytui/internal/ui"
	"monkeytui/internal/words"
)

func (m Model) Init() tea.Cmd {
	return tea.Batch(loadFile(), tick())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case successMsg:
		m.Words = msg
		if m.SessionDuration == 0 {
			m.SessionDuration = defaultSessionDuration
		}
		m.regenerateSession()
		m.Loaded = true
		return m, nil

	case tickMsg:
		if m.Started && !m.Finished {
			m.Elapsed = time.Since(m.StartTime)
			if m.Elapsed >= m.SessionDuration {
				m.Elapsed = m.SessionDuration
				m.Finished = true
			}
		}
		return m, tick()

	case errMsg:
		m.Err = msg
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case tea.KeyMsg:
		if msg.String() == "/" {
			m.ShowSettings = !m.ShowSettings
			return m, nil
		}

		if m.ShowSettings {
			switch msg.String() {
			case "esc":
				m.ShowSettings = false
			case "up", "k":
				m.SettingsCursor--
				if m.SettingsCursor < 0 {
					m.SettingsCursor = settingsItemCount - 1
				}
			case "down", "j":
				m.SettingsCursor++
				if m.SettingsCursor >= settingsItemCount {
					m.SettingsCursor = 0
				}
			case "enter", " ":
				m.applySettingsSelection()
			case "p":
				m.Punctuation = !m.Punctuation
				m.regenerateSession()
			case "1":
				m.SessionDuration = 15 * time.Second
				m.regenerateSession()
			case "2":
				m.SessionDuration = 30 * time.Second
				m.regenerateSession()
			case "3":
				m.SessionDuration = time.Minute
				m.regenerateSession()
			}
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "tab":
			m.RestartArmed = true
			return m, nil

		case "enter":
			if m.RestartArmed {
				m.regenerateSession()
				return m, tick()
			}
			return m, nil

		case "backspace":
			if m.Finished {
				break
			}
			if m.Cursor == 0 {
				break
			}
			m.Cursor--
			m.Chars[m.Cursor].State = Untouched
			m.Chars[m.Cursor].Typed = 0

		default:
			if m.Finished || m.Cursor >= len(m.Chars) {
				break
			}
			if !m.Started {
				m.Started = true
				m.StartTime = time.Now()
			}
			m.RestartArmed = false

			var typed rune
			if msg.String() == "space" {
				typed = ' '
			} else {
				runes := []rune(msg.String())
				if len(runes) != 1 {
					break
				}
				typed = runes[0]
			}

			m.Chars[m.Cursor].Typed = typed
			m.Attempts++
			if typed == m.Chars[m.Cursor].Expected {
				m.Chars[m.Cursor].State = Correct
				m.CorrectHits++
			} else {
				m.Chars[m.Cursor].State = Wrong
			}
			m.Cursor++
		}
	}

	return m, nil
}

func (m Model) View() tea.View {
	chars := make([]ui.Char, len(m.Chars))
	for i, ch := range m.Chars {
		state := ui.StateUntouched
		switch ch.State {
		case Correct:
			state = ui.StateCorrect
		case Wrong:
			state = ui.StateWrong
		}
		chars[i] = ui.Char{
			Expected: ch.Expected,
			Typed:    ch.Typed,
			State:    state,
		}
	}

	return ui.Render(ui.State{
		Err:             m.Err,
		Loaded:          m.Loaded,
		Width:           m.Width,
		Height:          m.Height,
		Chars:           chars,
		Cursor:          m.Cursor,
		ShowSettings:    m.ShowSettings,
		Punctuation:     m.Punctuation,
		SessionDuration: m.SessionDuration,
		Elapsed:         m.Elapsed,
		Started:         m.Started,
		Finished:        m.Finished,
		SettingsCursor:  m.SettingsCursor,
		WPM:             m.wpm(),
		Accuracy:        m.accuracy(),
	})
}

func loadFile() tea.Cmd {
	return func() tea.Msg {
		loadedWords, err := words.Load(words.DefaultPath)
		if err != nil {
			return errMsg(err)
		}
		return successMsg(loadedWords)
	}
}

func tick() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
