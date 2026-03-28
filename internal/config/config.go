package config

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Duration    int  `toml:"duration"`
	Punctuation bool `toml:"punctuation"`
}

func Default() Config {
	return Config{
		Duration:    30,
		Punctuation: false,
	}
}

func Path() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "monkeytui", "config.toml"), nil
}

func Load() (Config, error) {
	path, err := Path()
	if err != nil {
		return Default(), err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Default(), nil
		}
		return Default(), err
	}

	cfg := Default()
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return Default(), err
	}

	if cfg.Duration <= 0 {
		cfg.Duration = Default().Duration
	}

	return cfg, nil
}

func Save(cfg Config) error {
	path, err := Path()
	if err != nil {
		return err
	}

	if cfg.Duration <= 0 {
		cfg.Duration = Default().Duration
	}

	data, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o600)
}

func Reset() error {
	return Save(Default())
}
