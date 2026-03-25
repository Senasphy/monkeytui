package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"
	tea "charm.land/bubbletea/v2"
)

type CharState int

const (
	Untouched CharState = iota
	Correct
	Wrong
)

type Char struct {
	Expected rune
	Typed    rune
	State    CharState
}

type Model struct {
	StartTime time.Time
	Elapsed time.Duration
	Started bool
	Finished bool
	Chars  []Char
	Cursor int
	Loaded bool
	Height int 
	Width int
	Err    error
}

type WordList struct {
	Language string   `json:"language"`
	Words    []string `json:"words"`
}


type errMsg error
type successMsg string
type tickMsg time.Time

func (m Model) Init() tea.Cmd {
	return tea.Batch(loadFile(), tick())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case successMsg:
		sentence := string(msg)
		chars := make([]Char, len([]rune(sentence)))
		for i, r := range []rune(sentence) {
			chars[i] = Char{Expected: r, State: Untouched}
		}
		m.Chars = chars
		m.Loaded = true
		return m, nil

	case tickMsg:
		if !m.Finished {
			if m.Started {
				m.Elapsed = time.Since(m.StartTime)
			} 
			return m, tick()
		}
		return m, nil

	case errMsg:
		m.Err = msg
		return m, tea.Quit
	

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "backspace":
			if m.Cursor == 0 {
				break
			}
			m.Cursor--
			m.Chars[m.Cursor].State = Untouched
			m.Chars[m.Cursor].Typed = 0

		default:
			if m.Cursor >= len(m.Chars) {
				break
			}
			typed := []rune(msg.String())[0]
			if msg.String() == "space" {
				typed = ' '
			}
			m.Chars[m.Cursor].Typed = typed
			if typed == m.Chars[m.Cursor].Expected {
				m.Chars[m.Cursor].State = Correct
			} else {
				m.Chars[m.Cursor].State = Wrong
			}
			m.Cursor++
		}
	}

	return m, nil
}

// The view method is laid out in view.go

func loadFile() tea.Cmd {
	return func() tea.Msg {
		jsonFile, err := os.ReadFile("./words.json")
		if err != nil {
			return errMsg(err)
		}
		allWords := WordList{}
		err = json.Unmarshal(jsonFile, &allWords)
		if err != nil {
			return errMsg(err)
		}
		sentence := strings.Join(allWords.Words, " ")
		return successMsg(sentence)
	}
}

func tick() tea.Cmd {
	return tea.Every(time.Second, func (t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func main() {
	p := tea.NewProgram(&Model{})
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
