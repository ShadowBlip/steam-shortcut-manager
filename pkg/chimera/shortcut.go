package chimera

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadShortcuts will load all Chimera shortcuts from the given YAML file.
func LoadShortcuts(path string) ([]*Shortcut, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var shortcuts []*Shortcut
	err = yaml.Unmarshal(data, &shortcuts)
	if err != nil {
		return nil, err
	}
	return shortcuts, nil
}

// SaveShortcuts will save the given Chimera shortcuts to the given path.
func SaveShortcuts(path string, shortcuts []*Shortcut) error {
	bytes, err := yaml.Marshal(shortcuts)
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, 0644)
}

// ShortcutSetting is a function that mutates a Chimera Shortcut
type ShortcutSetting func(s *Shortcut)

// DefaultShortcut sets the default settings of a Chimera shortcut
func DefaultShortcut(s *Shortcut) {
	if s.Tags == nil {
		s.Tags = []string{}
	}
	s.Tags = append(s.Tags, "ChimeraOS Playable")
}

// NewShortcut will return a new Chimera Shortcut
func NewShortcut(name, exe string, settings ...ShortcutSetting) *Shortcut {
	shortcut := &Shortcut{Name: name, Cmd: exe}
	for _, setting := range settings {
		setting(shortcut)
	}

	return shortcut
}

// Shortcut is a structure for managing Chimera-managed shortcuts
type Shortcut struct {
	Background string   `yaml:"background,omitempty" json:"background,omitempty"`
	Banner     string   `yaml:"banner,omitempty" json:"banner,omitempty"`
	Cmd        string   `yaml:"cmd" json:"cmd"`
	Dir        string   `yaml:"dir" json:"dir"`
	Hidden     bool     `yaml:"hidden" json:"hidden"`
	Logo       string   `yaml:"logo,omitempty" json:"logo,omitempty"`
	Name       string   `yaml:"name" json:"name"`
	Poster     string   `yaml:"poster,omitempty" json:"poster,omitempty"`
	Tags       []string `yaml:"tags" json:"tags"`
}
