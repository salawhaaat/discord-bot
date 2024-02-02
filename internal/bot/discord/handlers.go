package discord

import (
	"fmt"
	"go-discord-bot/pkg/api/tictactoe"
	"go-discord-bot/pkg/api/weather"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (d *DiscordBot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || !strings.HasPrefix(m.Content, "/") {
		return
	}

	parts := strings.Fields(m.Content)
	if len(parts) < 1 {
		s.ChannelMessageSend(m.ChannelID, "Invalid command. Commands must have at least one word.")
		return
	}

	command := parts[0]

	switch command {
	case "/weather":
		go d.weather(s, m)
	case "/help":
		go d.help(s, m)
	case "/ttt":
		go d.playtictactoe(s, m)
	default:
		s.ChannelMessageSend(m.ChannelID, "Unknown command. Available commands: `/weather`, `/help`, `/ttt`.")
	}
}

func (d *DiscordBot) help(s *discordgo.Session, m *discordgo.MessageCreate) {
	helpMessage := "**Available Commands:**\n" +
		"- `/weather <city>`: Get weather information for the specified city.\n" +
		"- `/ttt start`: Start a new Tic-Tac-Toe game.\n" +
		"- `/ttt move <position>`: Make a move in the Tic-Tac-Toe game.\n" +
		"- `/help`: Display this help message."
	s.ChannelMessageSend(m.ChannelID, helpMessage)
}

func (d *DiscordBot) weather(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Check if the message starts with "/weather"
	if m.Author.ID == s.State.User.ID || !strings.HasPrefix(m.Content, "/weather") {
		return
	}

	parts := strings.Fields(m.Content)
	if len(parts) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Please provide a city after `/weather`.")
		return
	}

	city := strings.Join(parts[1:], " ")

	// Get weather information
	result, err := weather.GetWeather(d.cfg.WeatherApikey, city)
	if err != nil {
		fmt.Println("Error getting weather:", err)
		s.ChannelMessageSend(m.ChannelID, "Error getting weather information.")
		return
	}
	embed := &discordgo.MessageEmbed{
		Title:       "Weather Information",
		Description: result,
		Color:       0x3498db,
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func (d *DiscordBot) playtictactoe(s *discordgo.Session, m *discordgo.MessageCreate) {
	if d.tk == nil {
		d.tk = tictactoe.New()
	}

	if d.tk.IsGameOver {
		s.ChannelMessageSend(m.ChannelID, "Game over. Start a new game with `/ttt start`.")
		return
	}

	d.tk.Mutex.Lock()
	defer d.tk.Mutex.Unlock()

	fields := strings.Fields(m.Content)
	if len(fields) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Invalid command. Use `/ttt start` to begin a new game.")
		return
	}

	command := fields[1]

	switch command {
	case "start":
		d.tk = tictactoe.New()
		s.ChannelMessageSend(m.ChannelID, "New game started!\n"+d.tk.PrintBoard())
	case "move":
		if len(fields) < 3 {
			s.ChannelMessageSend(m.ChannelID, "Invalid move. Use `/ttt move <position>`.")
			return
		}
		position := fields[2]
		if !(position >= "1" && position <= "9") {
			s.ChannelMessageSend(m.ChannelID, "Invalid position. Please choose a number between 1 and 9.")
			return
		}

		if !d.tk.MakeMove(position) {
			s.ChannelMessageSend(m.ChannelID, "Invalid move. Position already taken. Try again.")
			return
		}

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Move made by %s!\n%s", d.tk.CurrentTurn, d.tk.PrintBoard()))

		if d.tk.CheckWin() {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Player %s wins!", d.tk.CurrentTurn))
			d.tk.IsGameOver = true
			return
		}

		if d.tk.CheckDraw() {
			s.ChannelMessageSend(m.ChannelID, "It's a draw! The game is over.")
			d.tk.IsGameOver = true
			return
		}

		d.tk.SwitchTurn()

		// Bot's Move
		botMove := d.tk.GetRandBotMove()
		if !d.tk.MakeMove(botMove) {
			s.ChannelMessageSend(m.ChannelID, "Error making a move for the bot.")
			return
		}
		d.tk.SwitchTurn()

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Bot's move:\n%s", d.tk.PrintBoard()))

		if d.tk.CheckWin() {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Bot wins!"))
			d.tk.IsGameOver = true
			return
		}

		if d.tk.CheckDraw() {
			s.ChannelMessageSend(m.ChannelID, "It's a draw! The game is over.")
			d.tk.IsGameOver = true
			return
		}
	default:
		s.ChannelMessageSend(m.ChannelID, "Invalid command. Use `/ttt start` to begin a new game or `/ttt move <position>` to make a move.")
	}
}
