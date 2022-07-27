package image

import (
	"os"

	"github.com/shadowblip/steam-shortcut-manager/pkg/image/kitty"
)

// Displayer is a function signature of a provider that can display images
type Displayer func(filename string) error

// Display the given image using a detected provider
var Display Displayer
var CanDisplay = false

func init() {
	// Set our displayer to Kitty if detected
	if os.Getenv("TERM") == "xterm-kitty" {
		Display = kitty.Display
		CanDisplay = true
	}
}
