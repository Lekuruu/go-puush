package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/internal/database"
	"github.com/Lekuruu/go-puush/internal/services"
)

// PuushError represents an error that can occur in the Puush API.
type PuushError struct {
	PuushCode  int
	StatusCode int
}

func (e PuushError) Error() string {
	return "PuushError: " + strconv.Itoa(e.PuushCode) + " (HTTP " + strconv.Itoa(e.StatusCode) + ")"
}

var (
	AuthenticationFailure    PuushError = PuushError{-1, http.StatusUnauthorized}
	RequestError             PuushError = PuushError{-2, http.StatusBadRequest}
	ServerError              PuushError = PuushError{-2, http.StatusInternalServerError}
	UnknownError             PuushError = PuushError{-2, http.StatusInternalServerError}
	NotImplementedError      PuushError = PuushError{-2, http.StatusNotImplemented}
	ChecksumError            PuushError = PuushError{-3, http.StatusBadRequest}
	InsufficientStorageError PuushError = PuushError{-4, http.StatusPaymentRequired}
)

// WritePuushError writes the given puush error struct to the response.
func WritePuushError(ctx *Context, error PuushError) {
	ctx.Response.WriteHeader(error.StatusCode)
	ctx.Response.Write([]byte(strconv.Itoa(error.PuushCode) + "\n"))
}

// UserAuthenticationFromKey attempts to authenticate a user using the provided API key.
func UserAuthenticationFromKey(key string, state *app.State, preload ...string) (*database.User, error) {
	user, err := services.FetchUserByApiKey(key, state, preload...)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UserAuthenticationFromContext attempts to authenticate a user based on the API key provided in the request.
func UserAuthenticationFromContext(ctx *Context, preload ...string) (*database.User, error) {
	err := ctx.Request.ParseForm()
	if err != nil {
		return nil, err
	}

	key := ctx.Request.FormValue("k")
	if key == "" {
		return nil, errors.New("missing api key")
	}

	return UserAuthenticationFromKey(key, ctx.State, preload...)
}

// UserPasswordAuthentication attempts to authenticate a user using their username and password.
func UserPasswordAuthentication(username string, password string, state *app.State) (*database.User, bool) {
	user, err := services.FetchUserByNameOrEmail(username, state)
	if err != nil {
		return nil, false
	}

	if !app.VerifyPasswordHash(password, user.Password) {
		return nil, false
	}

	if !user.Active {
		return nil, false
	}

	return user, true
}

// UserKeyAuthentication attempts to authenticate a user using an API key.
func UserKeyAuthentication(username string, key string, state *app.State) (*database.User, bool) {
	user, err := services.FetchUserByNameOrEmail(username, state)
	if err != nil {
		return nil, false
	}

	if user.ApiKey != key {
		return nil, false
	}

	if !user.Active {
		return nil, false
	}

	return user, true
}

// UserAuthenticationDynamic attempts to authenticate a user using either a password or an API key.
func UserAuthenticationDynamic(username string, password string, key string, state *app.State) (*database.User, bool) {
	if password != "" {
		return UserPasswordAuthentication(username, password, state)
	} else if key != "" {
		return UserKeyAuthentication(username, key, state)
	}
	return nil, false
}
