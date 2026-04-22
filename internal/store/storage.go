package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("record not found")
	QueryTimeoutDuration = time.Second * 5
)
type Storage struct {
	Posts interface {
		GetByID(context.Context, int64) (*Post, error)
		Create(context.Context, *Post) error
		DeletePost(context.Context, int64) error
		UpdatePost(context.Context, *Post) error
	}

	Users interface {
		GetByID(context.Context, int64) (*User, error)
		Create(context.Context, *User) error
	}

	Comments interface {
		Create(context.Context, *Comment) error
		GetPostById( context.Context, int64) ([]Comment, error)
	}
}


func NewStorage(db *sql.DB) Storage{
	return Storage {
		Posts: &PostStore{db},
		Users: &UserStore{db},
		Comments: &CommentStore{db},
	}
}