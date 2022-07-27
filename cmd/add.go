/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/shadowblip/steam-shortcut-manager/pkg/shortcut"
	"github.com/shadowblip/steam-shortcut-manager/pkg/steam"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <name> <exe>",
	Short: "Add a Steam shortcut to your steam library",
	Args:  cobra.ExactArgs(2),
	Long:  `Adds a Steam shortcut to your library`,
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		exe := args[1]

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

			// Generate a new shortcut from the cli flags
			newShortcut := newShortcutFromFlags(cmd, name, exe)
			shortcuts.Add(newShortcut)

			// Write the changes
			err = shortcut.Save(shortcuts, shortcutsPath)
			if err != nil {
				panic(err)
			}
		}
	},
}

// Creates a new shortcut object from command-line flags
func newShortcutFromFlags(cmd *cobra.Command, name, exe string) *shortcut.Shortcut {
	getString := func(name string) string {
		res, _ := cmd.Flags().GetString(name)
		return res
	}
	getBool := func(name string) int {
		res, _ := cmd.Flags().GetBool(name)
		return boolToInt(res)
	}
	shortcutConfiger := func(s *shortcut.Shortcut) {
		s.AllowDesktopConfig = getBool("allow-desktop-config")
		s.AllowOverlay = getBool("allow-overlay")
		s.FlatpakAppID = getString("flatpak-id")
		s.IsHidden = getBool("is-hidden")
		s.LaunchOptions = getString("launch-options")
		s.OpenVR = getBool("openvr")
		s.ShortcutPath = getString("shortcut-path")
		s.StartDir = getString("start-dir")
		s.Appid = int64(shortcut.CalculateAppID(exe, name))
		s.Icon = getString("icon")

		s.Tags = map[string]interface{}{}
		tags, _ := cmd.Flags().GetStringSlice("tags")
		for key, tag := range tags {
			s.Tags[fmt.Sprintf("%v", key)] = tag
		}
	}
	shortcut := shortcut.NewShortcut(name, exe, shortcutConfiger)
	return shortcut
}

// boolToInt will convert bool flag to int
func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Optional flags
	addCmd.Flags().Bool("allow-desktop-config", true, "Allow desktop config")
	addCmd.Flags().Bool("allow-overlay", true, "Allow steam overlay")
	addCmd.Flags().Bool("is-hidden", false, "Whether or not the shortcut is hidden")
	addCmd.Flags().String("flatpak-id", "", "Flatpak ID of the shortcut")
	addCmd.Flags().String("launch-options", "", "Launch options for the shortcut")
	addCmd.Flags().Bool("openvr", false, "Use OpenVR for the shortcut")
	addCmd.Flags().String("shortcut-path", "", "Path to the shortcut file for this application")
	addCmd.Flags().String("start-dir", "", "Working directory where the app is started")
	addCmd.Flags().String("icon", "", "Path to the icon to use for this application")
	addCmd.Flags().StringSlice("tags", []string{}, "Comma-separated list of tags")
	addCmd.Flags().String("user", "all", "Steam user ID to add the shortcut for")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
