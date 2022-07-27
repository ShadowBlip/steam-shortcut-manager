package steam

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
)

// GetSteamDir will return the base steam config directory
func GetBaseDir() (string, error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return dirname, err
	}

	return path.Join(dirname, ".steam", "steam"), nil
}

// GetSteamUserDir will return the steam userdata directory
func GetUserDir() (string, error) {
	steamDir, err := GetBaseDir()
	if err != nil {
		return steamDir, err
	}

	return path.Join(steamDir, "userdata"), nil
}

// GetUsers will return a list of steam user ids
func GetUsers() ([]string, error) {
	userDir, err := GetUserDir()
	if err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(userDir)
	if err != nil {
		return nil, err
	}

	users := []string{}
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		users = append(users, f.Name())
	}

	return users, nil
}

// GetImagesDir will return the steam images directory
func GetImagesDir(user string) (string, error) {
	userDir, err := GetUserDir()
	if err != nil {
		return "", err
	}

	return path.Join(userDir, user, "config", "grid"), nil
}

// GetShortcutsPath will return the path to the shortcuts file for the given
// user.
func GetShortcutsPath(user string) (string, error) {
	userDir, err := GetUserDir()
	if err != nil {
		return "", err
	}

	return path.Join(userDir, user, "config", "shortcuts.vdf"), nil
}

// Whether or not the user has a shortcuts file
func HasShortcuts(user string) bool {
	shortcutsPath, err := GetShortcutsPath(user)
	if err != nil {
		return false
	}
	if _, err := os.Stat(shortcutsPath); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
