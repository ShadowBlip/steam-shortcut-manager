/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/shadowblip/steam-shortcut-manager/pkg/shortcut"
	"github.com/shadowblip/steam-shortcut-manager/pkg/steam"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List currently registered Steam shortcuts",
	Long:  `Lists all of the shortcuts registered in Steam`,
	Run: func(cmd *cobra.Command, args []string) {
		users, err := steam.GetUsers()
		if err != nil {
			panic(err)
		}

		// Fetch all shortcuts
		results := map[string]*shortcut.Shortcuts{}
		for _, user := range users {
			if !steam.HasShortcuts(user) {
				continue
			}
			shortcutsPath, _ := steam.GetShortcutsPath(user)
			shortcuts, err := shortcut.Load(shortcutsPath)
			if err != nil {
				panic(err)
			}
			results[user] = shortcuts
		}

		// Print the output
		format := rootCmd.PersistentFlags().Lookup("output").Value.String()
		switch format {
		case "json", "term":
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
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
