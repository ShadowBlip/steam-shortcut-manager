/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
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
		fmt.Println("    ID:", data.ID)
		fmt.Println("    Style:", data.Style)
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
		fmt.Println("    ID:", data.ID)
		fmt.Println("    Style:", data.Style)
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
		fmt.Println("    ID:", data.ID)
		fmt.Println("    Style:", data.Style)
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
		fmt.Println("    ID:", data.ID)
		fmt.Println("    Style:", data.Style)
		fmt.Println("    Author:", data.Author.Name)
		fmt.Println("    URL:", data.URL)
		if image.CanDisplay {
			image.Display("/tmp/" + filename)
		}
	}
}

// search for SteamGridDB images
func search(cmd *cobra.Command, args []string, kind SearchType) {
	format := rootCmd.PersistentFlags().Lookup("output").Value.String()

	// Ensure we have a SteamGridDB API Key
	apiKey, _ := cmd.Flags().GetString("api-key")
	if apiKey == "" {
		cmd.Help()
		ExitError(fmt.Errorf("API key is required"), format)
	}

	// Create a SteamGridDB Client
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
	maxResults := getFlagInt(cmd, "max-results")
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
		maxImages := getFlagInt(cmd, "max-images")

		// Get all grid images
		if kind.Has(SearchGrids) {
			// Add any requested filters
			filters := []steamgriddb.FilterGrid{}
			if style := cmd.Flags().Lookup("style-grid").Value.String(); style != "" {
				filters = append(filters, steamgriddb.FilterGridStyle(style))
			}

			// Get the grids
			grids, err := client.GetGrids(appID, filters...)
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
			// Add any requested filters
			filters := []steamgriddb.FilterHeroes{}
			if style := cmd.Flags().Lookup("style-hero").Value.String(); style != "" {
				filters = append(filters, steamgriddb.FilterHeroesStyle(style))
			}

			// Get the heroes
			heroes, err := client.GetHeroes(appID, filters...)
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
			// Add any requested filters
			filters := []steamgriddb.FilterLogos{}
			if style := cmd.Flags().Lookup("style-logo").Value.String(); style != "" {
				filters = append(filters, steamgriddb.FilterLogosStyle(style))
			}

			// Get the logos
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
			// Add any requested filters
			filters := []steamgriddb.FilterIcons{}
			if style := cmd.Flags().Lookup("style-icon").Value.String(); style != "" {
				filters = append(filters, steamgriddb.FilterIconsStyle(style))
			}

			// Get the icons
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
		searchFlags := SearchGrids | SearchLogos | SearchIcons | SearchHeroes
		if ok, _ := cmd.Flags().GetBool("only-grids"); ok {
			searchFlags = SearchGrids
		}
		if ok, _ := cmd.Flags().GetBool("only-logos"); ok {
			searchFlags = SearchLogos
		}
		if ok, _ := cmd.Flags().GetBool("only-icons"); ok {
			searchFlags = SearchIcons
		}
		if ok, _ := cmd.Flags().GetBool("only-heroes"); ok {
			searchFlags = SearchHeroes
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
	searchCmd.PersistentFlags().IntP("max-results", "n", 1, "Number of search results to return")
	searchCmd.PersistentFlags().Int("max-images", 1, "Number of image results to return for a given image type")
	searchCmd.PersistentFlags().Bool("only-heroes", false, "Only include hero images in search")
	searchCmd.PersistentFlags().Bool("only-grids", false, "Only include grid images in search")
	searchCmd.PersistentFlags().Bool("only-icons", false, "Only include icon images in search")
	searchCmd.PersistentFlags().Bool("only-logos", false, "Only include logo images in search")
	searchCmd.MarkFlagsMutuallyExclusive("only-grids", "only-heroes", "only-icons", "only-logos")
	searchCmd.Flags().String("style-hero", "", `Optional hero style to search for ("alternate" "blurred" "material")`)
	searchCmd.Flags().String("style-grid", "", `Optional grid style to search for ("alternate" "blurred" "white_logo" "material" "no_logo")`)
	searchCmd.Flags().String("style-icon", "", `Optional icon style to search for ("official" "custom")`)
	searchCmd.Flags().String("style-logo", "", `Optional logo style to search for ("official" "white" "black" "custom")`)
}
