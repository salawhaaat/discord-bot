package discord

import (
	"fmt"
	"go-discord-bot/internal/bot"
	"go-discord-bot/internal/config"
	"go-discord-bot/pkg/api/tictactoe"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type DiscordBot struct {
	sess *discordgo.Session
	cfg  config.Config
	tk   *tictactoe.TicTacToe
}

func NewDiscordBot() bot.Bot {
	return &DiscordBot{}
}

func (d *DiscordBot) New() bot.Bot {
	return NewDiscordBot()
}

func (d *DiscordBot) Run() {
	err := d.sess.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer d.sess.Close()
	fmt.Println("The bot is online!")

	sc := make(chan os.Signal, 1) // ctrl c
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func (d *DiscordBot) Init(config config.Config) error {
	d.cfg = config
	sess, err := discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		log.Fatal(err)
	}
	d.sess = sess
	d.sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
	d.sess.AddHandler(d.messageCreate)

	return nil
}
