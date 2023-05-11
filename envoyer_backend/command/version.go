package command

import (
	"envoyer/config"
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of application",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.Config.AppName + " 0.2.2 -- HEAD")
	},
}
