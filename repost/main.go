package repost

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func Run(token, repostChannelID string) {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the repost func as a callback for MessageCreate events.
	dg.AddHandler(repost(repostChannelID))

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsDirectMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func repost(repostChannelID string) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore all messages created by the bot itself
		// This isn't required in this specific example but it's a good practice.
		if m.Author.ID == s.State.User.ID {
			return
		}

		cmdRepost := "!repost"
		if strings.HasPrefix(m.Content, cmdRepost) {
			fmt.Println("processing repost...")
			args := strings.Split(m.Content, " ")
			if len(args) != 3 {
				fmt.Println("missing parameters", len(args))
				return
			}
			channelID := args[1]
			messageID := args[2]
			message, err := s.ChannelMessage(channelID, messageID)
			if err != nil {
				fmt.Println(err)
				return
			}
			originChannel, err := s.Channel(channelID)
			if err != nil {
				fmt.Println(err)
				return
			}

			var repostMessage strings.Builder
			repostMessage.WriteString("-------------------------------------------------\n")
			repostMessage.WriteString(fmt.Sprintf("Message de %s du channel #%s\n", message.Author.Username, originChannel.Name))
			repostMessage.WriteString(fmt.Sprintf("```%s```", message.Content))
			_, err = s.ChannelMessageSend(repostChannelID, repostMessage.String())

			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
