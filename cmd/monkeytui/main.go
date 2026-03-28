package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	tea "charm.land/bubbletea/v2"
	"monkeytui/internal/config"
	"monkeytui/internal/model"
)

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	if handled, code := handleCLI(); handled {
		os.Exit(code)
	}

	m := model.New(rng)
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	m.ApplyConfig(cfg)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func handleCLI() (bool, int) {
	args := os.Args[1:]
	if len(args) == 0 {
		return false, 0
	}

	switch args[0] {
	case "config":
		return true, handleConfig(args[1:])
	case "reset":
		if err := config.Reset(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return true, 1
		}
		return true, 0
	default:
		fmt.Fprintln(os.Stderr, "unknown command")
		return true, 2
	}
}

func handleConfig(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: monkeytui config [get|set|reset]")
		return 2
	}

	switch args[0] {
	case "get":
		if len(args) != 2 {
			fmt.Fprintln(os.Stderr, "usage: monkeytui config get <key>")
			return 2
		}
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		switch args[1] {
		case "duration":
			fmt.Println(cfg.Duration)
		case "punctuation":
			fmt.Println(cfg.Punctuation)
		default:
			fmt.Fprintln(os.Stderr, "unknown key")
			return 2
		}
		return 0
	case "set":
		if len(args) != 3 {
			fmt.Fprintln(os.Stderr, "usage: monkeytui config set <key> <value>")
			return 2
		}
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		switch args[1] {
		case "duration":
			val, err := strconv.Atoi(args[2])
			if err != nil || val <= 0 {
				fmt.Fprintln(os.Stderr, "duration must be a positive integer")
				return 2
			}
			cfg.Duration = val
		case "punctuation":
			switch args[2] {
			case "true", "false":
				val, _ := strconv.ParseBool(args[2])
				cfg.Punctuation = val
			case "on":
				cfg.Punctuation = true
			case "off":
				cfg.Punctuation = false
			default:
				fmt.Fprintln(os.Stderr, "punctuation must be true|false|on|off")
				return 2
			}
		default:
			fmt.Fprintln(os.Stderr, "unknown key")
			return 2
		}
		if err := config.Save(cfg); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		return 0
	case "reset":
		if err := config.Reset(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		return 0
	default:
		fmt.Fprintln(os.Stderr, "unknown subcommand")
		return 2
	}
}
