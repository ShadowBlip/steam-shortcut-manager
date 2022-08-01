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

	"github.com/shadowblip/steam-shortcut-manager/pkg/chimera"
	"github.com/shadowblip/steam-shortcut-manager/pkg/shortcut"
	"github.com/shadowblip/steam-shortcut-manager/pkg/steam"
	"github.com/shadowblip/steam-shortcut-manager/pkg/steamgriddb"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <name> <exe>",
	Short: "Add a Steam shortcut to your steam library",
	Args:  cobra.ExactArgs(2),
	Long:  `Adds a Steam shortcut to your library`,
	Run: func(cmd *cobra.Command, args []string) {
		format := rootCmd.PersistentFlags().Lookup("output").Value.String()
		name := args[0]
		exe := args[1]

		// Fetch all users
		users, err := steam.GetUsers()
		if err != nil {
			ExitError(err, format)
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
				ExitError(err, format)
			}

			// Generate a new shortcut from the cli flags
			newShortcut := newShortcutFromFlags(cmd, name, exe)
			// Download images for the user if specified
			if download, _ := cmd.Flags().GetBool("download-images"); download {
				// Check that we have an API key
				apiKey, _ := cmd.Flags().GetString("api-key")
				if apiKey == "" {
					ExitError(fmt.Errorf("no API key specified"), format)
				}
				client := steamgriddb.NewClient(apiKey)
				downloaded, err := downloadImages(client, user, newShortcut)
				if err != nil {
					DebugPrintln("Error downloading images:", err)
					//ExitError(err, format)
				}

				// Update our shortcut with image paths if needed
				for imgType, path := range downloaded {
					switch imgType {
					case "icon":
						newShortcut.Icon = path
					}
				}
			}

			// Write the changes
			shortcuts.Add(newShortcut)
			err = shortcut.Save(shortcuts, shortcutsPath)
			if err != nil {
				ExitError(err, format)
			}
		}
	},
}

