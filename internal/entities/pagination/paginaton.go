package paginationEntity

var (
	InitSinceTime = "2025-01-01 00:00:00"
)

// PaginatedQuery represents query parameters for paginated requests
//
//	@Description	Query parameters for pagination and filtering
type PaginatedQuery struct {
	Limit  int      `json:"limit" validate:"gte=1,lte=20" example:"20"`         // Number of items per page (1-20)
	Offset int      `json:"offset" validate:"gte=0" example:"0"`                // Number of items to skip
	Sort   string   `json:"sort" validate:"oneof=asc desc" example:"desc"`      // Sort order (asc or desc)
	Tags   []string `json:"tags" validate:"max=5" example:"golang,programming"` // Filter by tags (max 5)
	Search string   `json:"search" validate:"max=100" example:"golang"`         // Search in title and content (max 100 chars)
	Since  string   `json:"since" example:"2024-01-01 00:00:00"`                // Filter posts created after this date
	Until  string   `json:"until" example:"2024-12-31 23:59:59"`                // Filter posts created before this date
}
