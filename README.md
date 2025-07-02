# OSS Watch 📊

[![Test](https://github.com/jackchuka/gh-oss-watch/workflows/Test/badge.svg)](https://github.com/jackchuka/gh-oss-watch/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/jackchuka/gh-oss-watch)](https://goreportcard.com/report/github.com/jackchuka/gh-oss-watch)

A GitHub CLI plugin that helps open-source maintainers stay on top of repository activity across multiple projects. Track stars, issues, pull requests, and more — all from your terminal.

## Features

- 🔍 **Multi-repo tracking** - Monitor multiple repositories from a single dashboard
- 📊 **Activity monitoring** - Track stars, issues, PRs, forks, etc...
- 🚀 **Diff-based status** - Only see _new_ activity since your last check
- 🎯 **Configurable events** - Choose which events to track per repository
- 📱 **Clean terminal UI** - Beautiful CLI output with emojis and clear formatting

## Installation

### Prerequisites

- [GitHub CLI](https://github.com/cli/cli) installed and authenticated
- Go 1.24+ (if building from source)

### GitHub CLI Extension

```bash
gh extension install jackchuka/gh-oss-watch
```

### Install from Source

```bash
git clone https://github.com/jackchuka/gh-oss-watch.git
cd gh-oss-watch
go build -o gh-oss-watch .
# Copy to your PATH or use directly
```

## Quick Start

1. **Initialize configuration:**

   ```bash
   gh oss-watch init
   ```

2. **Add repositories to watch:**

   ```bash
   gh oss-watch add facebook/react
   gh oss-watch add microsoft/vscode stars issues
   ```

3. **Check for new activity:**

   ```bash
   gh oss-watch status
   ```

4. **View dashboard:**
   ```bash
   gh oss-watch dashboard
   ```

## Commands

| Command                  | Description                        | Example                                    |
| ------------------------ | ---------------------------------- | ------------------------------------------ |
| `init`                   | Initialize config file             | `gh oss-watch init`                        |
| `add <repo> [events...]` | Add repo to watch list             | `gh oss-watch add owner/repo stars issues` |
| `set <repo> <events...>` | Configure events for repo          | `gh oss-watch set owner/repo forks`        |
| `remove <repo>`          | Remove repo from watch list        | `gh oss-watch remove owner/repo`           |
| `status`                 | Show new activity since last check | `gh oss-watch status`                      |
| `dashboard`              | Display summary across all repos   | `gh oss-watch dashboard`                   |

## Event Types

- **`stars`** - Repository stars
- **`issues`** - Issues created/reopened
- **`pull_requests`** - Pull requests opened
- **`forks`** - Repository forks

## Configuration

Configuration is stored in `~/.gh-oss-watch/config.yaml`:

```yaml
repos:
  - repo: facebook/react
    events:
      - stars
      - issues
      - pull_requests
  - repo: microsoft/vscode
    events:
      - stars
      - forks
```

## Example Output

### Status Command

```bash
$ gh oss-watch status

📈 facebook/react:
  ⭐ +23 stars (219,432 total)
  🐛 +5 issues (823 open)
  🔀 +12 pull requests (156 open)

📈 microsoft/vscode:
  ⭐ +45 stars (158,234 total)
  🍴 +8 forks (26,789 total)
```

### Dashboard Command

```bash
$ gh oss-watch dashboard

📊 OSS Watch Dashboard
======================

📁 facebook/react
   ⭐ Stars: 219,432
   🐛 Issues: 823
   🔀 Pull Requests: 156
   🍴 Forks: 43,234
   📅 Last Updated: 2024-01-15 14:23
   📢 Watching: stars, issues, pull_requests

📈 Total Across All Repos:
   ⭐ Total Stars: 377,666
   🐛 Total Issues: 1,456
   🔀 Total PRs: 289
   🍴 Total Forks: 70,023
```

## Automation

Perfect for cron jobs or CI/CD pipelines:

```bash
# Check for activity every hour
0 * * * * /path/to/gh-oss-watch status

# Weekly dashboard summary
0 9 * * 1 /path/to/gh-oss-watch dashboard | mail -s "Weekly OSS Summary" you@example.com
```

## Development

### Project Structure

```
.
├── cmd/                 # Command handlers
├── services/           # Business logic & interfaces
│   └── mock/          # Generated mocks
├── .github/workflows/ # CI/CD pipelines
└── main.go           # Entry point
```

### Building

```bash
go build -o gh-oss-watch .
```

### Testing

```bash
# Run tests
go test ./...

# Run tests with race detection
go test -race ./...

# Generate mocks
cd services && go generate
```

### Code Quality

```bash
# Format code
gofmt -s -w .

# Run linter
golangci-lint run

# Check formatting
gofmt -s -l .
```

## Architecture

- **Modular design** with clean separation of concerns
- **Dependency injection** for testability
- **Generated mocks** using mockgen for comprehensive testing
- **Interface-based architecture** enabling easy mocking and testing
- **Caching system** for efficient API usage and offline support

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for your changes
5. Ensure all tests pass (`go test ./...`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## License

This project is licensed under the MIT License.

## Acknowledgments

- Built with [GitHub CLI](https://github.com/cli/cli) for seamless GitHub integration
- Uses [gomock](https://github.com/golang/mock) for testing
- Inspired by the need for better OSS project monitoring tools

---

**Happy monitoring!** 🎉 If you find this tool useful, please consider giving it a ⭐ on GitHub.
