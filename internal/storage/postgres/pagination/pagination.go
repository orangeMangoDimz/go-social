package pagination

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	paginationEntity "github.com/orangeMangoDimz/go-social/internal/entities/pagination"
)

type PaginatedQuery paginationEntity.PaginatedQuery

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
		fq.Since = paginationEntity.InitSinceTime
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
