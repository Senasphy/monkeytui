package model

import (
	"fmt"
)

func (m Model) wpm() int {
	if !m.Started || m.Elapsed.Seconds() == 0 {
		return 0
	}

	correctChars := 0
	for _, ch := range m.Chars {
		if ch.State == Correct {
			correctChars++
		}
	}

	minutes := m.Elapsed.Minutes()
	if minutes == 0 {
		return 0
	}

	return int((float64(correctChars) / 5.0) / minutes)
}

func (m Model) accuracy() int {
	if m.Attempts == 0 {
		return 100
	}

	return (m.CorrectHits * 100) / m.Attempts
}

func (m Model) elapsed() string {
	if !m.Started {
		return "--"
	}

	secs := int(m.Elapsed.Seconds())
	return fmt.Sprintf("%d", secs)

}
