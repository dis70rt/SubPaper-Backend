package wallpaper
// TODO: Reuse the channel and token.
import (
	"sort"
	"strings"

	reddit "github.com/dis70rt/subpaper-backend/internal/Reddit"
)

func isValid(post *reddit.RedditPost) bool {
	return post.URL != "" && strings.Contains(post.URL, "i.redd.it") && !post.Is18
}

func (req *WallpaperRequest) FetchWallpaper() ([]WallpaperResponse, error) {
	client := reddit.NewClient()

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
			resp, err := client.GetReddit(path)
			if err != nil {
				resultsCh <- result{nil, err}
				return
			}

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
						preview = post.Preview.Images[0].Resolutions[2].URL
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

	return allWallpapers, nil
}
