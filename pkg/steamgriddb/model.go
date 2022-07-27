package steamgriddb

// Response is a generic SteamGridDB response
type Response struct {
	Success bool     `json:"success"`
	Errors  []string `json:"errors"`
}

// 'https://www.steamgriddb.com/api/v2/search/autocomplete/{term}'
type SearchResponse struct {
	Response
	Data []SearchResponseData `json:"data"`
}

type SearchResponseData struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Types    []string `json:"types"`
	Verified bool     `json:"verified"`
}

// https://www.steamgriddb.com/api/v2/grids/game/{gameId}
type GridResponse struct {
	Response
	Data []GridResponseData `json:"data"`
}

type GridResponseData struct {
	ID        int         `json:"id"`
	Score     int         `json:"score"`
	Style     string      `json:"style"`
	Width     int         `json:"width"`
	Height    int         `json:"height"`
	Nsfw      bool        `json:"nsfw"`
	Humor     bool        `json:"humor"`
	Notes     interface{} `json:"notes"`
	Mime      string      `json:"mime"`
	Language  string      `json:"language"`
	URL       string      `json:"url"`
	Thumb     string      `json:"thumb"`
	Lock      bool        `json:"lock"`
	Epilepsy  bool        `json:"epilepsy"`
	Upvotes   int         `json:"upvotes"`
	Downvotes int         `json:"downvotes"`
	Author    struct {
		Name    string `json:"name"`
		Steam64 string `json:"steam64"`
		Avatar  string `json:"avatar"`
	} `json:"author"`
}

// https://www.steamgriddb.com/api/v2/heroes/game/{gameId}
type HeroesResponse struct {
	Response
	Data []ImageResponseData `json:"data"`
}

type ImageResponseData struct {
	ID     int      `json:"id"`
	Score  int      `json:"score"`
	Style  string   `json:"style"`
	URL    string   `json:"url"`
	Thumb  string   `json:"thumb"`
	Tags   []string `json:"tags"`
	Author struct {
		Name    string `json:"name"`
		Steam64 string `json:"steam64"`
		Avatar  string `json:"avatar"`
	} `json:"author"`
}

// https://www.steamgriddb.com/api/v2/logos/game/{gameId}
type LogosResponse HeroesResponse

// https://www.steamgriddb.com/api/v2/icons/game/{gameId}
type IconsResponse HeroesResponse
