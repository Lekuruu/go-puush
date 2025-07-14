package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Lekuruu/go-puush/internal/database"
)

type AuthenticationRequest struct {
	Username string
	Password string
	Key      string
}

func NewAuthenticationRequest(request *http.Request) (*AuthenticationRequest, error) {
	err := request.ParseForm()
	if err != nil {
		return nil, err
	}

	username := request.FormValue("e")
	password := request.FormValue("p")
	key := request.FormValue("k")
	poop := request.FormValue("z")

	if poop != "poop" {
		return nil, errors.New("invalid request parameter 'z'")
	}

	if username == "" {
		return nil, errors.New("username is required")
	}

	if password == "" && key == "" {
		return nil, errors.New("either password or key must be provided")
	}

	return &AuthenticationRequest{
		Username: username,
		Password: password,
		Key:      key,
	}, nil
}

type AuthenticationResponse struct {
	AccountType        database.AccountType
	ApiKey             string
	DiskUsage          int64
	SubscriptionExpiry *time.Time
}

func (response *AuthenticationResponse) Serialize() []byte {
	var expiry string = ""

	if response.SubscriptionExpiry != nil {
		// Format in "Mon, 02 Jan 2006 15:04:05 MST" format
		expiry = response.SubscriptionExpiry.Format(time.RFC1123)
	}

	data := []string{
		strconv.Itoa(int(response.AccountType)),
		response.ApiKey,
		strconv.Itoa(int(response.DiskUsage)),
		expiry,
	}
	return []byte(strings.Join(data, ","))
}

func PuushAuthentication(ctx *Context) {
	request, err := NewAuthenticationRequest(ctx.Request)
	if err != nil {
		WritePuushError(ctx, -2, http.StatusBadRequest)
		return
	}

	user, success := UserAuthenticationDynamic(
		request.Username,
		request.Password,
		request.Key,
		ctx.State,
	)

	if !success {
		WritePuushError(ctx, -1, http.StatusUnauthorized)
		return
	}

	response := &AuthenticationResponse{
		AccountType:        user.Type,
		ApiKey:             user.ApiKey,
		DiskUsage:          user.DiskUsage,
		SubscriptionExpiry: user.SubscriptionEnd,
	}

	ctx.Response.WriteHeader(http.StatusOK)
	_, err = ctx.Response.Write(response.Serialize())
	if err != nil {
		WritePuushError(ctx, -3, http.StatusInternalServerError)
		return
	}
}
