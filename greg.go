package main

import (
//	"fmt"
	"os"
	"os/signal"
	"syscall"
	"strings"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"log"
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

type GregWeather struct {
	Query string 			`json:"query"`
	Coords GregWeatherCoords 	`json:"coords"`
	Restrict_au string 		`json:"restrict_au"`
	Api_key string 			`json:"api_key"`
	Location string			`json:"location"`
	Country string			`json:"country"`
	Source string			`json:"source"`
	Url string			`json:"url"`
	Station string			`json:"station"`
	Temp string			`json:"temp"`
	Feels string			`json:"feels"`
	Humidity string			`json:"humidity"`
	Rain string			`json:"rain"`
	Wind string			`json:"wind"`
	Summary string			`json:"summary"`
	Update string			`json:"update"`
	Icon string			`json:"icon"`
}

type GregWeatherCoords struct {
	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
}

func (g *Greg) Start() {
	discord, err := discordgo.New("Bot " + g.BotToken)
	if err != nil {
		log.Println("Error creating Discord session", err)
		os.Exit(1)
	}
	g.Session = discord
	g.Session.AddHandler(goGregGo)
	// Do some listening, mate
	err = g.Session.Open()
	if err != nil {
		log.Println("Error opening connection", err)
		os.Exit(2)
	}
	log.Println("Greg is now Gregging in:")
	channels := g.getAllChannels()
	for _, channel := range channels {
		log.Println("- (" + channel.Guild + ") " + channel.Name)
	}
	if os.Getenv("BOT_WEBHOOK_ID") != "" && os.Getenv("BOT_WEBHOOK_TOKEN") != "" {
		log.Println("Webhook logs are enabled")
	}
}

func (g *Greg) Stop() {
	log.Println("Greg is going to cease to Greg.")
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
		s.ChannelTyping(m.ChannelID)
		if result.Command == "weather" || result.Command == "w" {
			logAction(s, "<" + m.Author.Username + "> requested weather for '" + result.Params + "'")
			greggo := getGregWeather(m.Author.ID, result.Params)
			embed := &discordgo.MessageEmbed{
				Color: 0x0072bb,
				Title: "Weather for "+greggo.Location,
				Description: "It is currently " + greggo.Temp + "°C. Feels like " + greggo.Feels + "°C.\n"+greggo.Summary,
				Fields: []*discordgo.MessageEmbedField{
					&discordgo.MessageEmbedField{Name:"Humidity", Value:greggo.Humidity, Inline: true},
					&discordgo.MessageEmbedField{Name:"Rain", Value:greggo.Rain, Inline: true},
					&discordgo.MessageEmbedField{Name:"Wind", Value:greggo.Wind, Inline: true},
			        },
				Footer: &discordgo.MessageEmbedFooter{Text:"Requested by " + m.Author.Username + " Station: " + greggo.Station + " Last Update: " + greggo.Update, IconURL: m.Author.AvatarURL("")},
			}
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
		} else {
			s.ChannelMessageSend(m.ChannelID, "Greggo does not know what you're on about.")
		}
	}
}

func talkToRoy(command string, user string, query string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET","http://roy_api_1/", nil)
	if err != nil {
		return "greg is the worst", err
	}
	q := req.URL.Query()
	q.Add("source", "discord")
	q.Add("get", command)
	q.Add("uniqname", user)
	q.Add("query", query)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return "greg cant do", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "greg cant read", err
	}
	return string(body), nil
}

func getGregWeather(user string, params string) (GregWeather) {
	result, _ := talkToRoy("weather", string(user), params);
	gw := GregWeather{}
	json.Unmarshal([]byte(result), &gw)
	return gw
}

func parseMessage(m string, prefixes []string) (GregMessage) {
	result := GregMessage{Prefixed: false, Command: m, Params: ""}

	SplitMessage := strings.SplitN(strings.ToLower(m), " ", 3);
	MessagePrefix := SplitMessage[0]

	if len(SplitMessage) > 1 {
		result.Command = strings.ToLower(SplitMessage[1])
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

func logAction(s *discordgo.Session, message string) {
	log.Println(message)
	logToDiscord(s, message)
}

func logToDiscord(s *discordgo.Session, message string) {
	wid := os.Getenv("BOT_WEBHOOK_ID")
	wtoken := os.Getenv("BOT_WEBOOK_TOKEN")
	if wid != "" && wtoken != "" {
		params := discordgo.WebhookParams{Username: "Greg", Content: "```"+message+"```"}
		err := s.WebhookExecute(wid, wtoken, false, &params)
	        if err != nil {
			log.Println("logToDiscord failed: Unable to perform request")
	        }
	}
}

