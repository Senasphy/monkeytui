package model

import (
	"time"

	"monkeytui/internal/config"
)

func (m *Model) ApplyConfig(cfg config.Config) {
	if cfg.Duration > 0 {
		m.SessionDuration = time.Duration(cfg.Duration) * time.Second
	}
	m.Punctuation = cfg.Punctuation
}

func (m *Model) persistConfig() {
	_ = config.Save(config.Config{
		Duration:    int(m.SessionDuration.Seconds()),
		Punctuation: m.Punctuation,
	})
}
