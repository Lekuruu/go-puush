package routes

import (
	"net/http"
	"time"

	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/internal/database"
	"github.com/Lekuruu/go-puush/internal/services"
)

const sessionCookieName = "PHPSESSID"
const sessionDuration = time.Hour * 24 * 30

func GetUserSession(ctx *app.Context, preload ...string) (*database.User, error) {
	// Try to get the session token from the cookies
	sessionToken, err := ctx.Request.Cookie(sessionCookieName)
	if err != nil {
		return nil, err
	}

	// Validate the session using the token
	session, err := services.ValidateSession(sessionToken.Value, ctx.State)
	if err != nil {
		return nil, err
	}

	// Fetch the user associated with the session
	return services.FetchUserById(session.UserId, ctx.State, preload...)
}

func SetUserSession(user *database.User, ctx *app.Context) error {
	// Create a new session for the user
	session, err := services.CreateSession(user.Id, sessionDuration, ctx.State)
	if err != nil {
		return err
	}

	// Set the session token as a cookie in the response
	http.SetCookie(ctx.Response, &http.Cookie{
		Name:    sessionCookieName,
		Value:   session.Token,
		Expires: session.ExpiresAt,
		Path:    "/",
	})
	return nil
}

func ClearUserSession(ctx *app.Context) error {
	// Try to get the session token from the cookies
	sessionToken, err := ctx.Request.Cookie(sessionCookieName)
	if err != nil {
		return err
	}

	// Invalidate the session using the token
	// We don't need to handle the error here
	services.DeleteSession(sessionToken.Value, ctx.State)

	// Remove the session cookie from the response
	http.SetCookie(ctx.Response, &http.Cookie{
		Name:    sessionCookieName,
		Value:   "",
		Expires: time.Unix(0, 0),
	})
	return nil
}

func UserPasswordAuthentication(email string, password string, state *app.State) (*database.User, bool) {
	user, err := services.FetchUserByNameOrEmail(email, state)
	if err != nil {
		return nil, false
	}

	if !app.VerifyPasswordHash(password, user.Password) {
		return nil, false
	}

	return user, true
}

func UserKeyAuthentication(key string, state *app.State) (*database.User, bool) {
	user, err := services.FetchUserByApiKey(key, state)
	if err != nil {
		return nil, false
	}

	return user, true
}

func UserAuthenticationDynamic(email string, password string, key string, state *app.State) (*database.User, bool) {
	if key != "" {
		return UserKeyAuthentication(key, state)
	} else if password != "" {
		return UserPasswordAuthentication(email, password, state)
	}
	return nil, false
}
