package cmd

import (
	"discord-bots/repost"
	"github.com/spf13/cobra"
)

var token, repostChannelID string

func init() {
	rootCmd.AddCommand(startReposterCMD())
}

func startReposterCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "repost",
		Short:   "repost",
		Long:    "repost starts reposter bot.",
		Run:     runRepostBot,
		Example: "repost --token=my_token",
	}

	cmd.Flags().StringVarP(
		&token,
		"token",
		"t",
		"",
		"token of the bot application",
	)

	cmd.Flags().StringVarP(
		&repostChannelID,
		"channel",
		"c",
		"",
		"channel id in which you want to repost content",
	)

	return cmd
}

func runRepostBot(cmd *cobra.Command, args []string) {
	repost.Run(token, repostChannelID)
}
