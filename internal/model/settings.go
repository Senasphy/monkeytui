package model

import "time"

const settingsItemCount = 4

func (m *Model) applySettingsSelection() {
	switch m.SettingsCursor {
	case 0:
		m.setPunctuation(!m.Punctuation)
	case 1:
		m.setSessionDuration(15 * time.Second)
	case 2:
		m.setSessionDuration(30 * time.Second)
	case 3:
		m.setSessionDuration(time.Minute)
	}
}

func (m *Model) setPunctuation(value bool) {
	m.Punctuation = value
	m.regenerateSession()
	m.persistConfig()
}

func (m *Model) setSessionDuration(value time.Duration) {
	m.SessionDuration = value
	m.regenerateSession()
	m.persistConfig()
}
