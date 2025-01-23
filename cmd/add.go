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

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Store a password key-value pair",
	Long: `Store a password key-value pair in the secure storage.

This command allows you to securely store a password by associating it with a unique key.
For example, to store a password for a website, you can use:

  pm add
  Enter key: github_john.doe
  Enter value: mysecurepassword123

The key should be a unique identifier (e.g., service name or username),
and the value is the password you want to store.`,
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
			key, err = input.GetInput("Enter key")
			if err != nil {
				color.Red.Println(err)
				return
			}
		}

		//获取密码的值
		passwordValue, err := input.GetPasswordInput("Enter password")
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
		err = passwordInstance.SavePassword(key, passwordValue)
		if err != nil {
			color.Red.Println(err)
			return
		}
		color.Green.Println("Password saved successfully!")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
