package wallpaper

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

func SearchWallpapers(c *gin.Context) {
	var req WallpaperRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	}

	posts, err := req.FetchWallpaper()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, posts)
}