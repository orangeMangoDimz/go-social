package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/orangeMangoDimz/go-social/internal/store"
)

type postKey string

const postCtx postKey = "post"

// CreatePOstPayload represents the request payload for creating a new post
//
//	@Description	Request payload for creating a new post
type CreatePOstPayload struct {
	Title   string   `json:"title" validate:"required,max=100" example:"My First Post"`                     // Post title (max 100 characters)
	Content string   `json:"content" validate:"required,max=1000" example:"This is the content of my post"` // Post content (max 1000 characters)
	Tags    []string `json:"tags" example:"golang,programming"`                                             // Post tags
}

// createPostHandler creates a new post
//
//	@Summary		Create a new post
//	@Description	Create a new post with title, content and optional tags. Requires JWT authentication.
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreatePOstPayload	true	"Post creation data"
//	@Success		200		{object}	store.Post			"Created post"
//	@Failure		400		{object}	map[string]string	"Bad request"
//	@Failure		401		{object}	map[string]string	"Unauthorized - invalid or missing token"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Security		BearerAuth
//	@Router			/posts [post]
func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePOstPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := getUserFromContext(r)

	post := store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserId:  user.ID,
	}

	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, &post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, &post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// getPostHandler retrieves a post by ID with its comments
//
//	@Summary		Get post by ID
//	@Description	Get detailed information about a specific post including comments
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int					true	"Post ID"	example(1)
//	@Success		200		{object}	store.Post			"Post information with comments"
//	@Failure		404		{object}	map[string]string	"Post not found"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/posts/{postID} [get]
//
//	@Security		BearerAuth
func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromContext(r)
	if post == nil {
		app.notFoundResponse(w, r, errors.New("post not found"))
		return
	}

	comment, err := app.store.Comments.GetByPostID(r.Context(), post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comments = comment
	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// UpdatePostPayload represents the request payload for updating a post
//
//	@Description	Request payload for updating an existing post
type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100" example:"Updated Post Title"`      // Updated post title (max 100 characters)
	Content *string `json:"content" validate:"omitempty,max=1000" example:"Updated post content"` // Updated post content (max 1000 characters)
}

// getUserPostFeed retrieves the user's personalized post feed
//
//	@Summary		Get user's post feed
//	@Description	Get a paginated feed of posts from followed users and own posts. Requires JWT authentication.
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int					false	"Number of posts per page (1-20)"	default(20)		example(10)
//	@Param			offset	query		int					false	"Number of posts to skip"			default(0)		example(0)
//	@Param			sort	query		string				false	"Sort order (asc/desc)"				default(desc)	Enums(asc, desc)
//	@Param			search	query		string				false	"Search in title and content"		example("golang")
//	@Param			tags	query		string				false	"Comma-separated list of tags"		example("golang,programming")
//	@Param			since	query		string				false	"Posts created after this date"		example("2024-01-01 00:00:00")
//	@Param			until	query		string				false	"Posts created before this date"	example("2024-12-31 23:59:59")
//	@Success		200		{array}		store.Feed			"User's post feed"
//	@Failure		400		{object}	map[string]string	"Bad request"
//	@Failure		401		{object}	map[string]string	"Unauthorized - invalid or missing token"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Security		BearerAuth
//	@Router			/users/feed [get]
func (app *application) getUserPostFeed(w http.ResponseWriter, r *http.Request) {
	fq := store.PaginatedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := validate.Struct(fq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	user := getUserFromContext(r)

	userID := user.ID

	feed, err := app.store.Posts.GetUserFeed(ctx, int64(userID), fq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

// updatePostHandler updates an existing post
//
//	@Summary		Update a post
//	@Description	Update the title and/or content of an existing post. Requires JWT authentication.
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int					true	"Post ID"	example(1)
//	@Param			payload	body		UpdatePostPayload	true	"Post update data"
//	@Success		200		{object}	store.Post			"Updated post"
//	@Failure		400		{object}	map[string]string	"Bad request"
//	@Failure		401		{object}	map[string]string	"Unauthorized - invalid or missing token"
//	@Failure		404		{object}	map[string]string	"Post not found"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Security		BearerAuth
//	@Router			/posts/{postID} [patch]
func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromContext(r)
	if post == nil {
		app.notFoundResponse(w, r, store.ErrNotFound)
		return
	}

	var payload UpdatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}
	if payload.Content != nil {
		post.Content = *payload.Content
	}

	if err := app.store.Posts.Update(r.Context(), post); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// deletePostHandler deletes a post by ID
//
//	@Summary		Delete a post
//	@Description	Delete a specific post by its ID. Requires JWT authentication.
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path	int	true	"Post ID"	example(1)
//	@Success		204		"Post successfully deleted"
//	@Failure		401		{object}	map[string]string	"Unauthorized - invalid or missing token"
//	@Failure		404		{object}	map[string]string	"Post not found"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Security		BearerAuth
//	@Router			/posts/{postID} [delete]
func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromContext(r)
	if post == nil {
		app.notFoundResponse(w, r, store.ErrNotFound)
		return
	}

	if err := app.store.Posts.Delete(r.Context(), post.ID); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusAccepted, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) postContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParm := chi.URLParam(r, "postID")
		id, err := strconv.ParseInt(idParm, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}
		ctx := r.Context()

		post, err := app.store.Posts.GetById(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromContext(r *http.Request) *store.Post {
	post, ok := r.Context().Value(postCtx).(*store.Post)
	if !ok {
		return nil
	}
	return post
}
