package routes

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/internal/database"
	"github.com/Lekuruu/go-puush/internal/services"
)

func Account(ctx *app.Context) {
	user, err := GetUserSession(ctx, "Pools")
	if err != nil || user == nil {
		http.Redirect(ctx.Response, ctx.Request, "/login", http.StatusSeeOther)
		return
	}

	selectedPool, err := resolvePoolFromRequest(user, ctx)
	if err != nil {
		http.Redirect(ctx.Response, ctx.Request, "/account", http.StatusTemporaryRedirect)
		return
	}

	poolThumbnails, err := resolvePoolThumbnails(user.Pools, ctx)
	if err != nil {
		http.Redirect(ctx.Response, ctx.Request, "/account", http.StatusTemporaryRedirect)
		return
	}

	// Update view type if changed in query
	user.ViewType = resolveViewTypeFromRequest(user, ctx)
	services.UpdateUser(user, ctx.State)

	renderTemplate(ctx, "account/home", map[string]interface{}{
		"Title":          "account",
		"User":           user,
		"SelectedPool":   selectedPool,
		"PoolThumbnails": poolThumbnails,
	})
}

func AccountSettings(ctx *app.Context) {
	user, err := GetUserSession(ctx, "Pools")
	if err != nil || user == nil {
		http.Redirect(ctx.Response, ctx.Request, "/login", http.StatusSeeOther)
		return
	}

	renderTemplate(ctx, "account/settings", map[string]interface{}{
		"Title": "account - settings",
		"User":  user,
	})
}

func AccountSubscription(ctx *app.Context) {
	user, err := GetUserSession(ctx)
	if err != nil || user == nil {
		http.Redirect(ctx.Response, ctx.Request, "/login", http.StatusSeeOther)
		return
	}

	renderTemplate(ctx, "account/subscription", map[string]interface{}{
		"Title": "account - subscription history",
		"User":  user,
	})
}

func AccountGoPro(ctx *app.Context) {
	user, err := GetUserSession(ctx)
	if err != nil || user == nil {
		http.Redirect(ctx.Response, ctx.Request, "/login", http.StatusSeeOther)
		return
	}

	renderTemplate(ctx, "account/pro", map[string]interface{}{
		"Title": "go pro",
		"User":  user,
	})
}

func resolveViewTypeFromRequest(user *database.User, ctx *app.Context) database.ViewType {
	query := ctx.Request.URL.Query()

	if query.Has("list") {
		return database.ViewTypeList
	}
	if query.Has("grid") {
		return database.ViewTypeGrid
	}

	// Default to grid view
	return user.ViewType
}

func resolvePoolFromRequest(user *database.User, ctx *app.Context) (*database.Pool, error) {
	poolId := ctx.Request.URL.Query().Get("pool")
	if poolId == "" {
		return services.FetchPoolById(user.DefaultPoolId, ctx.State, "Uploads", "Uploads.Link")
	}

	id, err := strconv.Atoi(poolId)
	if err != nil {
		return nil, err
	}

	pool, err := services.FetchPoolById(id, ctx.State, "Uploads", "Uploads.Link")
	if err != nil {
		return nil, err
	}
	if pool == nil {
		return nil, errors.New("pool not found")
	}
	if pool.UserId != user.Id {
		return nil, errors.New("pool does not belong to user")
	}

	return pool, nil
}

func resolvePoolThumbnails(pools []*database.Pool, ctx *app.Context) (map[int]string, error) {
	thumbnails := make(map[int]string, len(pools))
	for _, pool := range pools {
		lastUpload, err := services.FetchLastPoolUpload(pool.Id, ctx.State, "Link")
		if err != nil || lastUpload == nil {
			thumbnails[pool.Id] = pool.Identifier
			continue
		}
		thumbnails[pool.Id] = lastUpload.Link.Identifier
	}
	return thumbnails, nil
}
