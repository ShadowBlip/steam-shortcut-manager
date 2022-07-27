package steamgriddb

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

const BASE_URL = "https://www.steamgriddb.com/api/v2"

var isDebug = os.Getenv("DEBUG") == "1"

// NewClient will return a new SteamGridDB Client
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
	}
}

// Client is a structure for querying the SteamGridDB API
type Client struct {
	apiKey string
	client http.Client
}

func (c *Client) debug(str string) {
	if !isDebug {
		return
	}
	fmt.Printf("%s\n", str)
}

// Get will perform a GET request to the given SteamGridDB API endpoint.
func (c *Client) Get(path string) (*http.Response, error) {
	return c.get(getUrl(path))
}

func (c *Client) get(url string) (*http.Response, error) {
	c.debug("GET " + url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Received non 200 response code")
	}
	return res, nil
}

// Download will download the given file to the provided path
func (c *Client) Download(url, path string) error {
	// Fetch the file
	res, err := c.get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Create a empty file
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the bytes to the file
	_, err = io.Copy(file, res.Body)
	if err != nil {
		return err
	}

	return nil
}

// CachedDownload will download only if the file does not already exist.
func (c *Client) CachedDownload(url, path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return c.Download(url, path)
	}
	return nil
}

// Search will return a list of search results for the given term
func (c *Client) Search(term string) (*SearchResponse, error) {
	res, err := c.Get("/search/autocomplete/" + url.QueryEscape(term))
	if err != nil {
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var results SearchResponse
	err = json.Unmarshal(body, &results)
	if err != nil {
		return nil, err
	}

	return &results, nil
}

// GetGrids will return the results of the grids for a given game ID
func (c *Client) GetGrids(gameID string) (*GridResponse, error) {
	res, err := c.Get("/grids/game/" + gameID)
	if err != nil {
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var results GridResponse
	err = json.Unmarshal(body, &results)
	if err != nil {
		return nil, err
	}

	return &results, nil
}

// GetHeroes will return the results of heroes for a given game ID
func (c *Client) GetHeroes(gameID string) (*HeroesResponse, error) {
	res, err := c.Get("/heroes/game/" + gameID)
	if err != nil {
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var results HeroesResponse
	err = json.Unmarshal(body, &results)
	if err != nil {
		return nil, err
	}

	return &results, nil
}

// GetLogos will return the results of logos for a given game ID
func (c *Client) GetLogos(gameID string) (*LogosResponse, error) {
	res, err := c.Get("/logos/game/" + gameID)
	if err != nil {
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var results LogosResponse
	err = json.Unmarshal(body, &results)
	if err != nil {
		return nil, err
	}

	return &results, nil
}

// GetIcons will return the results of icons for a given game ID
func (c *Client) GetIcons(gameID string) (*IconsResponse, error) {
	res, err := c.Get("/icons/game/" + gameID)
	if err != nil {
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var results IconsResponse
	err = json.Unmarshal(body, &results)
	if err != nil {
		return nil, err
	}

	return &results, nil
}

func getUrl(path string) string {
	return fmt.Sprintf("%s%s", BASE_URL, path)
}
