# Monopoly Banker App

A terminal-based application to manage player balances and transactions in a Monopoly game.

## Features

- Manage up to 6 players
- Track player balances (starting with 15,000,000)
- Handle various Monopoly transactions:
  - Buy/Pay to the Bank
  - Pay Rent (between players)
  - Receive money
  - Pay Mortgage (with 10% tax, rounded up to nearest 10,000)
  - Pass GO (receive 2,000,000)

## Requirements

- Go 1.26.2 or later

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/monopoly.git
   cd monopoly
   ```

2. Build the application:
   ```bash
   go build -o banker ./cmd/monopoly
   ```

## Usage

Run the application:
```bash
./banker
```

Follow the on-screen prompts to:
1. Enter the number of players (1-6)
2. Enter player names
3. Use the menu to perform transactions

## Menu Options

1. **Buy / Pay to The Bank**: Deduct amount from a player's balance
2. **Pay Rent**: Transfer money from one player to another
3. **Receive**: Add money to a player's balance
4. **Paying Mortgage**: Pay mortgage with 10% tax (rounded up)
5. **Passing GO**: Add 2,000,000 to a player's balance
0. **Exit**: Quit the application

## Project Structure

```
monopoly/
├── cmd/monopoly/
│   └── main.go          # Application entry point
├── internal/
│   ├── game/
│   │   ├── input.go     # Player setup and input validation
│   │   └── player.go    # Game logic and player management
│   └── ui/
│       └── terminal.go  # Terminal UI helpers
├── go.mod               # Go module file
└── README.md            # This file
```

## Contributing

Feel free to submit issues and pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.