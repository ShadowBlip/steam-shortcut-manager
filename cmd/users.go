/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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
