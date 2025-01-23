/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	zaplog "password_manager/common/log"
	"password_manager/service/aes"
	dbfilekit "password_manager/service/dbfile_Kit"
	"password_manager/service/input"
	"password_manager/service/password"
	secretkey "password_manager/service/secret_key"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a stored password",
	Long: `Update a stored password in the secure storage.

This command allows you to update the password associated with a specific key.
You will be prompted to enter the key and the new password. For example:

  pm update
  Enter key: github_john.doe
  Enter new password: ********

The key should match the one used when the password was originally stored.
If the key exists, the corresponding password will be updated.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			color.Red.Println("invalid input")
			return
		}
		//初始化日志模块
		if err := zaplog.LoggerInit(); err != nil {
			color.Red.Println(err)
			return
		}
		//获取密码的键
		key, err := input.GetInput("Enter key")
		if err != nil {
			color.Red.Println(err)
			return
		}
		//获取新密码
		newPassword, err := input.GetInput("Enter new password")
		if err != nil {
			color.Red.Println(err)
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
		err = passwordInstance.UpdatePassword(key, newPassword)
		if err != nil {
			color.Red.Println(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
