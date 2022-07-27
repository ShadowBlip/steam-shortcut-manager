/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/shadowblip/steam-shortcut-manager/pkg/image"
	"github.com/shadowblip/steam-shortcut-manager/pkg/steamgriddb"
	"github.com/spf13/cobra"
)

// SearchType is a bitmask of different kinds of searches
type SearchType uint8

func (b SearchType) Set(flag SearchType) SearchType    { return b | flag }
func (b SearchType) Clear(flag SearchType) SearchType  { return b &^ flag }
func (b SearchType) Toggle(flag SearchType) SearchType { return b ^ flag }
func (b SearchType) Has(flag SearchType) bool          { return b&flag != 0 }

const (
	SearchGrids SearchType = 1 << iota
	SearchHeroes
	SearchLogos
	SearchIcons
)

// Combined search output from SteamGridDB
type SearchOutput struct {
	Details steamgriddb.SearchResponseData  `json:"details"`
	Heroes  []steamgriddb.ImageResponseData `json:"heroes"`
	Grids   []steamgriddb.GridResponseData  `json:"grids"`
	Logos   []steamgriddb.ImageResponseData `json:"logos"`
	Icons   []steamgriddb.ImageResponseData `json:"icons"`
}

// Prints the search output to the terminal
func (s *SearchOutput) Print(client *steamgriddb.Client) {
	fmt.Println(s.Details.Name)
	fmt.Println("  App ID:", s.Details.ID)
	for _, data := range s.Grids {
		filename := path.Base(data.Thumb)
		if image.CanDisplay {
			err := client.CachedDownload(data.Thumb, fmt.Sprintf("/tmp/%v", filename))
			if err != nil {
				continue
			}
		}
		fmt.Println("  Grid Images")
		fmt.Println("    Author:", data.Author.Name)
		fmt.Println("    URL:", data.URL)
		if image.CanDisplay {
			image.Display("/tmp/" + filename)
		}
	}
	for _, data := range s.Logos {
		filename := path.Base(data.Thumb)
		if image.CanDisplay {
			err := client.CachedDownload(data.Thumb, fmt.Sprintf("/tmp/%v", filename))
			if err != nil {
				continue
			}
		}
		fmt.Println("  Logo Images")
		fmt.Println("    Author:", data.Author.Name)
		fmt.Println("    URL:", data.URL)
		if image.CanDisplay {
			image.Display("/tmp/" + filename)
		}
	}
	for _, data := range s.Icons {
		filename := path.Base(data.Thumb)
		if image.CanDisplay {
			err := client.CachedDownload(data.Thumb, fmt.Sprintf("/tmp/%v", filename))
			if err != nil {
				continue
			}
		}
		fmt.Println("  Icon Images")
		fmt.Println("    Author:", data.Author.Name)
		fmt.Println("    URL:", data.URL)
		if image.CanDisplay {
			image.Display("/tmp/" + filename)
		}
	}
	for _, data := range s.Heroes {
		filename := path.Base(data.Thumb)
		if image.CanDisplay {
			err := client.CachedDownload(data.Thumb, fmt.Sprintf("/tmp/%v", filename))
			if err != nil {
				continue
			}
		}
		fmt.Println("  Hero Images")
		fmt.Println("    Author:", data.Author.Name)
		fmt.Println("    URL:", data.URL)
		if image.CanDisplay {
			image.Display("/tmp/" + filename)
		}
	}
}

