package model

import "time"

const settingsItemCount = 4

func (m *Model) applySettingsSelection() {
	switch m.SettingsCursor {
	case 0:
		m.Punctuation = !m.Punctuation
		m.regenerateSession()
	case 1:
		m.SessionDuration = 15 * time.Second
		m.regenerateSession()
	case 2:
		m.SessionDuration = 30 * time.Second
		m.regenerateSession()
	case 3:
		m.SessionDuration = time.Minute
		m.regenerateSession()
	}
}
