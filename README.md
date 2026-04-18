# monkeytui

monkeytui is a terminal-based typing test with a clean, focused TUI.

I like terminals and TUIs more than GUIs, and I love typing. This project is my attempt to bring those two together in a UI I actually enjoy using.

## Features

- Fast, minimal TUI
- Timed sessions with adjustable duration
- Optional punctuation mode
- Accuracy and WPM stats
- Keyboard-first controls

## Platforms

Prebuilt binaries are available for macOS, Linux, and Windows.

## Install

### From GitHub Releases

1. Download the latest release for your platform.
2. Extract the archive.
3. Run the binary.

### With curl install scripts

`sh` (macOS/Linux):

```sh
curl -fsSL https://raw.githubusercontent.com/Senasphy/monkeytui/main/scripts/install.sh | sh
```

PowerShell (Windows):

```powershell
curl.exe -fsSL https://raw.githubusercontent.com/Senasphy/monkeytui/main/scripts/install.ps1 | powershell -NoProfile -ExecutionPolicy Bypass -Command -
```

## Usage

Run the binary:

```sh
./monkeytui
```

Controls:

- `ctrl+c` | `ctrl+z`: quit
- `tab` then `enter`: restart
- `/`: open settings
- `j` / `k` or arrow keys: navigate settings
- `enter` or `space`: apply setting

## Configuration

Settings are stored in a TOML file under your user config directory at `monkeytui/config.toml`. You can edit manually, but using the CLI commands below is the easier and recommended way to avoid mistakes.

Commands:

```sh
monkeytui config get duration
monkeytui config set duration 30
monkeytui config get punctuation
monkeytui config set punctuation on
monkeytui reset
```

Keys:

- `duration`: seconds per session
- `punctuation`: `on` or `off`

## Contributing

Issues and pull requests are welcome. Keep changes focused and avoid style-only churn.

## License

See `LICENSE`.
