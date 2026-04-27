# Monopoly Banker App

A terminal-based application to manage player balances and transactions in a Monopoly game, featuring a styled UI with ANSI colors and box-drawing borders.

## Features

- Manage up to 6 players
- Track player balances (starting with 15,000,000)
- Styled terminal UI with colors, box-drawing tables, and formatted numbers
- Inline amount entry — type `15M`, `500K`, or plain numbers
- Transaction history — last 5 actions shown on the main screen
- Undo — reverse the last transaction instantly
- Handle various Monopoly transactions:
  - Pay to the Bank (with optional mortgage tax)
  - Transfer between players (rent, trades, etc.)
  - Receive money from the Bank
  - Pass GO (receive 2,000,000)

## Requirements

- Go 1.26.2 or later
- A terminal with ANSI color support

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/moohalem/monopoly.git
   cd monopoly
   ```

2. Build the application:
   ```bash
   go build -o monopoly ./cmd/monopoly
   ```

   Or install it directly:
   ```bash
   go install ./cmd/monopoly
   ```

## Usage

Run the application:
```bash
./monopoly
```

Or if installed via `go install`:
```bash
monopoly
```

Follow the on-screen prompts to:
1. Enter the number of players (1–6)
2. Enter player names
3. Use the menu to perform transactions

## Menu Options

| Key | Action              | Description                                              |
|-----|---------------------|----------------------------------------------------------|
| `1` | Pay to Bank         | Deduct from a player's balance (optionally as mortgage)  |
| `2` | Transfer to Player  | Move money from one player to another                    |
| `3` | Receive from Bank   | Add money to a player's balance                          |
| `G` | Pass GO             | Add 2,000,000 to a player's balance                      |
| `U` | Undo Last           | Reverse the most recent transaction                      |
| `0` | Exit                | Quit the application                                     |

### Amount Entry

When prompted for an amount, you can use shorthand units:

| Input     | Interpreted As |
|-----------|----------------|
| `15M`     | 15,000,000     |
| `500K`    | 500,000        |
| `250000`  | 250,000        |

## Project Structure

```
monopoly/
├── cmd/monopoly/
│   └── main.go            # Application entry point
├── internal/
│   ├── game/
│   │   ├── initialize.go  # Welcome screen and player setup
│   │   └── player.go      # Game logic, menu, and transactions
│   └── ui/
│       ├── style.go       # ANSI colors, banner, and formatting helpers
│       └── terminal.go    # Terminal I/O (input scanner, clear screen)
├── go.mod                 # Go module file
├── LICENSE
└── README.md
```

## Contributing

Feel free to submit issues and pull requests.

## License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.