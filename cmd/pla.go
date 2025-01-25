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
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

// platformCmd represents the platform command
var platformCmd = &cobra.Command{
	Use:   "pla",
	Short: "List all passwords associated with a specific platform",
	Long: `List all passwords associated with a specific platform from the encrypted storage.

This command allows you to retrieve and display all passwords that are associated with
a specific platform (e.g., GitHub, Google, etc.). If no platform is provided as an argument,
you will be prompted to enter one.

Examples:
  - List passwords for a specific platform:
    pm pla GitHub

  - Enter platform interactively:
    pm pla
    Enter platform: GitHub

The output will display all keys and their corresponding passwords for the specified platform.
If no passwords are found for the platform, a message will be shown.`,
	Run: func(cmd *cobra.Command, args []string) {
		//初始化日志模块
		if err := zaplog.LoggerInit(); err != nil {
			color.Red.Println(err)
			return
		}

		var platform string
		var err error
		if len(args) > 0 {
			platform = args[0]
		} else {
			platform, err = input.GetInput("Enter platform")
			if err != nil {
				color.Red.Println(err)
				return
			}
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
		cnt := 0
		for k, v := range results {
			if !strings.Contains(v.Platform, platform) {
				continue
			}
			cnt++
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
		if cnt == 0 {
			color.Red.Println("No password found for platform " + platform)
		}
	},
}

func init() {
	rootCmd.AddCommand(platformCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// platformCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// platformCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
