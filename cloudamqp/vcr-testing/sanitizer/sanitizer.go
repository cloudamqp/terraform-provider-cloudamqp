package sanitizer

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func FilterSensitiveData(jsonBody, value, placeholder string) string {
	if len(value) == 0 {
		return jsonBody
	}
	return strings.ReplaceAll(jsonBody, value, placeholder)
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
