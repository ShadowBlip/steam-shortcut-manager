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
	"encoding/json"
	"fmt"

	"github.com/shadowblip/steam-shortcut-manager/pkg/chimera"
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
		format := rootCmd.PersistentFlags().Lookup("output").Value.String()
		users, err := steam.GetUsers()
		if err != nil {
			ExitError(err, format)
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
				ExitError(err, format)
			}
			results[user] = shortcuts
		}

		// Print the output
		switch format {
		case "term":
			for user, shortcuts := range results {
				if shortcuts.Shortcuts == nil || len(shortcuts.Shortcuts) == 0 {
					continue
				}
				fmt.Println("User:", user)
				for _, sc := range shortcuts.Shortcuts {
					fmt.Println("  ", sc.AppName)
					fmt.Println("    AppId:", sc.Appid)
					fmt.Println("    Executable:    ", sc.Exe)
					fmt.Println("    Launch Options:", sc.LaunchOptions)
				}
			}
		case "json":
			out, err := json.MarshalIndent(results, "", "  ")
			if err != nil {
				ExitError(err, format)
			}
			fmt.Println(string(out))
		default:
			panic("unknown output format: " + format)
		}
	},
}

// chimeraListCmd represents the list command
var chimeraListCmd = &cobra.Command{
	Use:   "list",
	Short: "List currently registered Chimera shortcuts",
	Long:  `Lists all of the shortcuts registered in Chimera`,
	Run: func(cmd *cobra.Command, args []string) {
		format := rootCmd.PersistentFlags().Lookup("output").Value.String()
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

		// Print the output
		switch format {
		case "term":
			for _, sc := range shortcuts {
				fmt.Println(sc.Name)
				fmt.Println("  Executable:", sc.Cmd)
			}
		case "json":
			out, err := json.MarshalIndent(shortcuts, "", "  ")
			if err != nil {
				ExitError(err, format)
			}
			fmt.Println(string(out))
		default:
			panic("unknown output format: " + format)
		}

	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	chimeraCmd.AddCommand(chimeraListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
