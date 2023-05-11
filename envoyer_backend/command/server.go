package command

import (
	"envoyer/config"
	"envoyer/dic"
	"envoyer/route"
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run server",
	Run: func(cmd *cobra.Command, args []string) {
		router := route.Setup(dic.Builder)
		fmt.Println("Running server on port: " + config.Config.Port)
		router.Run(":" + config.Config.Port)
	},
}
