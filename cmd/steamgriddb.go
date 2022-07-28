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
