/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// steamgriddbCmd represents the steamgriddb command
var steamgriddbCmd = &cobra.Command{
	Use:   "steamgriddb",
	Short: "Search and download artwork from SteamGridDB",
	Long:  `Search and download artwork from SteamGridDB`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(steamgriddbCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	steamgriddbCmd.PersistentFlags().StringP("api-key", "k", "", "SteamGridDB API Key")
	steamgriddbCmd.MarkFlagRequired("api-key")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// steamgriddbCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
