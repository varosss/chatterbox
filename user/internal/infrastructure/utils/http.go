package utils

import (
	"net/http"
	"strings"
)

func ParseSameSite(v string) http.SameSite {
	switch strings.ToLower(v) {

	case "none":
		return http.SameSiteNoneMode

	case "lax":
		return http.SameSiteLaxMode

	case "strict":
		return http.SameSiteStrictMode

	default:
		return http.SameSiteLaxMode
	}
}
