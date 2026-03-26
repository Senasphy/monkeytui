package model

import (
	"math/rand"
	"strings"
	"time"
	"unicode"
)

type CharState int

const (
	Untouched CharState = iota
	Correct
	Wrong
)

const (
	defaultSessionDuration = 30 * time.Second
	wordsPerSession        = 220
)

type Char struct {
	Expected rune
	Typed    rune
	State    CharState
}

type Model struct {
	StartTime       time.Time
	Elapsed         time.Duration
	Started         bool
	Finished        bool
	RestartArmed    bool
	Attempts        int
	CorrectHits     int
	Words           []string
	Punctuation     bool
	SessionDuration time.Duration
	ShowSettings    bool
	SettingsCursor  int
	Chars           []Char
	Cursor          int
	Loaded          bool
	Height          int
	Width           int
	Err             error
	rng             *rand.Rand
}

type errMsg error
type successMsg []string
type tickMsg time.Time

func New(rng *rand.Rand) *Model {
	if rng == nil {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	return &Model{rng: rng}
}

func capitalizeFirstWord(word string) string {
	runes := []rune(word)
	if len(runes) == 0 {
		return word
	}
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func (m *Model) buildSentence() string {
	if len(m.Words) == 0 {
		return ""
	}

	if m.rng == nil {
		m.rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	parts := make([]string, 0, wordsPerSession)
	capitalizeNext := m.Punctuation
	for i := 0; i < wordsPerSession; i++ {
		word := m.Words[m.rng.Intn(len(m.Words))]
		if m.Punctuation {
			if capitalizeNext || m.rng.Intn(100) < 3 {
				word = capitalizeFirstWord(word)
				capitalizeNext = false
			}

			if m.rng.Intn(100) < 10 {
				p := m.rng.Intn(100)
				punctuation := ","
				switch {
				case p < 55:
					punctuation = ","
				case p < 82:
					punctuation = "."
				case p < 91:
					punctuation = "?"
				case p < 97:
					punctuation = "!"
				case p < 99:
					punctuation = ";"
				default:
					punctuation = ":"
				}
				word += punctuation
				if punctuation == "." || punctuation == "?" || punctuation == "!" {
					capitalizeNext = true
				}
			}
		}
		parts = append(parts, word)
	}
	return strings.Join(parts, " ")
}

func (m *Model) setSentence(sentence string) {
	chars := make([]Char, len([]rune(sentence)))
	for i, r := range []rune(sentence) {
		chars[i] = Char{Expected: r, State: Untouched}
	}
	m.Chars = chars
}

func (m *Model) regenerateSession() {
	m.setSentence(m.buildSentence())
	m.resetSession()
}

func (m *Model) resetSession() {
	for i := range m.Chars {
		m.Chars[i].State = Untouched
		m.Chars[i].Typed = 0
	}
	m.StartTime = time.Time{}
	m.Elapsed = 0
	m.Started = false
	m.Finished = false
	m.RestartArmed = false
	m.Attempts = 0
	m.CorrectHits = 0
	m.Cursor = 0
}
