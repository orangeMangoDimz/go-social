package protocol

import (
	"net/http"

	postsEntity "github.com/orangeMangoDimz/go-social/internal/entities/posts"
	usersEntity "github.com/orangeMangoDimz/go-social/internal/entities/users"
)

type userKey string
type postKey string

const UserCtx userKey = "user"
const PostCtx postKey = "post"

func GetUserFromContext(r *http.Request) *usersEntity.User {
	user, ok := r.Context().Value(UserCtx).(*usersEntity.User)
	if !ok {
		return nil
	}
	return user
}

func GetPostFromContext(r *http.Request) *postsEntity.Post {
	post, ok := r.Context().Value(PostCtx).(*postsEntity.Post)
	if !ok {
		return nil
	}
	return post
}
