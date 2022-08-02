package steam

import (
	"errors"
	"fmt"
	"os"
	"path"
)

// ErrImageNotFound indicates that a grid images does not exist.
var ErrImageNotFound = errors.New("image not found")

// GetImagesDir will return the steam images directory
func GetImagesDir(user string) (string, error) {
	userDir, err := GetUserDir()
	if err != nil {
		return "", err
	}
	return path.Join(userDir, user, "config", "grid"), nil
}

// GetImageLandscape will return the landscape grid image
func GetImageLandscape(user, appId string) (string, error) {
	imagesDir, err := GetImagesDir(user)
	if err != nil {
		return "", err
	}

	// Check to see if the file exists with different extensions
	return checkForImage(path.Join(imagesDir, appId))
}

// GetImagePortrait will return the portrait grid image
func GetImagePortrait(user, appId string) (string, error) {
	imagesDir, err := GetImagesDir(user)
	if err != nil {
		return "", err
	}

	// Check to see if the file exists with different extensions
	return checkForImage(path.Join(imagesDir, fmt.Sprintf("%sp", appId)))
}

// GetImageHero will return the hero grid image
func GetImageHero(user, appId string) (string, error) {
	imagesDir, err := GetImagesDir(user)
	if err != nil {
		return "", err
	}

	// Check to see if the file exists with different extensions
	return checkForImage(path.Join(imagesDir, fmt.Sprintf("%s_hero", appId)))
}

// GetImageLogo will return the logo grid image
func GetImageLogo(user, appId string) (string, error) {
	imagesDir, err := GetImagesDir(user)
	if err != nil {
		return "", err
	}

	// Check to see if the file exists with different extensions
	return checkForImage(path.Join(imagesDir, fmt.Sprintf("%s_logo", appId)))
}

// checkForImage will check various image extensions for the given file path
// without an extension. Returns a ErrImageNotFound error if it does not exist.
func checkForImage(basePath string) (string, error) {
	knownExtensions := []string{"png", "jpg", "jpeg", "ico"}
	for _, ext := range knownExtensions {
		fileName := fmt.Sprintf("%s.%s", basePath, ext)
		if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
			continue
		}
		return fileName, nil
	}
	return "", ErrImageNotFound
}
