package request

type CustomerQueryRequest struct {
	Page    int    `form:"page,default=1"`
	Size    int    `form:"size,default=10"`
	Search  string `form:"search"`                  // Name or Email
	Mobile  string `form:"filter_mobile"`           // Renamed to filter_mobile
	Sort    string `form:"sort,default=created_at"` // Column to sort by
	SortDir string `form:"sort_dir,default=desc"`   // asc or desc
}
