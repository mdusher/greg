package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"strings"
	"github.com/bwmarrin/discordgo"
)

var BotToken = os.Getenv("BOT_TOKEN")
var BotPrefix = strings.Split(os.Getenv("BOT_PREFIX"),",")

func main() {
	discord, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		fmt.Println("Error creating Discord session", err)
		return
	}

	// Do stuff with messages
	discord.AddHandler(goGregGo)

	// Do some listening, mate
	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening connection", err)
		return
	}

	fmt.Println("Greg is now Gregging..")
	getAllChannels(discord)

	// Wait for the good old CTRL+C
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	fmt.Println("Greg is going to cease to Greg.")

	// Clean exit
	discord.Close()
}

func goGregGo(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Don't listen to yourself, Greg
	if m.Author.ID == s.State.User.ID {
		return
	}

	var isCommand, RestOfMessage = hasBotPrefix(m.Content, BotPrefix)
	if isCommand {
		s.ChannelMessageSend(m.ChannelID, "Greg recognises your Greggin' and you said: ```"+ RestOfMessage+"```")
	}
}

func hasBotPrefix(m string, prefixes []string) (bool, string) {
	var SplitMessage = strings.SplitN(strings.ToLower(m), " ", 2);
	var FirstToken = SplitMessage[0]
	var RestOfMessage string
	if len(SplitMessage) > 1 {
		RestOfMessage = SplitMessage[1]
	} else {
		RestOfMessage = ""
	}
	if strings.Trim(RestOfMessage, " ") == "" {
		return false, m
	}
	for _, p := range prefixes {
		var LowerPrefix = strings.ToLower(p)
		if strings.HasPrefix(FirstToken, LowerPrefix) || FirstToken == LowerPrefix {
			return true, RestOfMessage
		}
	}
	return false, m
}

func getAllChannels(s *discordgo.Session) {
	fmt.Println("Greg is in:")
	for _, guild := range s.State.Guilds {
		GuildChannels, _ := s.GuildChannels(guild.ID)
		for _, channel := range GuildChannels {
			if channel.Type != discordgo.ChannelTypeGuildText {
				continue
			}
			fmt.Println("- (" + guild.Name + ") " + channel.Name)
		}
	}
}
