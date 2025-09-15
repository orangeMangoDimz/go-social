package postsHandler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	payloadEntity "github.com/orangeMangoDimz/go-social/internal/entities/payload"
	postsEntity "github.com/orangeMangoDimz/go-social/internal/entities/posts"
	"github.com/orangeMangoDimz/go-social/internal/server/http/protocol"
	"github.com/orangeMangoDimz/go-social/internal/service"
	"github.com/orangeMangoDimz/go-social/internal/storage"
	"github.com/orangeMangoDimz/go-social/internal/storage/postgres/pagination"
	"go.uber.org/zap"
)

type httpHandler struct {
	PostService    service.PostsService
	CommentService service.CommentService
	logger         zap.SugaredLogger
}

func newHTTPHandler(postService service.PostsService, commentService service.CommentService, logger zap.SugaredLogger) *httpHandler {
	return &httpHandler{
		PostService:    postService,
		CommentService: commentService,
		logger:         logger,
	}
}

// createPostHandler godoc
//
//	@Summary		Create a new post
//	@Description	Create a new post with title, content and optional tags. Requires JWT authentication.
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			payload	body		github_com_orangeMangoDimz_go-social_internal_entities_payload.CreatePOstPayload	true	"Post creation data"
//	@Success		200		{object}	github_com_orangeMangoDimz_go-social_internal_entities_posts.Post					"Created post"
//	@Failure		400		{object}	map[string]string																	"Bad request"
//	@Failure		401		{object}	map[string]string																	"Unauthorized - invalid or missing token"
//	@Failure		500		{object}	map[string]string																	"Internal server error"
//	@Router			/posts [post]
func (h *httpHandler) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload payloadEntity.CreatePOstPayload

	if err := protocol.ReadJSON(w, r, &payload); err != nil {
		protocol.BadRequestResponse(w, r, err)
		return
	}

	var validate = validator.New()
	if err := validate.Struct(payload); err != nil {
		protocol.BadRequestResponse(w, r, err)
		return
	}

	user := protocol.GetUserFromContext(r)

	post := postsEntity.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserId:  user.ID,
	}

	ctx := r.Context()
	if err := h.PostService.Create(ctx, &post); err != nil {
		protocol.InternalServerError(w, r, err)
		return
	}

	if err := protocol.JsonResponse(w, http.StatusOK, &post); err != nil {
		protocol.InternalServerError(w, r, err)
		return
	}
}

// getPostHandler godoc
//
//	@Summary		Get post by ID
//	@Description	Get detailed information about a specific post including comments
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			postID	path		int																	true	"Post ID"	example(1)
//	@Success		200		{object}	github_com_orangeMangoDimz_go-social_internal_entities_posts.Post	"Post information with comments"
//	@Failure		404		{object}	map[string]string													"Post not found"
//	@Failure		500		{object}	map[string]string													"Internal server error"
//	@Router			/posts/{postID} [get]
func (h *httpHandler) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := protocol.GetPostFromContext(r)
	if post == nil {
		protocol.NotFoundResponse(w, r, errors.New("post not found"))
		return
	}

	ctx := r.Context()
	comment, err := h.CommentService.GetByPostID(ctx, post.ID)
	if err != nil {
		protocol.InternalServerError(w, r, err)
		return
	}

	post.Comments = comment
	if err := protocol.JsonResponse(w, http.StatusOK, post); err != nil {
		protocol.InternalServerError(w, r, err)
		return
	}
}

