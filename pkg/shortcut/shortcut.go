package shortcut

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/wakeful-cloud/vdf"
)

// Load the given shortcuts file
func Load(file string) (*Shortcuts, error) {
	// Read the VDF bytes
	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// Parse the VDF file
	vdfMap, err := vdf.ReadVdf(bytes)
	if err != nil {
		return nil, err
	}

	// Covert to JSON so we can map it to a struct
	rawJSON, err := json.Marshal(vdfMap)
	if err != nil {
		return nil, err
	}

	// Unmarshal to a struct
	var shortcuts Shortcuts
	err = json.Unmarshal(rawJSON, &shortcuts)
	if err != nil {
		return nil, err
	}

	return &shortcuts, nil
}

// Save the given shortcuts file
func Save(shortcuts *Shortcuts, file string) error {
	// Convert the struct to JSON so we can map it to a VDF map
	rawJSON, err := json.Marshal(shortcuts)
	if err != nil {
		return fmt.Errorf("Unable to marshal to JSON: %v", err)
	}

	// Marshal the shortcut into a VDF map
	var vdfMap map[string]interface{}
	err = json.Unmarshal(rawJSON, &vdfMap)
	if err != nil {
		return fmt.Errorf("Unable to unmarshal to VDF Map: %v", err)
	}

	// Save the shortcuts
	rawVdf, err := vdf.WriteVdf(ensureVDFMap(vdfMap))
	if err != nil {
		return fmt.Errorf("Unable to convert VDF to bytes: %v", err)
	}

	// Write the file
	err = os.WriteFile(file, rawVdf, 0666)
	if err != nil {
		return fmt.Errorf("Unable to write VDF file: %v", err)
	}

	return nil
}

// Ensures the given map is a vdf.Map
func ensureVDFMap(m map[string]interface{}) vdf.Map {
	var newMap vdf.Map = vdf.Map{}
	for k, v := range m {
		switch v.(type) {
		case int, int64:
			newMap[k] = v.(uint32)
		case float64:
			newMap[k] = uint32(v.(float64))
		case map[string]interface{}:
			newMap[k] = ensureVDFMap(v.(map[string]interface{}))
		default:
			newMap[k] = v
		}
	}
	return newMap
}
