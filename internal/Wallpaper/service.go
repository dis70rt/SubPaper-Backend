package wallpaper

import (
	"fmt"
	"html"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	reddit "github.com/dis70rt/subpaper-backend/internal/Reddit"
	"github.com/patrickmn/go-cache"
)

type Service struct {
    client 	*reddit.RedditClient
	cache 	*cache.Cache
}

func NewService(client *reddit.RedditClient, cache *cache.Cache) *Service {
    return &Service{
        client: client,
		cache: cache,
    }
}

func (req *WallpaperRequest) buildSearchPath(subreddit string) string {
	req.setDefaults()
	params := url.Values{}
	params.Set("q", req.Query)
	params.Set("sort", req.Sort)
	params.Set("t", req.TimeFilter)
	params.Set("limit", strconv.Itoa(req.Limit))
	params.Set("restrict_sr", "1")

	return fmt.Sprintf("/r/%s/search?%s", subreddit, params.Encode())
}

func (service *Service) FetchWallpaper(req *WallpaperRequest) ([]WallpaperResponse, error) {
	
	cacheKey := fmt.Sprintf("wallpaper:%s:%s:%s:%s:%d", 
        req.Type, req.Query, req.Sort, req.TimeFilter, req.Limit)
	
	if cached, found := service.cache.Get(cacheKey); found {
        return cached.([]WallpaperResponse), nil
    }

	subreddits := map[string][]string{
		"anime":  {"Animewallpaper", "awwnime", "AnimeWallpapersSFW"},
		"mobile": {"Amoledbackgrounds", "iphonewallpapers", "iWallpaper"},
	}

	var targetSubs []string
	switch req.Type {
	case "anime":
		targetSubs = subreddits["anime"]
	case "mobile":
		targetSubs = subreddits["mobile"]
	}

	type result struct {
		wallpapers []WallpaperResponse
		err        error
	}

	resultsCh := make(chan result, len(targetSubs))
	defer close(resultsCh)

	for _, sub := range targetSubs {
		subreddit := sub
		go func() {
			path := req.buildSearchPath(subreddit)
			resp, err := service.client.GetReddit(path)
			if err != nil {
				resultsCh <- result{nil, err}
				return
			}

			wallpapers := extractPosts(resp)
			resultsCh <- result{wallpapers, nil}
		}()
	}

	var allWallpapers []WallpaperResponse
	for i := 0; i < len(targetSubs); i++ {
		res := <-resultsCh
		if res.err != nil {
			continue
		}
		allWallpapers = append(allWallpapers, res.wallpapers...)
	}

	sort.Slice(allWallpapers, func(i, j int) bool {
		return allWallpapers[i].Score > allWallpapers[j].Score
	})

	service.cache.Set(cacheKey, allWallpapers, 12*time.Hour)
	return allWallpapers, nil
}

func isValid(post *reddit.RedditPost) bool {
	return post.URL != "" && strings.Contains(post.URL, "i.redd.it") && !post.Is18
}

func extractPosts(resp *reddit.RedditAPIResponse) []WallpaperResponse {
    var wallpapers []WallpaperResponse
	for _, child := range resp.Data.Children {
		post := child.Data
		if isValid(&post) {
			width, height := 0, 0
			preview := ""
			if post.Preview.Enabled && len(post.Preview.Images) > 0 {
				images := post.Preview.Images
				last := images[len(images)-1]
				width = last.Source.Width
				height = last.Source.Height
				preview = html.UnescapeString(post.Preview.Images[0].Resolutions[2].URL)
			}

			if width >= 1080 && height >= 1920 {
				wallpapers = append(wallpapers, WallpaperResponse{
					ID:      post.ID,
					Post:    post.Post,
					Preview: preview,
					URL:     post.URL,
					Score:   post.Score,
					Height:  height,
					Width:   width,
				})
			}
		}
	}
    return wallpapers
}
