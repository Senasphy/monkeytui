package main

import (
	"fmt"
)

func (m Model) wpm() int {
		if !m.Started || m.Elapsed.Seconds() == 0 {
        return 0
    }

	correctChars := 0;
	for _, ch := range m.Chars {
		if ch.State == Correct {
			correctChars++
		}
	}

	minutes := m.Elapsed.Minutes()
	return int(float64((correctChars)/len(m.Chars))/minutes)
}

func (m Model) accuracy() int {
	typed := 0
	correct := 0

	for _, ch := range m.Chars {
		if ch.State != Untouched {
			typed++
		} 
		if ch.State == Correct {
			correct++
		}
	}

	if typed == 0 {
		return 100
	}

	return (correct *100) /typed
}

func (m Model) elapsed() string {
	if !m.Started {
		return fmt.Sprintf("--")
	}

	secs := int(m.Elapsed.Seconds())
	return fmt.Sprintf("%d", secs)

}
