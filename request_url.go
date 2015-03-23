package pusher

import (
	"net/url"
	"strings"
)

const auth_version = "1.0"

func unsigned_params(key, timestamp string, body []byte, additional_queries map[string]string) url.Values {
	params := url.Values{
		"auth_key":       {key},
		"auth_timestamp": {timestamp},
		"auth_version":   {auth_version},
	}

	if body != nil {
		params.Add("body_md5", md5Signature(body))
	}

	if additional_queries != nil {
		for key, values := range additional_queries {
			params.Add(key, values)
		}
	}

	return params

}

func unescape_url(_url url.Values) string {
	unesc, _ := url.QueryUnescape(_url.Encode())
	return unesc
}

func createRequestUrl(method, host, path, key, secret, timestamp string, body []byte, additional_queries map[string]string) string {
	params := unsigned_params(key, timestamp, body, additional_queries)

	string_to_sign := strings.Join([]string{method, path, unescape_url(params)}, "\n")

	auth_signature := hmacSignature(string_to_sign, secret)

	params.Add("auth_signature", auth_signature)

	if host == "" {
		host = "http://api.pusherapp.com"
	}

	endpoint, _ := url.Parse(host + path)

	endpoint.RawQuery = unescape_url(params)

	return endpoint.String()
}
