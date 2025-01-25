/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	zaplog "password_manager/common/log"
	"password_manager/service/aes"
	dbfilekit "password_manager/service/dbfile_Kit"
	"password_manager/service/password"
	secretkey "password_manager/service/secret_key"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all stored passwords and their details",
	Long: `List all stored passwords and their associated details from the encrypted storage.

This command retrieves and displays all stored password entries, including the key, password, 
and optional platform information. Each entry is displayed in a clear format for easy reading. 
For example:

  pm list

Output:
  github_john.doe (GitHub) : my-secure-password123
  email_jane.doe : another-password

If a platform is associated with a password, it will be shown in parentheses next to the key.`,
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
		results, err := passwordInstance.GetAllPasswords()
		if err != nil {
			color.Red.Println(err)
			return
		}
		for k, v := range results {
			if v.Platform == "" {
				color.Blue.Printf("\n" + k + " : ")
				color.Green.Printf(v.Password + "\n")
			} else {
				color.Blue.Printf("\n" + k)
				color.Cyan.Printf(" (" + v.Platform + ") : ")
				color.Green.Printf(v.Password + "\n")
			}

		}
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
