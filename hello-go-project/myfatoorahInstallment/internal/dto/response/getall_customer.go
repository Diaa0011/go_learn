package response

type PaginationMetadata struct {
	TotalRecords int64 `json:"total_records" example:"100"`
	TotalPages   int   `json:"total_pages" example:"10"`
	CurrentPage  int   `json:"current_page" example:"1"`
	Size         int   `json:"index" example:"10"`
}

type PaginatedResponse[T any] struct {
	Data []T                `json:"data"`
	Meta PaginationMetadata `json:"meta"`
}
