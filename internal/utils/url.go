package utils

import (
	"fmt"
	"net/url"
)

// AddQueryParams adds query parameters to a URL.
func AddQueryParams(u string, params map[string]interface{}) string {
	if len(params) == 0 {
		return u
	}

	if u == "" {
		u = "/"
	}

	uu, err := url.Parse(u)
	if err != nil {
		return u
	}

	q := uu.Query()
	for k, v := range params {
		switch v.(type) {
		case []string:
			for _, vv := range v.([]string) {
				q.Add(k, vv)
			}
		default:
			q.Add(k, fmt.Sprintf("%v", v))
		}
	}
	uu.RawQuery = q.Encode()

	return uu.String()
}

// AddQueryParam adds a query parameter to a URL.
func AddQueryParam(u, k string, v interface{}) string {
	return AddQueryParams(u, map[string]interface{}{k: v})
}
