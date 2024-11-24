# Bitpin TUI

A feature-rich Terminal User Interface (TUI) client for the Bitpin cryptocurrency exchange, built with Go. This application provides a comprehensive interface for trading, market data visualization, and account management.

## Features

### Authentication & Session Management
- Secure API key-based authentication
- Automatic token refresh handling
- Session persistence across launches
- Database-backed session storage

### Market Data
- Real-time order book visualization
- Live market ticker data
- Recent trades history
- Favorite symbols tracking
- Configurable update intervals
- Customizable market depth views

### Trading Features
- View account balances with live updates
- Create market and limit orders
- View and manage open orders
- Single-order cancellation
- Batch order cancellation ("Pussy Out" feature)
- Order filtering by symbol and side

### User Interface
- Intuitive terminal-based interface
- Interactive forms and menus
- Color-coded data visualization
- Real-time data updates
- Keyboard shortcuts for navigation
- Mouse interaction support
- Modal-based notifications and confirmations

## Installation

### Prerequisites
- Go 1.21 or higher
- SQLite3
- UPX (optional, for binary compression)

### Building from Source

1. Clone the repository:
```bash
git clone https://github.com/amiwrpremium/bitpin-tui.git
cd bitpin-tui
```

2. Build the application:
```bash
go build -o bitpin-tui
```

### Pre-built Binaries

Pre-built binaries are available for the following platforms:
- macOS (Apple Silicon and Intel)
- Windows (64-bit)
- Linux (64-bit)

Download the latest release from the [Releases](https://github.com/amiwrpremium/bitpin-tui/releases) page.

## Usage

1. Launch the application:
```bash
./bitpin-tui
```

2. Choose how you want to use the application:
    - Access public market data immediately without authentication (order book, tickers, recent trades)
    - Authenticate with your Bitpin API credentials (API Key and Secret) to access private features like trading and balance management

3. Navigate the interface using:
    - Arrow keys for movement
    - Enter to select
    - ESC to return to main menu
    - Ctrl+C to exit

## Main Menu Options

### Authenticated Users
- **Get Balances**: View account balances with live updates
- **Get Open Orders**: View active orders with filtering
- **Create Order**: Place new market or limit orders
- **Cancel Order**: Cancel specific orders by ID
- **Pussy Out**: Batch cancel all open orders
- **Refresh**: Re-authenticate and refresh session
- **Logout**: Clear session and exit

### All Users
- **Get Order Book**: View real-time order book data
- **Get Ticker**: View current market tickers
- **Get Recent Trades**: View trade history
- **Reset**: Reset database to default state
- **Quit**: Exit application

## Project Structure
```
.
├── bitpin_client/     # Bitpin API client implementation
├── db/               # Database operations and models
├── tui/             # Terminal UI components and handlers
├── utils/           # Utility functions
└── scripts/         # Build and release scripts
```

## Development

### Key Components
- **Bitpin Client**: Handles all API interactions with automatic token refresh
- **Database Layer**: Manages session persistence and favorite symbols
- **TUI Components**: Provides interactive UI elements using tview
- **Utils**: Includes JWT handling, string manipulation, and time utilities

### Building for Release
The project includes a release script (`scripts/release.sh`) that:
- Builds binaries for multiple platforms
- Applies build optimizations
- Compresses binaries using UPX
- Prepares releases for distribution

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Notice

This project is an unofficial, independent terminal interface for the Bitpin exchange. It is:
- Not affiliated with, endorsed by, or connected to Bitpin in any way
- Provided "as-is" without any warranty or guarantee of reliability
- A community project created for educational and convenience purposes
- Not responsible for any trading losses or technical issues that may occur

Users should:
- Use this software at their own risk
- Always verify critical information through official channels
- Keep their API credentials secure and never share them
- Be aware that cryptocurrency trading carries inherent risks

## Author

AMiWR ([@amiwrpremium](https://github.com/amiwrpremium))

---

🌟 If you find this project useful, please consider giving it a star on GitHub!
