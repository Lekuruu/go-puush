package routes

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/internal/database"
	"github.com/Lekuruu/go-puush/internal/services"
)

func MoveDialog(ctx *app.Context) {
	user, err := GetUserSession(ctx, "Pools")
	if err != nil || user == nil {
		http.Redirect(ctx.Response, ctx.Request, "/login", http.StatusSeeOther)
		return
	}

	poolThumbnails, err := resolvePoolThumbnails(user.Pools, ctx)
	if err != nil {
		renderText(500, "Server error", ctx)
		return
	}

	targetUploads, err := resolveTargetUploadsFromQuery(ctx)
	if err != nil {
		renderText(500, "Server error", ctx)
		return
	}

	renderTemplate(ctx, "ajax/move", map[string]interface{}{
		"PoolThumbnails": poolThumbnails,
		"TargetUploads":  targetUploads,
		"User":           user,
	})
}

func MoveUpload(ctx *app.Context) {
	user, err := GetUserSession(ctx)
	if err != nil || user == nil {
		http.Redirect(ctx.Response, ctx.Request, "/login", http.StatusSeeOther)
		return
	}

	err = ctx.Request.ParseForm()
	if err != nil {
		renderText(400, "Bad request", ctx)
		return
	}

	targetPoolIdStr := ctx.Request.FormValue("p")
	if targetPoolIdStr == "" {
		renderText(400, "Bad request", ctx)
		return
	}

	targetPoolId, err := strconv.Atoi(targetPoolIdStr)
	if err != nil {
		renderText(400, "Bad request", ctx)
		return
	}

	targetPool, err := services.FetchPoolById(targetPoolId, ctx.State)
	if err != nil || targetPool == nil {
		renderText(400, "Pool was not found", ctx)
		return
	}
	if targetPool.UserId != user.Id {
		renderText(403, "You do not own that pool", ctx)
		return
	}

	targetUploads, err := resolveTargetUploadsFromForm(ctx)
	if err != nil {
		renderText(500, "Server error", ctx)
		return
	}
	if len(targetUploads) == 0 {
		renderText(400, "No uploads were specified", ctx)
		return
	}

	for _, upload := range targetUploads {
		if upload.Pool.UserId != user.Id {
			renderText(403, "You do not own that upload", ctx)
			return
		}
		if upload.PoolId == targetPool.Id {
			continue
		}

		upload.PoolId = targetPool.Id
		err = services.UpdateUploadPool(upload.Id, targetPool.Id, ctx.State)
		if err != nil {
			renderText(500, "Server error", ctx)
			return
		}
	}

	err = services.UpdatePoolUploadCounts(user, ctx.State)
	if err != nil {
		renderText(500, "Server error", ctx)
		return
	}

	renderRaw(200, "text/html", []byte{}, ctx)
}

func DeleteUpload(ctx *app.Context) {
	user, err := GetUserSession(ctx)
	if err != nil || user == nil {
		http.Redirect(ctx.Response, ctx.Request, "/login", http.StatusSeeOther)
		return
	}

	err = ctx.Request.ParseForm()
	if err != nil {
		renderText(400, "Bad request", ctx)
		return
	}

	targetUploads, err := resolveTargetUploadsFromForm(ctx)
	if err != nil {
		renderText(500, "Server error", ctx)
		return
	}
	if len(targetUploads) == 0 {
		renderText(400, "No uploads were specified", ctx)
		return
	}

	for _, upload := range targetUploads {
		if upload.Pool.UserId != user.Id {
			renderText(403, "You do not own that upload", ctx)
			return
		}

		err = services.DeleteUpload(upload, ctx.State)
		if err != nil {
			renderText(500, "Server error", ctx)
			return
		}
	}

	err = services.UpdatePoolUploadCounts(user, ctx.State)
	if err != nil {
		renderText(500, "Server error", ctx)
		return
	}

	renderRaw(200, "text/html", []byte{}, ctx)
}

func UpdateDefaultPool(ctx *app.Context) {
	user, err := GetUserSession(ctx)
	if err != nil || user == nil {
		http.Redirect(ctx.Response, ctx.Request, "/login", http.StatusSeeOther)
		return
	}

	err = ctx.Request.ParseForm()
	if err != nil {
		renderText(400, "Bad request", ctx)
		return
	}

	targetPoolIdStr := ctx.Request.FormValue("p")
	if targetPoolIdStr == "" {
		renderText(400, "Bad request", ctx)
		return
	}

	targetPoolId, err := strconv.Atoi(targetPoolIdStr)
	if err != nil {
		renderText(400, "Bad request", ctx)
		return
	}

	targetPool, err := services.FetchPoolById(targetPoolId, ctx.State)
	if err != nil || targetPool == nil {
		renderText(400, "Pool was not found", ctx)
		return
	}
	if targetPool.UserId != user.Id {
		renderText(403, "You do not own that pool", ctx)
		return
	}

	user.DefaultPoolId = targetPool.Id
	err = services.UpdateUserDefaultPool(user.Id, targetPool.Id, ctx.State)
	if err != nil {
		renderText(500, "Server error", ctx)
		return
	}

	renderTemplate(ctx, "ajax/success", nil)
}

func resolveTargetUploadsFromQuery(ctx *app.Context) ([]*database.Upload, error) {
	identifiersString := ctx.Request.URL.Query().Get("i")
	if identifiersString == "" {
		return []*database.Upload{}, nil
	}

	identifiers := strings.Split(identifiersString, ",")

	links, err := services.FetchManyShortLinksByIdentifiers(identifiers, ctx.State, "Upload")
	if err != nil {
		return nil, err
	}

	var uploads []*database.Upload
	for _, link := range links {
		uploads = append(uploads, link.Upload)
	}
	return uploads, nil
}

func resolveTargetUploadsFromForm(ctx *app.Context) ([]*database.Upload, error) {
	identifiers := ctx.Request.Form["i[]"]
	if len(identifiers) == 0 {
		return []*database.Upload{}, nil
	}

	links, err := services.FetchManyShortLinksByIdentifiers(identifiers, ctx.State, "Upload", "Upload.Pool")
	if err != nil {
		return nil, err
	}

	var uploads []*database.Upload
	for _, link := range links {
		uploads = append(uploads, link.Upload)
	}
	return uploads, nil
}
