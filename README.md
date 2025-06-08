<h1 align="center"> crictty </h1>

<p align="center">
Beautiful, minimal TUI cricket scorecard viewer
</p>

<br>

<div align="center" style="width: 80%; margin: auto;">

<table>
  <tr>
    <th>1st Innings</th>
    <th>2nd Innings</th>
  </tr>
  <tr>
    <td><img src="assets/bat1.png" width="300"></td>
    <td><img src="assets/bat2.png" width="300"></td>
  </tr>
  <tr>
    <td><img src="assets/bowl1.png" width="300"></td>
    <td><img src="assets/bowl2.png" width="300"></td>
  </tr>
</table>

</div>

---

## Features

- **Live Cricket Scores:** Real-time updates from Cricbuzz
- **Match Details:** Team scores, current batsmen, bowler figures
- **Complete Scorecards:** Detailed batting and bowling statistics
- **Innings Navigation:** Browse through all innings with ease
- **Multi-Match Support:** Switch between multiple live matches
- **Clean Interface:** Minimal, terminal-friendly design

## Installation

### Using `go install`

```bash
go install github.com/ashish0kumar/crictty@latest
```

### From Source

```bash
git clone https://github.com/ashish0kumar/crictty.git
cd crictty
go build
sudo mv crictty /usr/local/bin/
crictty -h
```

## Usage

```bash
# View all live matches
crictty

# View a specific match
crictty --match-id 118928

# Set refresh rate to 30 seconds
crictty --tick-rate 30000

# Show help
crictty --help
```

> [!TIP]
> To use the `--match-id` flag, open the specific match page on [Cricbuzz](https://www.cricbuzz.com), and extract the match ID from the URL <br>
`https://www.cricbuzz.com/live-cricket-scorecard/<id>/...`

### Controls

| Key | Action |
|-----|--------|
| **`←`** **`→`** | Switch between matches |
| **`↑`** **`↓`** | Navigate innings |
| **`b`** | Toggle batting/bowling view |
| **`q`** | Quit application |

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling and layout
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [goquery](https://github.com/PuerkitoBio/goquery) - HTML parsing

## Acknowledgments

- [Cricbuzz](https://www.cricbuzz.com) for providing cricket data
- [Charm](https://charm.sh/) for the excellent TUI libraries

## Contributing

Contributions are always welcome! Feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT License](LICENSE)

<br><br>

<p align="center">
	<img src="https://raw.githubusercontent.com/catppuccin/catppuccin/main/assets/footers/gray0_ctp_on_line.svg?sanitize=true" />
</p>
