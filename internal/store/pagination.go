package store

import (
	"net/http"
	"strconv"
	"strings"
	"time"
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

func (fq PaginatedQuery) Parse(r *http.Request) (PaginatedQuery, error) {
	qs := r.URL.Query()

	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, nil
		}
		fq.Limit = l
	}

	offset := qs.Get("offset")
	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return fq, nil
		}
		fq.Offset = o
	}

	sort := qs.Get("sort")
	if sort != "" {
		fq.Sort = sort
	}

	tags := qs.Get("tags")
	if tags != "" {
		fq.Tags = strings.Split(tags, ",")
	}

	search := qs.Get("search")
	if search != "" {
		fq.Search = search
	}

	since := qs.Get("since")
	if since != "" {
		fq.Since = parseTime(since)
	} else {
		// Init time
		fq.Since = InitSinceTime
	}

	until := qs.Get("until")
	if until != "" {
		fq.Until = parseTime(until)
	} else {
		fq.Until = time.Now().Format(time.DateTime)
	}

	return fq, nil
}

func parseTime(s string) string {
	t, err := time.Parse(time.DateTime, s)
	if err != nil {
		return ""
	}

	return t.Format(time.DateTime)
}
