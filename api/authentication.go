package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/internal/database"
)

// /api/auth handles user authentication for the puush api service.
// It supports both password & api key authentication and returns
// the user's account type, api key, disk usage, and subscription expiry.
func PuushAuthentication(ctx *app.Context) {
	request, err := NewAuthenticationRequest(ctx.Request)
	if err != nil {
		WritePuushError(ctx, RequestError)
		return
	}

	user, success := UserAuthenticationDynamic(
		request.Username,
		request.Password,
		request.Key,
		ctx.State,
	)

	if !success {
		WritePuushError(ctx, AuthenticationFailure)
		return
	}

	response := &AuthenticationResponse{
		AccountType:        user.Type,
		ApiKey:             user.ApiKey,
		DiskUsage:          user.DiskUsage,
		SubscriptionExpiry: user.SubscriptionEnd,
	}
	WritePuushResponse(ctx, response)
}

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
		// Format in "2006-01-02 15:04:05" format
		expiry = response.SubscriptionExpiry.Format(time.DateTime)
	}

	data := []string{
		strconv.Itoa(int(response.AccountType)),
		response.ApiKey,
		expiry,
		strconv.Itoa(int(response.DiskUsage)),
	}
	return []byte(strings.Join(data, ","))
}
