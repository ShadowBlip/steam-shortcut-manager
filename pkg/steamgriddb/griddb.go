package steamgriddb

import (
	"encoding/json"
	"fmt"
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
	c.debug("GET " + getUrl(path))
	req, err := http.NewRequest("GET", getUrl(path), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	return c.client.Do(req)
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
