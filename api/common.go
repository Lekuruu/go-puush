package api

import (
	"errors"
	"strconv"

	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/internal/database"
	"github.com/Lekuruu/go-puush/internal/services"
)

// UserAuthenticationFromContext attempts to authenticate a user based on the API key provided in the request.
func UserAuthenticationFromContext(ctx *Context) (*database.User, error) {
	err := ctx.Request.ParseForm()
	if err != nil {
		return nil, err
	}

	key := ctx.Request.FormValue("k")
	if key == "" {
		return nil, errors.New("missing api key")
	}

	user, err := services.FetchUserByApiKey(key, ctx.State)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UserPasswordAuthentication attempts to authenticate a user using their username and password.
func UserPasswordAuthentication(username string, password string, state *app.State) (*database.User, bool) {
	user, err := services.FetchUserByName(username, state)
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
	user, err := services.FetchUserByName(username, state)
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

// WritePuushError writes the given status/error code to the response.
func WritePuushError(ctx *Context, errorCode int, statusCode int) {
	ctx.Response.WriteHeader(statusCode)
	ctx.Response.Write([]byte(strconv.Itoa(int(errorCode)) + "\n"))
}
