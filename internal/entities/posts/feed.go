package postsEntity

// Feed represents a post in the user's feed with additional metadata
//
//	@Description	Post feed item with comment count
type Feed struct {
	Post
	TotalComments int64 `json:"total_comment" example:"5"` // Total number of comments on this post
}
