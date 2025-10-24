package reddit



type RedditAPIResponse struct {
	Data struct {
		After    string       `json:"after"`
		Before   string       `json:"before"`
		Children []struct {
			Data RedditPost `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type RedditPost struct {
	ID        string `json:"id"`
	Post      string `json:"permalink"`
	Preview   Preview `json:"preview,omitempty"`
	URL       string `json:"url"`
	Score     int    `json:"score"`
	Height    int    `json:"height,omitempty"`
	Width     int    `json:"width,omitempty"`
	Is18      bool   `json:"over_18"`
	Title     string `json:"title"`
	Subreddit string `json:"subreddit"`
	NumComments int  `json:"num_comments"`
}

type Preview struct {
	Enabled bool `json:"enabled"`
	Images  []struct {
		Source struct {
			URL    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"source"`
		Resolutions []struct {
			URL    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"resolutions"`

	} `json:"images"`
}
