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
	Short: "Update a stored password or platform",
	Long: `Update a stored password or platform information in the encrypted storage.

This command allows you to update the password or platform associated with a specific key.
You will be prompted to enter the key, the new password (optional), and the new platform (optional). 
For example:

  pm update
  Enter key: github_john.doe
  Enter new password (optional): ********
  Enter new platform (optional): GitHub

If the key exists, the corresponding password or platform will be updated. 
If no new password or platform is provided, the existing values will remain unchanged.`,
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
		//获取新密码
		newPassword, err := input.GetOptionalPassword("Enter new password(optional, press Enter to skip)")
		if err != nil {
			color.Red.Println(err)
			return
		}
		//获取新平台信息
		newPlatform, err := input.GetOptionalInput("Enter new platform(optional, press Enter to skip)")
		if err != nil {
			color.Red.Println(err)
			return
		}
		if newPlatform == "" && newPassword == "" {
			color.Red.Println("New password and new platform cannot be empty at the same time")
			return
		}
		//获取新用户名
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
		err = passwordInstance.UpdatePassword(key, newPassword, newPlatform)
		if err != nil {
			color.Red.Println(err)
			return
		}
		err = kitInstance.BackupDB()
		if err != nil {
			color.Red.Println(err)
			return
		}
		color.Green.Println("backup successfully")
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
