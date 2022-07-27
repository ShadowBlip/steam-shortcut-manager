/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/shadowblip/steam-shortcut-manager/pkg/shortcut"
	"github.com/shadowblip/steam-shortcut-manager/pkg/steam"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a Steam shortcut from your library",
	Long:  `Remove a Steam shortcut from your library`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		// Fetch all users
		users, err := steam.GetUsers()
		if err != nil {
			panic(err)
		}

		// Check to see if we're fetching for just one user
		onlyForUser := cmd.Flags().Lookup("user").Value.String()

		// Fetch all shortcuts
		for _, user := range users {
			if !steam.HasShortcuts(user) {
				continue
			}
			if onlyForUser != "all" && onlyForUser != user {
				continue
			}

			shortcutsPath, _ := steam.GetShortcutsPath(user)
			shortcuts, err := shortcut.Load(shortcutsPath)
			if err != nil {
				panic(err)
			}

			// Find the shortcut to remove by name
			shortcutsList := []shortcut.Shortcut{}
			for _, sc := range shortcuts.Shortcuts {
				if sc.AppName == name {
					continue
				}
				shortcutsList = append(shortcutsList, sc)
			}

			// Create a new shortcuts object that we will save
			newShortcuts := &shortcut.Shortcuts{
				Shortcuts: map[string]shortcut.Shortcut{},
			}
			for key, sc := range shortcutsList {
				newShortcuts.Shortcuts[fmt.Sprintf("%v", key)] = sc
			}

			// Write the changes
			err = shortcut.Save(newShortcuts, shortcutsPath)
			if err != nil {
				panic(err)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	// Here you will define your flags and configuration settings.
	removeCmd.Flags().String("user", "all", "Steam user ID to remove the shortcut for")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
