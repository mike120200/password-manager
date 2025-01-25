/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pm",
	Short: "A secure password manager for storing and managing encrypted credentials",
	Long: `pm is a secure password manager designed to store and manage your credentials in an encrypted format.

This tool allows you to:
  - Securely store passwords with unique keys and optional platform information.
  - Retrieve stored passwords by their associated keys.
  - Update existing passwords or platform details.
  - List all stored credentials for easy management.
  - Delete a stored password by its key.
  - Backup all stored credentials to a file.
  - Restore credentials from a backup file.
  - Automatically backup all credentials every 500 seconds while the program is running.
  - List passwords associated with a specific platform.

All data is encrypted using AES encryption, ensuring your passwords are safe and protected.

Examples:
  - Store a new password:    pm add
  - Retrieve a password:     pm query
  - Update a password:       pm update
  - List all passwords:      pm list
  - Delete a password:       pm del
  - Backup all passwords:    pm backup
  - Restore from backup:     pm restore
  - List passwords by platform: pm pla

Automatic Backup:
  - While the program is running, a backup of all credentials will be created every 500 seconds.
  - The backup file will be saved in the same directory as the main database file.

For more information on a specific command, use 'pm [command] --help'.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.password_manager.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
