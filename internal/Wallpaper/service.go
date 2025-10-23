package wallpaper

import (
	"strings"

	reddit "github.com/dis70rt/subpaper-backend/internal/Reddit"
)

func isValid(post *reddit.RedditPost) bool {
	return post.URL != "" && strings.Contains(post.URL, "i.redd.it") && !post.Is18
}

func (req *WallpaperRequest) FetchWallpaper() ([]WallpaperResponse, error) {
	client := reddit.NewClient()

	// subreddits := map[string][]string{
	// 	"anime":  {"Animewallpaper", "awwnime", "AnimeWallpapersSFW"},
	// 	"mobile": {"Amoledbackgrounds", "iphonewallpapers", "iWallpaper"},
	// }

	path := req.buildSearchPath("Animewallpaper")
	resp, err := client.GetReddit(path)
	if err != nil {
		return nil, err
	}

	var wallpapers []WallpaperResponse
	for _, child := range resp.Data.Children {
		post := child.Data

		if isValid(&post) {
			width, height := 0, 0
			if post.Preview.Enabled && len(post.Preview.Images) > 0 {
				images := post.Preview.Images
				last := images[len(images)-1]
				width = last.Source.Width
				height = last.Source.Height
			}

			if width >= 1080 && height >= 1920 {
				wallpapers = append(wallpapers, WallpaperResponse{
					ID:      post.ID,
					Post:    post.Post,
					Preview: "",
					URL:     post.URL,
					Score:   post.Score,
					Height:  height,
					Width:   width,
				})
			}
		}
	}

	return wallpapers, nil
}
