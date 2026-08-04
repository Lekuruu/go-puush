package database

import (
	"net/url"
	"strings"
)

func appendDatabaseOptions(path string, options url.Values) string {
	separator := "?"
	if strings.Contains(path, "?") {
		separator = "&"
	}
	if strings.HasSuffix(path, "?") || strings.HasSuffix(path, "&") {
		separator = ""
	}
	return path + separator + options.Encode()
}
