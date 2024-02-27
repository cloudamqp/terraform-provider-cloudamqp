package sanitizer

import (
	"fmt"
	"net/url"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func Fields(jsonBody string) string {
	blockedFields := []string{"apikey", "password"}
	for _, field := range blockedFields {
		if gjson.Get(jsonBody, field).Exists() {
			jsonBody, _ = sjson.Set(jsonBody, field, "***")
		}
	}
	return jsonBody
}

func URL(jsonBody string) string {
	urlFields := []string{"url", "urls.external", "urls.internal"}
	for _, urlField := range urlFields {
		field := gjson.Get(jsonBody, urlField)
		if field.Exists() {
			u, _ := url.Parse(field.String())
			sanitizedUrl := fmt.Sprintf("%s://%s:***@%s%s", u.Scheme, u.User.Username(), u.Host, u.Path)
			jsonBody, _ = sjson.Set(jsonBody, urlField, sanitizedUrl)
		}
	}
	return jsonBody
}