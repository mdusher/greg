package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"strings"
	"github.com/bwmarrin/discordgo"
)

//
// Greg
//

type Greg struct {
        Session *discordgo.Session
        BotToken string
        BotPrefix []string
}

type GregChannel struct {
	Name string
	Guild string
}

type GregMessage struct {
	Prefixed bool
	Command string
	Params string
}

func (g *Greg) Start() {
	discord, err := discordgo.New("Bot " + g.BotToken)
	if err != nil {
		fmt.Println("Error creating Discord session", err)
		os.Exit(1)
	}
	g.Session = discord
	g.Session.AddHandler(goGregGo)
	// Do some listening, mate
	err = g.Session.Open()
	if err != nil {
		fmt.Println("Error opening connection", err)
		os.Exit(2)
	}
	fmt.Println("Greg is now Gregging in:")
	channels := g.getAllChannels()
	for _, channel := range channels {
		fmt.Println("- (" + channel.Guild + ") " + channel.Name)
	}
}

func (g *Greg) Stop() {
	fmt.Println("Greg is going to cease to Greg.")
	g.Session.Close()
}

func (g *Greg) getAllChannels() ([]GregChannel) {
	channelList := []GregChannel{}
	for _, guild := range g.Session.State.Guilds {
		GuildChannels, _ := g.Session.GuildChannels(guild.ID)
		for _, channel := range GuildChannels {
			if channel.Type != discordgo.ChannelTypeGuildText {
				continue
			}
			channelList = append(channelList, GregChannel{Name: channel.Name, Guild: guild.Name})
		}
	}
	return channelList
}

//
// Main
//

var BotToken = os.Getenv("BOT_TOKEN")
var BotPrefix = strings.Split(os.Getenv("BOT_PREFIX"),",")

func main() {
	greg := Greg{BotToken: BotToken, BotPrefix: BotPrefix}
	greg.Start()

	// Wait for the good old CTRL+C
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	greg.Stop()
}

func goGregGo(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Don't listen to yourself, Greg
	if m.Author.ID == s.State.User.ID {
		return
	}

	var result = parseMessage(m.Content, BotPrefix)
	if result.Prefixed {
		s.ChannelMessageSend(m.ChannelID, "Greg recognises your Greggin' and your command was `"+result.Command+"` with parameters `"+result.Params+"`")
	}
}

func parseMessage(m string, prefixes []string) (GregMessage) {
	result := GregMessage{Prefixed: false, Command: m, Params: ""}

	SplitMessage := strings.SplitN(strings.ToLower(m), " ", 3);
	MessagePrefix := SplitMessage[0]

	if len(SplitMessage) > 1 {
		result.Command = SplitMessage[1]
	}
	if len(SplitMessage) > 2 {
		result.Params = SplitMessage[2]
	}
	if strings.Trim(result.Command, " ") == "" {
		result.Prefixed = true
	}
	for _, p := range prefixes {
		var LowerPrefix = strings.ToLower(p)
		if strings.HasPrefix(MessagePrefix, LowerPrefix) || MessagePrefix == LowerPrefix {
			result.Prefixed = true
		}
	}
	return result
}

