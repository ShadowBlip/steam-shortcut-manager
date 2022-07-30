/*
MIT License

Copyright © 2022 William Edwards <shadowapex at gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package cmd

import (
	"fmt"

	"github.com/shadowblip/steam-shortcut-manager/pkg/chimera"
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

// removeCmd represents the remove command
var chimeraRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a Chimera shortcut from your library",
	Long:  `Remove a Chimera shortcut from your library`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		format := rootCmd.PersistentFlags().Lookup("output").Value.String()
		DebugPrintln("Using output format:", format)
		if !chimera.HasChimera() {
			ExitError(fmt.Errorf("no chimera config found at %v", chimera.ConfigDir), format)
		}

		// Get the platform flag
		platform := chimeraCmd.PersistentFlags().Lookup("platform").Value.String()

		// Read from the given shortcuts file
		shortcuts, err := chimera.LoadShortcuts(chimera.GetShortcutsFile(platform))
		if err != nil {
			ExitError(err, format)
		}

		// Find the shortcut to remove by name
		shortcutsList := []*chimera.Shortcut{}
		for _, sc := range shortcuts {
			if sc.Name == name {
				continue
			}
			shortcutsList = append(shortcutsList, sc)
		}

		// Save the shortcuts
		err = chimera.SaveShortcuts(chimera.GetShortcutsFile(platform), shortcutsList)
		if err != nil {
			ExitError(err, format)
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
	chimeraCmd.AddCommand(chimeraRemoveCmd)

	// Here you will define your flags and configuration settings.
	removeCmd.Flags().String("user", "all", "Steam user ID to remove the shortcut for")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
