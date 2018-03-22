package httputilx

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"testing"
)

func TestS(t *testing.T) {
	r := Request{&http.Request{
		Method:           http.MethodGet,
		URL:              &url.URL{},
		Proto:            "",
		ProtoMajor:       0,
		ProtoMinor:       0,
		Header:           map[string][]string{},
		Body:             nil,
		GetBody:          nil,
		ContentLength:    0,
		TransferEncoding: nil,
		Close:            false,
		Host:             "",
		Form:             map[string][]string{},
		PostForm:         map[string][]string{},
		MultipartForm:    &multipart.Form{},
		Trailer:          map[string][]string{},
		RemoteAddr:       "",
		RequestURI:       "",
	}}

	j, err := json.MarshalIndent(r, "", "    ")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(j))
}
