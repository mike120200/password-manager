/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	zaplog "password_manager/common/log"
	dbfilekit "password_manager/service/dbfile_Kit"
	secretkey "password_manager/service/secret_key"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore credentials from a backup file",
	Long: `Restore credentials from a backup file to the main database.

This command allows you to restore all previously backed-up credentials from a backup file
into the main database. The backup file should be located in the same directory as the
main database file and named 'backup.json'.

Example:
  pm restore

The restore process will overwrite any existing data in the main database with the data
from the backup file. Use this command with caution.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			color.Red.Println("invalid args")
			return
		}
		//初始化日志模块
		if err := zaplog.LoggerInit(); err != nil {
			color.Red.Println(err)
			return
		}
		//初始化密钥模块
		secretKeyInstance := secretkey.NewSecretKey()

		//初始化数据库模块
		kitInstance := dbfilekit.NewDBKit(secretKeyInstance)
		if err := kitInstance.RestoreDB(); err != nil {
			color.Red.Println(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// restoreCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// restoreCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
