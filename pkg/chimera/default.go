package chimera

import (
	"fmt"
	"os"
	"path"
)

var homeDir, _ = os.UserHomeDir()

var ConfigDir = path.Join(homeDir, ".local/share/chimera")
var ShortcutsDir = path.Join(ConfigDir, "shortcuts")
var ImagesDir = path.Join(ConfigDir, "images")

var BannerDir = path.Join(ImagesDir, "banner")
var LogoDir = path.Join(ImagesDir, "logo")
var PosterDir = path.Join(ImagesDir, "poster")
var BackgroundDir = path.Join(ImagesDir, "background")

var SupportedPlatforms = []string{"flathub"}

// HasChimera will return whether or not Chimera has a configuration directory
func HasChimera() bool {
	if _, err := os.Stat(ConfigDir); !os.IsNotExist(err) {
		return true
	}
	return false
}

// EnsureShortcutsFileExists will create the Chimera shortcuts file for the
// given platform if it does not already exist.
func EnsureShortcutsFileExists(platform string) error {
	fileName := GetShortcutsFile(platform)
	_, err := os.Stat(fileName)
	if !os.IsNotExist(err) {
		return nil
	}
	err = os.MkdirAll(path.Dir(fileName), 0755)
	if err != nil {
		return err
	}
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

// GetShortcutsFile will return the path to the shortcuts file for the given
// platform (e.g. flathub)
func GetShortcutsFile(platform string) string {
	return path.Join(ShortcutsDir, fmt.Sprintf("chimera.%s.yaml", platform))
}

// IsPlatformSupported will return whether or not the given Chimera platform
// is supported by the shortcut manager.
func IsPlatformSupported(platform string) bool {
	for _, plat := range SupportedPlatforms {
		if plat == platform {
			return true
		}
	}
	return false
}
