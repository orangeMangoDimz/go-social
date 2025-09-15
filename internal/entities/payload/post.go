package payloadEntity

// CreatePOstPayload represents the request payload for creating a new post
//
//	@Description	Request payload for creating a new post
type CreatePOstPayload struct {
	Title   string   `json:"title" validate:"required,max=100" example:"My First Post"`                     // Post title (max 100 characters)
	Content string   `json:"content" validate:"required,max=1000" example:"This is the content of my post"` // Post content (max 1000 characters)
	Tags    []string `json:"tags" example:"golang,programming"`                                             // Post tags
}

// UpdatePostPayload represents the request payload for updating a post
//
//	@Description	Request payload for updating an existing post
type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100" example:"Updated Post Title"`      // Updated post title (max 100 characters)
	Content *string `json:"content" validate:"omitempty,max=1000" example:"Updated post content"` // Updated post content (max 1000 characters)
}
