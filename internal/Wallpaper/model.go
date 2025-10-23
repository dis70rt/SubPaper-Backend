package wallpaper

type WallpaperResponse struct {
	ID 		string	`json:"id"`
	Post 	string	`json:"post"`
	Preview string	`json:"preview"`
	URL 	string	`json:"url"`
	Score 	int		`json:"score"`
	Height 	int		`json:"height"`
	Width 	int		`json:"width"`
}

type WallpaperRequest struct {
	Type 		string	`form:"type"`
	Query 		string	`form:"q"`
	Sort 		string 	`form:"sort" default:"top"`
	Limit 		int		`form:"limit" default:"50"`
	TimeFilter 	string	`form:"t" default:"all"`
}

func (req *WallpaperRequest) setDefaults() {
	if req.Sort == "" {
		req.Sort = "top"
	}
	if req.Limit == 0 {
		req.Limit = 50
	}
	if req.TimeFilter == "" {
		req.TimeFilter = "all"
	}
}

