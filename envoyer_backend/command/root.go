package command

import (
	"envoyer/config"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Short: "Application description",
	Long: `Long
                application
                description`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please use command server. Like go run main.go server")
	},
}

func Execute() {
	rootCmd.Use = config.Config.AppCommand
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
