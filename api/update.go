package api

import "net/http"

func PuushWindowsUpdate(ctx *Context) {
	txt, err := ctx.State.Storage.ReadUpdateConfigurationWindows()
	if err != nil {
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.Response.Header().Set("Content-Type", "text/plain")
	ctx.Response.WriteHeader(http.StatusOK)
	_, err = ctx.Response.Write([]byte(txt))
	if err != nil {
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func PuushMacOSRssFeed(ctx *Context) {
	rss, err := ctx.State.Storage.ReadUpdateConfigurationMacOS()
	if err != nil {
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.Response.Header().Set("Content-Type", "application/rss+xml")
	ctx.Response.WriteHeader(http.StatusOK)
	_, err = ctx.Response.Write([]byte(rss))
	if err != nil {
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}
}