// search for SteamGridDB images
func search(cmd *cobra.Command, args []string, kind SearchType) {
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

	// Error if not success
	if !results.Success {
		panic(fmt.Errorf("%v", results.Errors))
	}

	// Filter our results
	maxResults := getFlagInt(cmd, "num-results")
	if maxResults > len(results.Data) {
		maxResults = len(results.Data)
	}
	results.Data = results.Data[:maxResults]

	// Create a structure to hold our results
	searchResult := map[string]*SearchOutput{}

	// Get all images for each found result
	for _, result := range results.Data {
		appID := fmt.Sprintf("%v", result.ID)
		searchResult[appID] = &SearchOutput{
			Details: result,
			Grids:   []steamgriddb.GridResponseData{},
			Heroes:  []steamgriddb.ImageResponseData{},
			Logos:   []steamgriddb.ImageResponseData{},
			Icons:   []steamgriddb.ImageResponseData{},
		}
		maxImages := getFlagInt(cmd, "num-images")

		// Get all grid images
		if kind.Has(SearchGrids) {
			grids, err := client.GetGrids(appID)
			if err != nil {
				panic(err)
			}
			num := maxImages
			if num > len(grids.Data) {
				num = len(grids.Data)
			}
			grids.Data = grids.Data[:num]
			searchResult[appID].Grids = grids.Data
		}

		// Get all hero images
		if kind.Has(SearchHeroes) {
			heroes, err := client.GetHeroes(appID)
			if err != nil {
				panic(err)
			}
			num := maxImages
			if num > len(heroes.Data) {
				num = len(heroes.Data)
			}
			heroes.Data = heroes.Data[:num]
			searchResult[appID].Heroes = heroes.Data
		}

		// Get all logo images
		if kind.Has(SearchLogos) {
			logos, err := client.GetLogos(appID)
			if err != nil {
				panic(err)
			}
			num := maxImages
			if num > len(logos.Data) {
				num = len(logos.Data)
			}
			logos.Data = logos.Data[:num]
			searchResult[appID].Logos = logos.Data
		}

		// Get all icon images
		if kind.Has(SearchIcons) {
			icons, err := client.GetIcons(appID)
			if err != nil {
				panic(err)
			}
			num := maxImages
			if num > len(icons.Data) {
				num = len(icons.Data)
			}
			icons.Data = icons.Data[:num]
			searchResult[appID].Icons = icons.Data
		}

	}

	// Print the output
	format := rootCmd.PersistentFlags().Lookup("output").Value.String()
	switch format {
	case "term":
		for _, result := range searchResult {
			result.Print(client)
		}
	case "json":
		out, err := json.MarshalIndent(searchResult, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(out))
	default:
		panic("unknown output format: " + format)
	}
}

func getFlagInt(cmd *cobra.Command, name string) int {
	result, _ := cmd.PersistentFlags().GetInt(name)
	if result == 0 {
		result, _ = cmd.Parent().PersistentFlags().GetInt(name)
	}
	return result
}

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search --api-key <key> <name>",
	Short: "Search SteamGridDB for images",
	Args:  cobra.MinimumNArgs(1),
	Long:  `Search SteamGridDB for images. Returns all image types.`,
	Run: func(cmd *cobra.Command, args []string) {
		var searchFlags SearchType
		if ok, _ := cmd.Flags().GetBool("grids"); ok {
			searchFlags = searchFlags.Set(SearchGrids)
		}
		if ok, _ := cmd.Flags().GetBool("logos"); ok {
			searchFlags = searchFlags.Set(SearchLogos)
		}
		if ok, _ := cmd.Flags().GetBool("icons"); ok {
			searchFlags = searchFlags.Set(SearchIcons)
		}
		if ok, _ := cmd.Flags().GetBool("heroes"); ok {
			searchFlags = searchFlags.Set(SearchHeroes)
		}
		search(cmd, args, searchFlags)
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
	searchCmd.PersistentFlags().IntP("num-results", "n", 1, "Number of search results to return")
	searchCmd.PersistentFlags().Int("num-images", 1, "Number of image results to return for a given image type")
	searchCmd.PersistentFlags().Bool("grids", true, "Include grid images in search")
	searchCmd.PersistentFlags().Bool("heroes", true, "Include hero images in search")
	searchCmd.PersistentFlags().Bool("icons", true, "Include icon images in search")
	searchCmd.PersistentFlags().Bool("logos", true, "Include logo images in search")
}
