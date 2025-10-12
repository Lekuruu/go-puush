package routes

import (
	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/internal/database"
	"github.com/Lekuruu/go-puush/internal/services"
)

const GalleryPoolName = "Public"

type GalleryFeedResponse struct {
	Count   int           `json:"count"`
	Objects []GalleryItem `json:"objects"`
}

type GalleryItem struct {
	Index      int    `json:"index"`
	Identifier string `json:"id"`
	Size       string `json:"size"`
	Filename   string `json:"name"`
}

func NewGalleryFeed(pool *database.Pool) *GalleryFeedResponse {
	feed := &GalleryFeedResponse{
		Count:   0,
		Objects: []GalleryItem{},
	}

	for i, upload := range pool.Uploads {
		if upload.Link == nil {
			// NOTE: This should not happen, but just in case if the link
			// 		 is missing we skip this upload to avoid a crash.
			continue
		}

		feed.Objects = append(feed.Objects, GalleryItem{
			Index:      i + 1,
			Identifier: upload.Link.Identifier,
			Size:       upload.SizeHumanReadable(),
			Filename:   upload.Filename,
		})
	}

	feed.Count = len(feed.Objects)
	return feed
}

func Gallery(ctx *app.Context) {
	username := ctx.Vars["username"]
	if username == "" {
		renderText(404, "page not found", ctx)
		return
	}

	user, err := services.FetchUserByName(username, ctx.State)
	if err != nil || user == nil {
		renderText(404, "page not found", ctx)
		return
	}

	pool, err := services.FetchPoolByUserAndName(user.Id, GalleryPoolName, ctx.State)
	if err != nil || pool == nil {
		renderText(404, "page not found", ctx)
		return
	}

	renderTemplate(ctx, "gallery", map[string]interface{}{
		"Title": "gallery",
		"User":  user,
		"Pool":  pool,
	})
}

func GalleryFeed(ctx *app.Context) {
	username := ctx.Vars["username"]
	if username == "" {
		renderText(404, "page not found", ctx)
		return
	}

	user, err := services.FetchUserByName(username, ctx.State)
	if err != nil || user == nil {
		renderText(404, "page not found", ctx)
		return
	}

	pool, err := services.FetchPoolByUserAndName(user.Id, GalleryPoolName, ctx.State, "Uploads", "Uploads.Link")
	if err != nil || pool == nil {
		renderText(404, "page not found", ctx)
		return
	}

	renderJson(200, NewGalleryFeed(pool), ctx)
}
