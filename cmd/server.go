package cmd

import (
	bot2 "github.com/sepuka/campaner/internal/bot"
	"github.com/sepuka/campaner/internal/def"
	"github.com/sepuka/campaner/internal/def/bot"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var (
	serverCmd = &cobra.Command{
		Use: `server`,
		RunE: func(cmd *cobra.Command, args []string) error {
			botListener, err := def.Container.SafeGet(bot.BotDef)
			if err != nil {
				return err
			}

			return botListener.(*bot2.Bot).Listen()
		},
	}
)
