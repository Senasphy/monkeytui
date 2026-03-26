package main

import (
	"log"
	"math/rand"
	"time"

	tea "charm.land/bubbletea/v2"
	"monkeytui/internal/model"
)

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	p := tea.NewProgram(model.New(rng))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
