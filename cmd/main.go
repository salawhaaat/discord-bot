package main

import (
	"go-discord-bot/internal/bot"
	"go-discord-bot/internal/bot/discord"
	"go-discord-bot/internal/config"
	"log"
)

func main() {
	cfg := config.New()
	bot := bot.Bot(discord.NewDiscordBot())
	err := bot.Init(*cfg)
	if err != nil {
		log.Fatal(err)
	}
	bot.Run()

}
