# Discord Bot with Tic-Tac-Toe and Weather Commands

This Go-based Discord bot offers Tic-Tac-Toe and Weather commands. Below is a brief overview of the project.

## Project Structure

- **cmd**: Main entry point of the application.
- **config**: Manages configuration settings.
- **internal/bot**: Implements core bot functionality.
  - **internal/bot/discord**: Handles Discord-specific features.
- **pkg/api/tictactoe**: Provides Tic-Tac-Toe game logic.
- **pkg/api/weather**: Interfaces with an external weather API.

## Getting Started

1. **Clone and Build:**

   ```bash
   git clone <repository-url>
   cd <project-directory>
   go build ./cmd -config=config/local.toml
   ```

### Available Commands

```bash
    /weather <city>
    /ttt star # Start a new game.
    /ttt move <position> # Make a move.
    /help
```
