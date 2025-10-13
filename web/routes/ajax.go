package routes

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/internal/database"
	"github.com/Lekuruu/go-puush/internal/services"
)

type AjaxError struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

var ErrorUsernameAlreadySet = AjaxError{Error: true, Message: "You have already set a username."}
var ErrorPasswordIncorrect = AjaxError{Error: true, Message: "Current password incorrect."}
var ErrorUsernameTaken = AjaxError{Error: true, Message: "That username is already taken."}
var ErrorServerError = AjaxError{Error: true, Message: "An internal server error occurred."}
var ErrorBadRequest = AjaxError{Error: true, Message: "Bad request."}
var NoError = AjaxError{Error: false, Message: ""}

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

		// Remove assets from storage, do nothing on error
		ctx.State.Storage.RemoveThumbnail(upload.Key())
		ctx.State.Storage.RemoveUpload(upload.Key())
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

func ChangePassword(ctx *app.Context) {
	user, err := GetUserSession(ctx)
	if err != nil || user == nil {
		http.Redirect(ctx.Response, ctx.Request, "/login", http.StatusSeeOther)
		return
	}

	err = ctx.Request.ParseForm()
	if err != nil {
		renderJson(200, ErrorBadRequest, ctx)
		return
	}

	currentPassword := ctx.Request.FormValue("c")
	newPassword := ctx.Request.FormValue("p")

	if currentPassword == "" || newPassword == "" {
		renderJson(200, ErrorBadRequest, ctx)
		return
	}

	if !app.VerifyPasswordHash(currentPassword, user.Password) {
		renderJson(200, ErrorPasswordIncorrect, ctx)
		return
	}

	newPasswordHash, err := app.CreatePasswordHash(newPassword)
	if err != nil {
		renderJson(200, ErrorServerError, ctx)
		return
	}

	err = services.UpdateUserPassword(user.Id, newPasswordHash, ctx.State)
	if err != nil {
		renderJson(200, ErrorServerError, ctx)
		return
	}

	renderRaw(200, "text/html", []byte{}, ctx)
}

func CheckUsername(ctx *app.Context) {
	err := ctx.Request.ParseForm()
	if err != nil {
		renderJson(200, ErrorBadRequest, ctx)
		return
	}

	username := ctx.Request.FormValue("u")
	if username == "" {
		renderJson(200, ErrorBadRequest, ctx)
		return
	}

	existingUser, _ := services.FetchUserByName(username, ctx.State)
	if existingUser != nil {
		renderJson(200, ErrorUsernameTaken, ctx)
		return
	}

	renderJson(200, NoError, ctx)
}

func ClaimUsername(ctx *app.Context) {
	user, err := GetUserSession(ctx)
	if err != nil || user == nil {
		http.Redirect(ctx.Response, ctx.Request, "/login", http.StatusSeeOther)
		return
	}

	if user.Name != "" {
		renderJson(200, ErrorUsernameAlreadySet, ctx)
		return
	}

	err = ctx.Request.ParseForm()
	if err != nil {
		renderJson(200, ErrorBadRequest, ctx)
		return
	}

	username := ctx.Request.FormValue("u")
	if username == "" {
		renderJson(200, ErrorBadRequest, ctx)
		return
	}

	existingUser, _ := services.FetchUserByName(username, ctx.State)
	if existingUser != nil {
		renderJson(200, ErrorUsernameTaken, ctx)
		return
	}

	// TODO: Validate username (allowed characters, length, etc.)
	user.Name = username
	user.UsernameSetupReminder = false
	err = services.UpdateUser(user, ctx.State)
	if err != nil {
		renderJson(200, ErrorServerError, ctx)
		return
	}

	renderJson(200, NoError, ctx)
}

func StopAskingAboutUsername(ctx *app.Context) {
	user, err := GetUserSession(ctx)
	if err != nil || user == nil {
		http.Redirect(ctx.Response, ctx.Request, "/login", http.StatusSeeOther)
		return
	}

	user.UsernameSetupReminder = false
	err = services.UpdateUser(user, ctx.State)
	if err != nil {
		renderJson(200, ErrorServerError, ctx)
		return
	}

	renderJson(200, NoError, ctx)
}

func resolveTargetUploadsFromQuery(ctx *app.Context) ([]*database.Upload, error) {
	identifiersString := ctx.Request.URL.Query().Get("i")
	if identifiersString == "" {
		return []*database.Upload{}, nil
	}

	identifiers := strings.Split(identifiersString, ",")

	uploads, err := services.FetchManyUploadsByIdentifiers(identifiers, ctx.State)
	if err != nil {
		return nil, err
	}

	return uploads, nil
}

func resolveTargetUploadsFromForm(ctx *app.Context) ([]*database.Upload, error) {
	identifiers := ctx.Request.Form["i[]"]
	if len(identifiers) == 0 {
		return []*database.Upload{}, nil
	}

	uploads, err := services.FetchManyUploadsByIdentifiers(identifiers, ctx.State, "Pool")
	if err != nil {
		return nil, err
	}

	return uploads, nil
}
