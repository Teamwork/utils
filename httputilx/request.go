package httputilx

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/url"
)

// Request embeds http.Request and adds methods for serialisation.
type Request struct {
	*http.Request
}

// MarshalJSON ..
func (r Request) MarshalJSON() ([]byte, error) {
	//r.URL
	//r.Header

	return json.Marshal(map[string]interface{}{
		"Method":           r.Method,
		"URL":              &url.URL{},
		"Proto":            "",
		"ProtoMajor":       0,
		"ProtoMinor":       0,
		"Header":           map[string][]string{},
		"Body":             nil,
		"GetBody":          nil,
		"ContentLength":    0,
		"TransferEncoding": nil,
		"Close":            false,
		"Host":             "",
		"Form":             map[string][]string{},
		"PostForm":         map[string][]string{},
		"MultipartForm":    &multipart.Form{},
		"Trailer":          map[string][]string{},
		"RemoteAddr":       "",
		"RequestURI":       "",
	})
}

// Header embeds http.Header and adds methods for serialisation.
type Header struct {
	http.Header
}

// Values embeds url.Values and adds methods for serialisation.
type Values struct {
	url.Values
}
