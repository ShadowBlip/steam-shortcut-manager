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
	"path"
	"path/filepath"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/shadowblip/steam-shortcut-manager/pkg/chimera"
	"github.com/shadowblip/steam-shortcut-manager/pkg/shortcut"
	"github.com/shadowblip/steam-shortcut-manager/pkg/steam"
	"github.com/shadowblip/steam-shortcut-manager/pkg/steamgriddb"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func discoverDownloadDir(cmd *cobra.Command, format string) string {
	// Handle when download dir is directly specified
	destinationDir, _ := cmd.PersistentFlags().GetString("destination-dir")
	if destinationDir != "" {
		return destinationDir
	}

	// Handle downloading for steam user IDs
	user, _ := cmd.PersistentFlags().GetString("user")
	if user != "" {
		users, err := steam.GetUsers()
		if err != nil {
			ExitError(err, format)
		}
		if !contains(users, user) {
			ExitError(fmt.Errorf("user not found"), format)
		}
		downloadDir, err := steam.GetImagesDir(user)
		if err != nil {
			ExitError(err, format)
		}
		return downloadDir
	}

	ExitError(fmt.Errorf("unable to discover download directory"), format)
	return ""
}

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download --api-key=<key> <name>",
	Short: "Download SteamGridDB images for a given app",
	Long:  `Download SteamGridDB images for a given app`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		format := rootCmd.PersistentFlags().Lookup("output").Value.String()
		name := args[0]

		// Ensure we have a SteamGridDB API Key
		apiKey, _ := cmd.Flags().GetString("api-key")
		if apiKey == "" {
			cmd.Help()
			ExitError(fmt.Errorf("API key is required"), format)
		}
		appId, _ := cmd.Flags().GetInt("app-id")
		if appId == 0 {
			cmd.Help()
			ExitError(fmt.Errorf("Shortcut app id is required"), format)
		}

		// Create a shortcut used to download the files
		sc := &shortcut.Shortcut{Appid: int64(appId), AppName: name}

		// Create a SteamGridDB client
		client := steamgriddb.NewClient(apiKey)

		// Get all steam users
		users, err := steam.GetUsers()
		if err != nil {
			ExitError(err, format)
		}

		// TODO: Cache and symlink instead of downloading for each user
		results := map[string]interface{}{}
		var errors error
		for _, user := range users {
			// Check that we have an API key
			apiKey, _ := cmd.Flags().GetString("api-key")
			if apiKey == "" {
				ExitError(fmt.Errorf("no API key specified"), format)
			}
			DebugPrintln("Downloading images for shortcut")
			downloaded, err := downloadImages(client, user, sc)
			if err != nil {
				DebugPrintln("Error downloading images:", err)
				errors = multierror.Append(errors, err)
			}
			results[user] = downloaded
		}
		if errors != nil {
			ExitError(err, format)
		}

		// Print the output
		switch format {
		case "term":
			fmt.Println(results)
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

// downloadImages will download images for the given shortcut
// TODO: Handle errors better
func downloadImages(client *steamgriddb.Client, user string, sc *shortcut.Shortcut) (map[string]string, error) {
	// This map will contain the paths to our downloaded images
	downloaded := map[string]string{}
	var errors error

	// Get the image directory for the user.
	gridDir, err := steam.GetImagesDir(user)
	if err != nil {
		return nil, err
	}
	DebugPrintln("Discovered images dir:", gridDir)

	// Search for the app images
	results, err := client.Search(sc.AppName)
	if err != nil {
		return nil, err
	}
	// TODO: Log or return no image results
	if len(results.Data) == 0 {
		return nil, fmt.Errorf("no results found for %v", sc.AppName)
	}
	DebugPrintln(fmt.Sprintf("Found %v results for %s", len(results.Data), sc.AppName))

	// Get the first result
	// TODO: Enable showing different results?
	gameID := fmt.Sprintf("%v", results.Data[0].ID)
	steamAppID := fmt.Sprintf("%v", sc.Appid)

	// Download the grid images. Steam uses a portrait and landscape image
	// that is displays in the library.
	grids, err := client.GetGrids(gameID)
	if err != nil {
		errors = multierror.Append(errors, err)
		grids = &steamgriddb.GridResponse{Data: []steamgriddb.GridResponseData{}}
	}
	portraitGrids := steamgriddb.FilterGridVertical()(grids)
	for _, data := range portraitGrids {
		ext := filepath.Ext(data.URL)
		imgFile := path.Join(gridDir, fmt.Sprintf("%sp%s", steamAppID, ext))
		DebugPrintln("Downloading portrait grid image...")
		err := client.CachedDownload(data.URL, imgFile)
		if err != nil {
			errors = multierror.Append(errors, err)
			continue
		}
		downloaded["gridP"] = imgFile
		break
	}
	landscapeGrids := steamgriddb.FilterGridHorizontal()(grids)
	for _, data := range landscapeGrids {
		ext := filepath.Ext(data.URL)
		imgFile := path.Join(gridDir, fmt.Sprintf("%s%s", steamAppID, ext))
		DebugPrintln("Downloading landscape grid image...")
		err := client.CachedDownload(data.URL, imgFile)
		if err != nil {
			errors = multierror.Append(errors, err)
			continue
		}
		downloaded["gridL"] = imgFile
		break
	}

	// Download the hero image. The hero image is used as a banner at the
	// top of the app page in the Steam UI.
	heroes, err := client.GetHeroes(gameID)
	if err != nil {
		errors = multierror.Append(errors, err)
		heroes = &steamgriddb.HeroesResponse{Data: []steamgriddb.ImageResponseData{}}
	}
	for _, data := range heroes.Data {
		ext := filepath.Ext(data.URL)
		imgFile := path.Join(gridDir, fmt.Sprintf("%s_hero%s", steamAppID, ext))
		DebugPrintln("Downloading hero grid image...")
		err := client.CachedDownload(data.URL, imgFile)
		if err != nil {
			errors = multierror.Append(errors, err)
			continue
		}
		downloaded["hero"] = imgFile
		break
	}

	// Download the logo image. Logo images are used in the Steam overlay menu.
	logos, err := client.GetLogos(gameID)
	if err != nil {
		errors = multierror.Append(errors, err)
		logos = &steamgriddb.LogosResponse{Data: []steamgriddb.ImageResponseData{}}
	}
	for _, data := range logos.Data {
		ext := filepath.Ext(data.URL)
		imgFile := path.Join(gridDir, fmt.Sprintf("%s_logo%s", steamAppID, ext))
		DebugPrintln("Downloading logo grid image...")
		err := client.CachedDownload(data.URL, imgFile)
		if err != nil {
			errors = multierror.Append(errors, err)
			continue
		}
		downloaded["logo"] = imgFile
		break
	}

	// Download the icon image. Icon images are used in some part of the UI.
	icons, err := client.GetIcons(gameID)
	if err != nil {
		errors = multierror.Append(errors, err)
		icons = &steamgriddb.IconsResponse{Data: []steamgriddb.ImageResponseData{}}
	}
	for _, data := range icons.Data {
		ext := filepath.Ext(data.URL)
		imgFile := path.Join(gridDir, fmt.Sprintf("%s-icon%s", steamAppID, ext))
		err := client.CachedDownload(data.URL, imgFile)
		if err != nil {
			errors = multierror.Append(errors, err)
			continue
		}
		downloaded["icon"] = imgFile
		break
	}

	return downloaded, errors
}

// downloadChimeraImages will download images for the given shortcut. This
// will return the paths of each type of image we downloaded.
// TODO: Handle errors better
func downloadChimeraImages(flags *pflag.FlagSet, client *steamgriddb.Client, platform string, sc *chimera.Shortcut) (map[string]string, error) {
	// Get the download directories
	posterDir := path.Join(chimera.PosterDir, platform)
	bannerDir := path.Join(chimera.BannerDir, platform)
	logoDir := path.Join(chimera.LogoDir, platform)
	backgroundDir := path.Join(chimera.BackgroundDir, platform)

	// This map will contain the paths to our downloaded images
	downloaded := map[string]string{}

	// Set the base file name based on platform
	var fileBaseName string

	// Shortcuts for different platforms are handled differently
	switch platform {
	case "flathub":
		fileBaseName = flags.Lookup("flatpak-id").Value.String()
	default:
		fileBaseName = sc.Name
	}

	// Search for the app images
	results, err := client.Search(sc.Name)
	if err != nil {
		return downloaded, nil
	}
	// TODO: Log or return no image results
	if len(results.Data) == 0 {
		return downloaded, nil
	}

	// Get the first result
	// TODO: Enable showing different results?
	gameID := fmt.Sprintf("%v", results.Data[0].ID)

	// Download the grid image. Grid images are "poster" Chimera images
	grids, err := client.GetGrids(gameID)
	if err != nil {
		return nil, err
	}
	posters := steamgriddb.FilterGridVertical()(grids)
	for _, data := range posters {
		ext := filepath.Ext(data.URL)
		imgFile := path.Join(posterDir, fmt.Sprintf("%s%s", fileBaseName, ext))
		err := client.CachedDownload(data.URL, imgFile)
		if err != nil {
			continue
		}
		downloaded["poster"] = imgFile
		break
	}

	// We also need to get the "banner" image which is a grid with
	// a landscape orientation.
	banners := steamgriddb.FilterGridHorizontal()(grids)
	for _, data := range banners {
		ext := filepath.Ext(data.URL)
		imgFile := path.Join(bannerDir, fmt.Sprintf("%s%s", fileBaseName, ext))
		err := client.CachedDownload(data.URL, imgFile)
		if err != nil {
			continue
		}
		downloaded["banner"] = imgFile
		break
	}

	// Download the hero image. Hero images are "background" Chimera images
	heroes, err := client.GetHeroes(gameID)
	if err != nil {
		return nil, err
	}
	for _, data := range heroes.Data {
		ext := filepath.Ext(data.URL)
		imgFile := path.Join(backgroundDir, fmt.Sprintf("%s%s", fileBaseName, ext))
		err := client.CachedDownload(data.URL, imgFile)
		if err != nil {
			continue
		}
		downloaded["background"] = imgFile
		break
	}

	// Download the logo image. Logo images are "logo" Chimera images
	logos, err := client.GetLogos(gameID)
	if err != nil {
		return nil, err
	}
	for _, data := range logos.Data {
		ext := filepath.Ext(data.URL)
		imgFile := path.Join(logoDir, fmt.Sprintf("%s%s", fileBaseName, ext))
		err := client.CachedDownload(data.URL, imgFile)
		if err != nil {
			continue
		}
		downloaded["logo"] = imgFile
		break
	}

	return downloaded, nil
}

func init() {
	steamgriddbCmd.AddCommand(downloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	downloadCmd.Flags().IntP("app-id", "i", 0, "Steam App ID to download images for")

	downloadCmd.PersistentFlags().Bool("only-hero", false, "Only download hero image")
	downloadCmd.PersistentFlags().Bool("only-grid", false, "Only download grid image")
	downloadCmd.PersistentFlags().Bool("only-icon", false, "Only download icon image")
	downloadCmd.PersistentFlags().Bool("only-logo", false, "Only download logo image")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
