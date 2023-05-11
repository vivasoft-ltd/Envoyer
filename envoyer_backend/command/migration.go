package command

import (
	"envoyer/config/service_name"
	"envoyer/dic"
	"envoyer/model/entity"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func init() {
	rootCmd.AddCommand(migrationCmd)
}

var migrationCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run migration",
	Run: func(cmd *cobra.Command, args []string) {
		db := dic.Container.Get(service_name.DbService).(*gorm.DB)
		_ = db.AutoMigrate(&entity.Application{})
		_ = db.AutoMigrate(&entity.Client{})
		_ = db.AutoMigrate(&entity.Event{})
		_ = db.AutoMigrate(&entity.Template{})
		_ = db.AutoMigrate(&entity.User{})
		_ = db.AutoMigrate(&entity.Provider{})
		_ = db.AutoMigrate(&entity.ErrorLog{})
		//_ = db.AutoMigrate(&entity.SmsPriority{})
	},
}
