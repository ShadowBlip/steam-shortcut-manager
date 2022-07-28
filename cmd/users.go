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

	"github.com/shadowblip/steam-shortcut-manager/pkg/steam"
	"github.com/spf13/cobra"
)

// usersCmd represents the users command
var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "List current Steam user IDs",
	Long:  `List current Steam user IDs`,
	Run: func(cmd *cobra.Command, args []string) {
		format := rootCmd.PersistentFlags().Lookup("output").Value.String()
		users, err := steam.GetUsers()
		if err != nil {
			ExitError(err, format)
		}

		// Print the output
		switch format {
		case "term":
			for _, user := range users {
				fmt.Println(user)
			}
		case "json":
			out, err := json.MarshalIndent(users, "", "  ")
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
	rootCmd.AddCommand(usersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// usersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// usersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