// downloadImages will download images for the given shortcut
// TODO: Handle errors better
func downloadImages(client *steamgriddb.Client, user string, sc *shortcut.Shortcut) (map[string]string, error) {
	// This map will contain the paths to our downloaded images
	downloaded := map[string]string{}

	// Get the image directory for the user.
	gridDir, err := steam.GetImagesDir(user)
	if err != nil {
		return nil, err
	}

	// Search for the app images
	results, err := client.Search(sc.AppName)
	if err != nil {
		return nil, err
	}
	// TODO: Log or return no image results
	if len(results.Data) == 0 {
		return nil, err
	}

	// Get the first result
	// TODO: Enable showing different results?
	gameID := fmt.Sprintf("%v", results.Data[0].ID)
	steamAppID := fmt.Sprintf("%v", sc.Appid)

	// Download the grid images. Steam uses a portrait and landscape image
	// that is displays in the library.
	grids, err := client.GetGrids(gameID)
	if err != nil {
		//return nil, err
		grids = &steamgriddb.GridResponse{Data: []steamgriddb.GridResponseData{}}
	}
	portraitGrids := steamgriddb.FilterGridVertical()(grids)
	for _, data := range portraitGrids {
		ext := filepath.Ext(data.URL)
		imgFile := path.Join(gridDir, fmt.Sprintf("%sp%s", steamAppID, ext))
		err := client.CachedDownload(data.URL, imgFile)
		if err != nil {
			continue
		}
		downloaded["gridP"] = imgFile
		break
	}
	landscapeGrids := steamgriddb.FilterGridHorizontal()(grids)
	for _, data := range landscapeGrids {
		ext := filepath.Ext(data.URL)
		imgFile := path.Join(gridDir, fmt.Sprintf("%s%s", steamAppID, ext))
		err := client.CachedDownload(data.URL, imgFile)
		if err != nil {
			continue
		}
		downloaded["gridL"] = imgFile
		break
	}

	// Download the hero image. The hero image is used as a banner at the
	// top of the app page in the Steam UI.
	heroes, err := client.GetHeroes(gameID)
	if err != nil {
		//return nil, err
		heroes = &steamgriddb.HeroesResponse{Data: []steamgriddb.ImageResponseData{}}
	}
	for _, data := range heroes.Data {
		ext := filepath.Ext(data.URL)
		imgFile := path.Join(gridDir, fmt.Sprintf("%s_hero%s", steamAppID, ext))
		err := client.CachedDownload(data.URL, imgFile)
		if err != nil {
			continue
		}
		downloaded["hero"] = imgFile
		break
	}

	// Download the logo image. Logo images are used in the Steam overlay menu.
	logos, err := client.GetLogos(gameID)
	if err != nil {
		//return nil, err
		logos = &steamgriddb.LogosResponse{Data: []steamgriddb.ImageResponseData{}}
	}
	for _, data := range logos.Data {
		ext := filepath.Ext(data.URL)
		imgFile := path.Join(gridDir, fmt.Sprintf("%s_logo%s", steamAppID, ext))
		err := client.CachedDownload(data.URL, imgFile)
		if err != nil {
			continue
		}
		downloaded["logo"] = imgFile
		break
	}

	// Download the icon image. Icon images are used in some part of the UI.
	icons, err := client.GetIcons(gameID)
	if err != nil {
		//return nil, err
		icons = &steamgriddb.IconsResponse{Data: []steamgriddb.ImageResponseData{}}
	}
	for _, data := range icons.Data {
		ext := filepath.Ext(data.URL)
		imgFile := path.Join(gridDir, fmt.Sprintf("%s-icon%s", steamAppID, ext))
		err := client.CachedDownload(data.URL, imgFile)
		if err != nil {
			continue
		}
		downloaded["icon"] = imgFile
		break
	}

	return downloaded, nil
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

// chimeraAddCmd represents the add command
var chimeraAddCmd = &cobra.Command{
	Use:   "add <name> <exe>",
	Short: "Add a Chimera shortcut to your steam library",
	Args:  cobra.ExactArgs(2),
	Long:  `Adds a Chimera shortcut to your library`,
	Run: func(cmd *cobra.Command, args []string) {
		format := rootCmd.PersistentFlags().Lookup("output").Value.String()
		name := args[0]
		exe := args[1]

		// Ensure we have a Chimera install
		if !chimera.HasChimera() {
			ExitError(fmt.Errorf("no chimera config found at %v", chimera.ConfigDir), format)
		}

		// Get the platform flag
		platform := chimeraCmd.PersistentFlags().Lookup("platform").Value.String()
		DebugPrintln("Using Chimera platform:", platform)

		// Ensure the Chimera shortcuts file exists
		err := chimera.EnsureShortcutsFileExists(platform)
		if err != nil {
			ExitError(err, format)
		}

		// Check that we have required params for platform
		switch platform {
		case "flathub":
			if id, _ := cmd.Flags().GetString("flatpak-id"); id == "" {
				ExitError(fmt.Errorf("flatpak-id required for flathub platform"), format)
			}
		}

		// Read from the given shortcuts file
		shortcutsFile := chimera.GetShortcutsFile(platform)
		DebugPrintln("Using shortcuts file:", shortcutsFile)
		shortcuts, err := chimera.LoadShortcuts(shortcutsFile)
		if err != nil {
			ExitError(err, format)
		}

		// Create the new shortcut to add
		newShortcut := newChimeraShortcutFromFlags(cmd, name, exe)

		// Download images for the user if specified
		if download, _ := cmd.Flags().GetBool("download-images"); download {
			DebugPrintln("Requested to download images for shortcut")
			// Check that we have an API key
			apiKey, _ := cmd.Flags().GetString("api-key")
			if apiKey == "" {
				ExitError(fmt.Errorf("no API key specified"), format)
			}

			// Download the images
			client := steamgriddb.NewClient(apiKey)
			downloaded, err := downloadChimeraImages(cmd.Flags(), client, platform, newShortcut)
			if err != nil {
				ExitError(err, format)
			}

			// Update our shortcut with image paths
			for imgType, path := range downloaded {
				switch imgType {
				case "poster":
					newShortcut.Poster = path
				case "background":
					newShortcut.Background = path
				case "banner":
					newShortcut.Banner = path
				case "logo":
					newShortcut.Logo = path
				}
			}
		}

		// Save the shortcuts
		shortcuts = append(shortcuts, newShortcut)
		err = chimera.SaveShortcuts(shortcutsFile, shortcuts)
		if err != nil {
			ExitError(err, format)
		}

		// Print the output
		switch format {
		case "term":
			fmt.Println(newShortcut.Name)
			fmt.Println("  Executable:", newShortcut.Cmd)
			fmt.Println("  Poster:", newShortcut.Poster)
			fmt.Println("  Banner:", newShortcut.Banner)
			fmt.Println("  Logo:", newShortcut.Logo)
			fmt.Println("  Background:", newShortcut.Background)
		case "json":
			out, err := json.MarshalIndent(newShortcut, "", "  ")
			if err != nil {
				ExitError(err, format)
			}
			fmt.Println(string(out))
		default:
			panic("unknown output format: " + format)
		}
	},
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

// Creates a new Chimera shortcut entry from command-line flags
func newChimeraShortcutFromFlags(cmd *cobra.Command, name, exe string) *chimera.Shortcut {
	getString := func(name string) string {
		res, _ := cmd.Flags().GetString(name)
		return res
	}
	shortcutConfiger := func(s *chimera.Shortcut) {
		s.Dir = getString("start-dir")
		s.Hidden, _ = cmd.Flags().GetBool("is-hidden")

		s.Tags = []string{}
		tags, _ := cmd.Flags().GetStringSlice("tags")
		for _, tag := range tags {
			s.Tags = append(s.Tags, tag)
		}
	}
	shortcut := chimera.NewShortcut(name, exe, shortcutConfiger)
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
	chimeraCmd.AddCommand(chimeraAddCmd)

	// Normal add flags
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
	addCmd.Flags().StringP("chimera-shortcut", "c", "~/.local/share/chimera/shortcuts/chimera.flathub.yaml", "Optional path to Chimera shortcut config")

	addCmd.Flags().StringP("api-key", "k", "", "SteamGridDB API Key")
	addCmd.Flags().BoolP("download-images", "i", false, "Auto-download artwork from SteamGridDB for shortcut (requires SteamGridDB API Key)")

	// Chimera add flags
	chimeraAddCmd.Flags().String("start-dir", "~", "Working directory where the app is started")
	chimeraAddCmd.Flags().Bool("is-hidden", false, "Whether or not the shortcut is hidden")
	chimeraAddCmd.Flags().StringSlice("tags", []string{}, "Comma-separated list of tags")
	chimeraAddCmd.Flags().String("flatpak-id", "", "Flatpak ID of the shortcut (if platform 'flathub')")

	chimeraAddCmd.Flags().StringP("api-key", "k", "", "SteamGridDB API Key")
	chimeraAddCmd.Flags().BoolP("download-images", "i", false, "Auto-download artwork from SteamGridDB for shortcut (requires SteamGridDB API Key)")
}
