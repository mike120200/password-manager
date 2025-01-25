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

// queryCmd represents the find command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Retrieve a stored password by key",
	Long: `Retrieve a stored password and its associated details by key from the encrypted storage.

This command allows you to securely retrieve a password and its optional platform information
using the unique key (e.g., service name or username) that was used when the password was stored.
For example:

  pm query
  Enter key: github_john.doe

  -------

  pm query github_john.doe

If the key exists, the corresponding password and platform (if available) will be decrypted
and displayed in the following format:
  github_john.doe (GitHub) : my-secure-password123`,
	Run: func(cmd *cobra.Command, args []string) {
		//初始化日志模块
		if err := zaplog.LoggerInit(); err != nil {
			color.Red.Println(err)
			return
		}
		var key string
		if len(args) != 0 {
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
			color.Red.Println("Key or account cannot be empty")
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
		password, platform, err := passwordInstance.GetPasswordWithKey(key)
		if err != nil {
			color.Red.Println(err)
			return
		}
		fmt.Println()
		if platform == "" {
			color.Blue.Println(key + " : " + password)
		} else {
			color.Blue.Printf("\n" + key + "(" + platform + ")" + " : " + password + "\n")
		}
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// findCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// findCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
