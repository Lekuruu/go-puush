package api

func PuushErrorSubmission(ctx *Context) {
	WritePuushError(ctx, NotImplementedError)
}