// getUserPostFeed godoc
//
//	@Summary		Get user's post feed
//	@Description	Get a paginated feed of posts from followed users and own posts. Requires JWT authentication.
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			limit	query		int																	false	"Number of posts per page (1-20)"	default(20)		example(10)
//	@Param			offset	query		int																	false	"Number of posts to skip"			default(0)		example(0)
//	@Param			sort	query		string																false	"Sort order (asc/desc)"				default(desc)	Enums(asc, desc)
//	@Param			search	query		string																false	"Search in title and content"		example("golang")
//	@Param			tags	query		string																false	"Comma-separated list of tags"		example("golang,programming")
//	@Param			since	query		string																false	"Posts created after this date"		example("2024-01-01 00:00:00")
//	@Param			until	query		string																false	"Posts created before this date"	example("2024-12-31 23:59:59")
//	@Success		200		{array}		github_com_orangeMangoDimz_go-social_internal_entities_posts.Feed	"User's post feed"
//	@Failure		400		{object}	map[string]string													"Bad request"
//	@Failure		401		{object}	map[string]string													"Unauthorized - invalid or missing token"
//	@Failure		500		{object}	map[string]string													"Internal server error"
//	@Router			/posts/feed [get]
func (h *httpHandler) getUserPostFeed(w http.ResponseWriter, r *http.Request) {
	fmt.Print("test")
	fq := pagination.PaginatedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	fq, err := fq.Parse(r)
	fmt.Print("err: ", err)
	if err != nil {
		protocol.BadRequestResponse(w, r, err)
		return
	}

	var validate = validator.New()
	if err := validate.Struct(fq); err != nil {
		protocol.BadRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	user := protocol.GetUserFromContext(r)
	fmt.Print("USER ID: ", user.ID)

	userID := user.ID
	feed, err := h.PostService.GetUserFeed(ctx, int64(userID), fq)
	if err != nil {
		protocol.InternalServerError(w, r, err)
		return
	}

	if err := protocol.JsonResponse(w, http.StatusOK, feed); err != nil {
		protocol.InternalServerError(w, r, err)
		return
	}

}

// updatePostHandler godoc
//
//	@Summary		Update a post
//	@Description	Update the title and/or content of an existing post. Requires JWT authentication.
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			postID	path		int																					true	"Post ID"	example(1)
//	@Param			payload	body		github_com_orangeMangoDimz_go-social_internal_entities_payload.UpdatePostPayload	true	"Post update data"
//	@Success		200		{object}	github_com_orangeMangoDimz_go-social_internal_entities_posts.Post					"Updated post"
//	@Failure		400		{object}	map[string]string																	"Bad request"
//	@Failure		401		{object}	map[string]string																	"Unauthorized - invalid or missing token"
//	@Failure		404		{object}	map[string]string																	"Post not found"
//	@Failure		500		{object}	map[string]string																	"Internal server error"
//	@Router			/posts/{postID} [patch]
func (h *httpHandler) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := protocol.GetPostFromContext(r)
	if post == nil {
		protocol.NotFoundResponse(w, r, storage.ErrNotFound)
		return
	}

	var payload payloadEntity.UpdatePostPayload
	if err := protocol.ReadJSON(w, r, &payload); err != nil {
		protocol.BadRequestResponse(w, r, err)
		return
	}

	var validate = validator.New()
	if err := validate.Struct(payload); err != nil {
		protocol.BadRequestResponse(w, r, err)
		return
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}
	if payload.Content != nil {
		post.Content = *payload.Content
	}

	ctx := r.Context()
	err := h.PostService.Update(ctx, post)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrNotFound):
			protocol.NotFoundResponse(w, r, err)
		default:
			protocol.InternalServerError(w, r, err)
		}
		return
	}

	if err := protocol.JsonResponse(w, http.StatusOK, post); err != nil {
		protocol.InternalServerError(w, r, err)
		return
	}

}

// deletePostHandler godoc
//
//	@Summary		Delete a post
//	@Description	Delete a specific post by its ID. Requires JWT authentication.
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			postID	path	int	true	"Post ID"	example(1)
//	@Success		204		"Post successfully deleted"
//	@Failure		401		{object}	map[string]string	"Unauthorized - invalid or missing token"
//	@Failure		404		{object}	map[string]string	"Post not found"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/posts/{postID} [delete]
func (h *httpHandler) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	post := protocol.GetPostFromContext(r)
	if post == nil {
		protocol.NotFoundResponse(w, r, storage.ErrNotFound)
		return
	}

	ctx := r.Context()
	err := h.PostService.Delete(ctx, post.ID)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrNotFound):
			protocol.NotFoundResponse(w, r, err)
		default:
			protocol.InternalServerError(w, r, err)
		}
		return
	}

	if err := protocol.JsonResponse(w, http.StatusAccepted, nil); err != nil {
		protocol.InternalServerError(w, r, err)
		return
	}
}

func (h *httpHandler) postContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParm := chi.URLParam(r, "postID")
		id, err := strconv.ParseInt(idParm, 10, 64)
		if err != nil {
			protocol.InternalServerError(w, r, err)
			return
		}
		ctx := r.Context()

		post, err := h.PostService.GetById(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, storage.ErrNotFound):
				protocol.NotFoundResponse(w, r, err)
			default:
				protocol.InternalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, protocol.PostCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
