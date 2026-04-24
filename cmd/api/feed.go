package main

import (
	"net/http"

	"github.com/Davidmuthee12/socials/internal/store"
)

// GetUserFeed godoc
//
//	@Summary		Get user feed
//	@Description	Returns a paginated feed of posts for a user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int		false	"Items per page (1-20)"
//	@Param			offset	query		int		false	"Offset"
//	@Param			sort	query		string	false	"Sort order"	Enums(asc, desc)
//	@Param			tags	query		string	false	"Comma-separated tags"
//	@Param			search	query		string	false	"Search term"
//	@Param			since	query		string	false	"Lower bound datetime (YYYY-MM-DD HH:MM:SS)"
//	@Param			until	query		string	false	"Upper bound datetime (YYYY-MM-DD HH:MM:SS)"
//	@Success		200		{array}		store.PostWithMetadata
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(fq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	feed, err := app.store.Posts.GetUserFeed(ctx, int64(26), fq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}
}
