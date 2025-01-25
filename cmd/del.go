/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	zaplog "password_manager/common/log"
	"password_manager/service/aes"
	dbfilekit "password_manager/service/dbfile_Kit"
	"password_manager/service/input"
	"password_manager/service/password"
	secretkey "password_manager/service/secret_key"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

// delCmd represents the del command
var delCmd = &cobra.Command{
	Use:   "del",
	Short: "Delete a stored password by its key",
	Long: `Delete a stored password from the secure database.

This command allows you to remove a password that was previously stored using a unique key.
Before deletion, a confirmation prompt will be displayed to prevent accidental deletions.

Example usage:

  pm del

You will be prompted to enter the key associated with the password you want to delete.
If the key exists, the password will be permanently removed from the database.

Warning: This action is irreversible.`,
	Run: func(cmd *cobra.Command, args []string) {
		//初始化日志模块
		if err := zaplog.LoggerInit(); err != nil {
			color.Red.Println(err)
			return
		}
		var key string
		if len(args) > 0 {
			key = args[0]
		} else {
			var err error
			//获取密码的键
			key, err = input.GetInput("Enter key or account")
			if err != nil {
				color.Red.Println(err)
				return
			}
		}
		if key == "" {
			color.Red.Println("key is empty")
			return
		}
		actionConfirm, err := input.GetInput("Are you sure to delete key:" + key + " (y/n)")
		if err != nil {
			color.Red.Println(err)
			return
		}
		if actionConfirm == "n" {
			fmt.Println("delete key:" + key + " cancel")
			return
		} else if actionConfirm != "y" {
			fmt.Println("invalid input")
			return
		}
		//初始化密钥模块
		secretKeyInstance := secretkey.NewSecretKey()

		//初始化数据库模块
		kitInstance := dbfilekit.NewDBKit(secretKeyInstance)
		if err := kitInstance.Init(); err != nil {
			color.Red.Println(err)
			return
		}

		//获取数据库
		db, err := kitInstance.GetDB()
		if err != nil {
			color.Red.Println(err)
			return
		}

		//获取密钥
		secretKey, err := secretKeyInstance.GetSecretKey()
		if err != nil {
			color.Red.Println(err)
			return
		}
		//初始化加密模块
		aesInstance := aes.NewAesService(secretKey)
		//初始化密码保存模块
		passwordInstance := password.NewPasswordService(aesInstance, db)
		err = passwordInstance.DeletePassword(key)
		if err != nil {
			color.Red.Println(err)
			return
		}
		color.Green.Println("delete key:" + key + " success")
	},
}

func init() {
	rootCmd.AddCommand(delCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// delCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// delCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
