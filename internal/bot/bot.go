package bot

import "go-discord-bot/internal/config"

type Bot interface {
	Run()
	Init(cfg config.Config) error
	New() Bot
}
