package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var prefix string
var commands map[string]interface{}
var messageEvents map[string]interface{}

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("There was an error reading your bot config file: %s", err.Error())
	}
	contents, _ := io.ReadAll(file)
	var botRules map[string]interface{}
	err = json.Unmarshal(contents, &botRules)
	if err != nil {
		log.Fatalf("There was an error parsing your bot config file: %s. This is usually because of a trailing comma, or other json related error", err.Error())
	}
	if botRules["commands"] != nil {
		commandrules := botRules["commands"].(map[string]interface{})
		prefix = commandrules["prefix"].(string)
		commands = commandrules["commands"].(map[string]interface{})
	}
	if botRules["messages"] != nil {
		messageEvents = botRules["messages"].(map[string]interface{})
	}

	bot, _ := discordgo.New("Bot " + botRules["token"].(string))
	bot.Identify.Intents = discordgo.IntentsGuildMessages
	bot.AddHandler(messageCreate)
	err = bot.Open()
	if err != nil {
		log.Fatalf("Error opening connection to discord%s", err.Error())
	}
	fmt.Println("Connection opened")
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	c := m.Content
	if val, ok := messageEvents[c]; ok {
		specialReply(s, m, val)
	}
	if commands == nil {
		return
	}
	if !strings.HasPrefix(c, prefix) {
		return
	}
	cmdstr := c[len(prefix):]
	if commands[cmdstr] == nil {
		return
	}
	steps := commands[cmdstr]
	specialReply(s, m, steps)
}

func specialReply(s *discordgo.Session, m *discordgo.MessageCreate, msi interface{}) {
	switch reflect.TypeOf(msi).Kind() {
	case reflect.String:
		{
			s.ChannelMessageSendReply(m.ChannelID, msi.(string), m.Reference())
		}
	case reflect.Map:
		{
			mapsteps := msi.(map[string]interface{})
			if mapsteps["reply"] != nil {
				s.ChannelMessageSendReply(m.ChannelID, mapsteps["reply"].(string), m.Reference())
			}
			if mapsteps["dm"] != nil {
				c, err := s.UserChannelCreate(m.Author.ID)
				if err != nil {
					fmt.Printf("Error creating private message with %s. %e", m.Author.Email, err)
				}
				s.ChannelMessageSend(c.ID, mapsteps["dm"].(string))
			}
		}
	}
}
