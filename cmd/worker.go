package cmd

import (
	bot2 "github.com/sepuka/campaner/internal/bot"
	"github.com/sepuka/campaner/internal/def"
	"github.com/sepuka/campaner/internal/def/bot"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(workerCmd)
}

var (
	workerCmd = &cobra.Command{
		Use: `worker`,
		RunE: func(cmd *cobra.Command, args []string) error {
			botWorker, err := def.Container.SafeGet(bot.WorkerDef)
			if err != nil {
				return err
			}

			return botWorker.(*bot2.Worker).Work()
		},
	}
)
