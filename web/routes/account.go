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
	// TODO: Implement authentication check
	// http.Redirect(ctx.Response, ctx.Request, "/login", http.StatusTemporaryRedirect)

	// For testing purposes, render the account page with user ID 1
	user, err := services.FetchUserById(1, ctx.State, "Pools")
	if err != nil || user == nil {
		renderText(404, "page not found", ctx)
		return
	}

	selectedPool, err := resolvePoolFromRequest(user, ctx)
	if err != nil {
		http.Redirect(ctx.Response, ctx.Request, "/account", http.StatusTemporaryRedirect)
		return
	}

	renderTemplate(ctx, "account/home", map[string]interface{}{
		"Title":        "account",
		"User":         user,
		"SelectedPool": selectedPool,
		"ViewType":     resolveViewTypeFromRequest(ctx),
	})
}

func resolveViewTypeFromRequest(ctx *app.Context) string {
	query := ctx.Request.URL.Query()

	if query.Has("list") {
		return "list"
	}
	if query.Has("grid") {
		return "grid"
	}

	// Default to grid view
	return "grid"
}

func resolvePoolFromRequest(user *database.User, ctx *app.Context) (*database.Pool, error) {
	poolId := ctx.Request.URL.Query().Get("p")
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
