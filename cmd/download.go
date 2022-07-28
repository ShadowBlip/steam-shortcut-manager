/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/shadowblip/steam-shortcut-manager/pkg/steam"
	"github.com/spf13/cobra"
)

func discoverDownloadDir(cmd *cobra.Command, format string) string {
	// Handle when download dir is directly specified
	destinationDir, _ := cmd.PersistentFlags().GetString("destination-dir")
	if destinationDir != "" {
		return destinationDir
	}

	// Handle downloading for steam user IDs
	user, _ := cmd.PersistentFlags().GetString("user")
	if user != "" {
		users, err := steam.GetUsers()
		if err != nil {
			ExitError(err, format)
		}
		if !contains(users, user) {
			ExitError(fmt.Errorf("user not found"), format)
		}
		downloadDir, err := steam.GetImagesDir(user)
		if err != nil {
			ExitError(err, format)
		}
		return downloadDir
	}

	ExitError(fmt.Errorf("unable to discover download directory"), format)
	return ""
}

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download --api-key=<key>",
	Short: "Download SteamGridDB images for a given app",
	Long:  `Download SteamGridDB images for a given app`,
	Run: func(cmd *cobra.Command, args []string) {
		format := rootCmd.PersistentFlags().Lookup("output").Value.String()

		// Ensure we have a SteamGridDB API Key
		apiKey, _ := cmd.Flags().GetString("api-key")
		if apiKey == "" {
			cmd.Help()
			ExitError(fmt.Errorf("API key is required"), format)
		}

		// Get the directory to download to
		downloadDir := discoverDownloadDir(cmd, format)

		fmt.Println(downloadDir)
	},
}

func init() {
	steamgriddbCmd.AddCommand(downloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	downloadCmd.PersistentFlags().StringP("destination-dir", "d", "", "Destination directory to download images")
	downloadCmd.PersistentFlags().StringP("user", "u", "", "Steam user ID to download images for")
	downloadCmd.MarkFlagsMutuallyExclusive("destination-dir", "user")

	downloadCmd.PersistentFlags().Bool("only-hero", false, "Only download hero image")
	downloadCmd.PersistentFlags().Bool("only-grid", false, "Only download grid image")
	downloadCmd.PersistentFlags().Bool("only-icon", false, "Only download icon image")
	downloadCmd.PersistentFlags().Bool("only-logo", false, "Only download logo image")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
