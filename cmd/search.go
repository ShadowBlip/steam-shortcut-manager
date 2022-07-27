/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/shadowblip/steam-shortcut-manager/pkg/steamgriddb"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search --api-key <key> <name>",
	Short: "Search SteamGridDB for images",
	Args:  cobra.MinimumNArgs(1),
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, _ := cmd.Flags().GetString("api-key")
		if apiKey == "" {
			cmd.Help()
			fmt.Println("Error: API Key is required")
			os.Exit(1)
		}
		client := steamgriddb.NewClient(apiKey)
		results, err := client.Search(args[0])
		if err != nil {
			panic(err)
		}

		// Print the output
		format := rootCmd.PersistentFlags().Lookup("output").Value.String()
		switch format {
		case "term", "json":
			out, err := json.MarshalIndent(results, "", "  ")
			if err != nil {
				panic(err)
			}
			fmt.Println(string(out))
		default:
			panic("unknown output format: " + format)
		}

	},
}

func init() {
	steamgriddbCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
