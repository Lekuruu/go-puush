package routes

import (
	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/internal/database"
	"github.com/Lekuruu/go-puush/internal/services"
)

const GalleryPoolName = "Gallery"

type GalleryFeedResponse struct {
	Count   int           `json:"count"`
	Objects []GalleryItem `json:"objects"`
}

type GalleryItem struct {
	Index      int    `json:"index"`
	Identifier string `json:"id"`
	Size       string `json:"size"`
	Filename   string `json:"name"`
	MimeType   string `json:"mime_type"`
}

func NewGalleryFeed(pool *database.Pool) *GalleryFeedResponse {
	feed := &GalleryFeedResponse{
		Count:   0,
		Objects: []GalleryItem{},
	}

	for i, upload := range pool.Uploads {
		feed.Objects = append(feed.Objects, GalleryItem{
			Index:      i + 1,
			Identifier: upload.Identifier,
			Size:       upload.SizeHumanReadable(),
			Filename:   upload.Filename,
			MimeType:   upload.MimeType,
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

	pool, err := services.FetchPoolByUserAndName(user.Id, GalleryPoolName, ctx.State, "Uploads")
	if err != nil || pool == nil {
		renderText(404, "page not found", ctx)
		return
	}

	renderJson(200, NewGalleryFeed(pool), ctx)
}
