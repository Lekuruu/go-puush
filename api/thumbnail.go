package api

// /api/thumb returns a thumbnail for the given upload ID.
// If the thumbnail is not available, it tries to generate one.
func PuushThumbnail(ctx *Context) {
	WritePuushError(ctx, NotImplementedError)
}
