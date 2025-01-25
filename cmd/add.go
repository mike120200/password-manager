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
	Short: "Securely store a password",
	Long: `Securely store a password in the encrypted storage.

This command allows you to save a password by associating it with a unique identifier (key).
For example, to store a password for a website, you can use:

  pm add
  Enter key: github_john.doe
  Enter password: my-secure-password123
  Enter platform (optional, press Enter to skip): GitHub

The key should be a unique identifier (e.g., service name, username, or account),
and the value is the password you want to store. The optional platform field can 
be used to specify the service or application associated with the password.
If you do not need to specify a platform, simply press Enter to skip.`,
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
			color.Red.Println("Key or account cannot be empty")
			return
		}
		//获取密码的值
		passwordValue, err := input.GetPasswordInput("Enter password")
		if err != nil {
			color.Red.Println(err)
			return
		}
		//输入平台
		platform, err := input.GetOptionalInput("Enter platform (optional, press Enter to skip)")
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
		err = passwordInstance.SavePassword(key, passwordValue, platform)
		if err != nil {
			color.Red.Println(err)
			return
		}
		//备份
		err = kitInstance.BackupDB()
		if err != nil {
			color.Red.Println(err)
			return
		}
		color.Green.Println("Password saved  and backup successfully!")
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
