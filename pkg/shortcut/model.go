package shortcut

import (
	"fmt"
	"strconv"
)

/*
  "shortcuts": {
    "0": {
      "AllowDesktopConfig": 1,
      "AllowOverlay": 1,
      "AppName": "Insomnia",
      "Devkit": 0,
      "DevkitGameID": "",
      "DevkitOverrideAppID": 0,
      "Exe": "\"/usr/bin/flatpak\"",
      "FlatpakAppID": "",
      "IsHidden": 0,
      "LastPlayTime": 0,
      "LaunchOptions": "run --branch=stable --arch=x86_64 --command=/app/bin/insomnia --file-forwarding rest.insomnia.Insomnia --no-sandbox @@u @@",
      "OpenVR": 0,
      "ShortcutPath": "/var/lib/flatpak/exports/share/applications/rest.insomnia.Insomnia.desktop",
      "StartDir": "\"/usr/bin/\"",
      "appid": 3417544970,
      "icon": "",
      "tags": {}
    }
  }
}
*/

type Shortcuts struct {
	Shortcuts map[string]Shortcut `json:"shortcuts"`
}

// Add will add the given shortcut
func (s *Shortcuts) Add(shortcut *Shortcut) error {
	nextKey, err := s.getNextKey()
	if err != nil {
		return err
	}
	s.Shortcuts[nextKey] = *shortcut

	return nil
}

// Get the next shortcut id
func (s *Shortcuts) getNextKey() (string, error) {
	highestKey := -1
	for key := range s.Shortcuts {
		keyNum, err := strconv.Atoi(key)
		if err != nil {
			return "", fmt.Errorf("Non-number shortcut key: %v", err)
		}
		if keyNum > highestKey {
			highestKey = keyNum
		}
	}

	return fmt.Sprintf("%v", highestKey+1), nil
}

// ShortcutSetting is a function that mutates a Shortcut
type ShortcutSetting func(s *Shortcut)

// DefaultShortcut sets the default settings of a shortcut
func DefaultShortcut(s *Shortcut) {
	s.AllowDesktopConfig = 1
	s.AllowOverlay = 1
}

// NewShortcut will return a new Steam Shortcut
func NewShortcut(name, exe string, settings ...ShortcutSetting) *Shortcut {
	shortcut := &Shortcut{AppName: name, Exe: exe}
	for _, setting := range settings {
		setting(shortcut)
	}

	return shortcut
}

type Shortcut struct {
	AllowDesktopConfig  int                    `json:"AllowDesktopConfig"`
	AllowOverlay        int                    `json:"AllowOverlay"`
	AppName             string                 `json:"AppName"`
	Devkit              int                    `json:"Devkit"`
	DevkitGameID        string                 `json:"DevkitGameID"`
	DevkitOverrideAppID int                    `json:"DevkitOverrideAppID"`
	Exe                 string                 `json:"Exe"`
	FlatpakAppID        string                 `json:"FlatpakAppID"`
	IsHidden            int                    `json:"IsHidden"`
	LastPlayTime        int                    `json:"LastPlayTime"`
	LaunchOptions       string                 `json:"LaunchOptions"`
	OpenVR              int                    `json:"OpenVR"`
	ShortcutPath        string                 `json:"ShortcutPath"`
	StartDir            string                 `json:"StartDir"`
	Appid               int64                  `json:"appid"`
	Icon                string                 `json:"icon"`
	Tags                map[string]interface{} `json:"tags"`
}
