package httputilx

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"testing"

	"github.com/teamwork/test"
)

func TestDumpBody(t *testing.T) {
	numg0 := runtime.NumGoroutine()

	cases := []struct {
		Req  http.Request
		Body interface{} // optional []byte or func() io.ReadCloser to populate Req.Body

		WantDump string
		ReadN    int64
		NoBody   bool // if true, set DumpRequest{,Out} body to false
	}{

		// HTTP/1.1 => chunked coding; body; empty trailer
		{
			Req: http.Request{
				Method: "GET",
				URL: &url.URL{
					Scheme: "http",
					Host:   "www.google.com",
					Path:   "/search",
				},
				ProtoMajor:       1,
				ProtoMinor:       1,
				TransferEncoding: []string{"chunked"},
			},
			Body:     []byte("abcdef"),
			WantDump: chunk("abcdef") + chunk(""),
			ReadN:    -1,
		},
		{
			Req: http.Request{
				Method: "GET",
				URL: &url.URL{
					Scheme: "http",
					Host:   "www.google.com",
					Path:   "/search",
				},
				ProtoMajor:       1,
				ProtoMinor:       1,
				TransferEncoding: []string{"chunked"},
			},
			Body:     []byte("abcdef"),
			WantDump: chunk("a") + chunk(""),
			ReadN:    1,
		},
		{
			Req: http.Request{
				Method: "GET",
				URL: &url.URL{
					Scheme: "http",
					Host:   "www.google.com",
					Path:   "/search",
				},
				ProtoMajor:       1,
				ProtoMinor:       1,
				TransferEncoding: []string{"chunked"},
			},
			Body:     []byte("abcdef"),
			WantDump: chunk("ab") + chunk(""),
			ReadN:    2,
		},

		// Request with Body > 8196 (default buffer size)
		{
			Req: http.Request{
				Method: "POST",
				URL: &url.URL{
					Scheme: "http",
					Host:   "post.tld",
					Path:   "/",
				},
				Header: http.Header{
					"Content-Length": []string{"8193"},
				},
				ContentLength: 8193,
				ProtoMajor:    1,
				ProtoMinor:    1,
			},
			Body:     bytes.Repeat([]byte("a"), 8193),
			WantDump: strings.Repeat("a", 8193),
			ReadN:    -1,
		},
		{
			Req: http.Request{
				Method: "POST",
				URL: &url.URL{
					Scheme: "http",
					Host:   "post.tld",
					Path:   "/",
				},
				Header: http.Header{
					"Content-Length": []string{"8193"},
				},
				ContentLength: 8193,
				ProtoMajor:    1,
				ProtoMinor:    1,
			},
			Body:     bytes.Repeat([]byte("a"), 8293),
			WantDump: strings.Repeat("a", 8193),
			ReadN:    8193,
		},
		{
			Req: http.Request{
				Method: "POST",
				URL: &url.URL{
					Scheme: "http",
					Host:   "post.tld",
					Path:   "/",
				},
				Header: http.Header{
					"Content-Length": []string{"10"},
				},
				ContentLength: 10,
				ProtoMajor:    1,
				ProtoMinor:    1,
			},
			Body:     bytes.Repeat([]byte("a"), 10),
			WantDump: strings.Repeat("a", 10),
			ReadN:    15,
		},
		{
			Req: http.Request{
				Method: "GET",
				URL: &url.URL{
					Scheme: "http",
					Host:   "www.google.com",
					Path:   "/search",
				},
				ProtoMajor:       1,
				ProtoMinor:       1,
				TransferEncoding: []string{"chunked"},
			},
			Body:     []byte("abcdef"),
			WantDump: chunk("abcdef") + chunk(""),
			ReadN:    100,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			setBody := func() {
				if tc.Body == nil {
					return
				}
				switch b := tc.Body.(type) {
				case []byte:
					tc.Req.Body = ioutil.NopCloser(bytes.NewReader(b))
				case func() io.ReadCloser:
					tc.Req.Body = b()
				default:
					t.Fatalf("unsupported Body of %T", tc.Body)
				}
			}
			setBody()
			if tc.Req.Header == nil {
				tc.Req.Header = make(http.Header)
			}

			setBody()
			dump, err := DumpBody(&tc.Req, tc.ReadN)
			if err != nil {
				t.Fatal(err)
			}
			if string(dump) != tc.WantDump {
				t.Errorf("\nwant: %#v\ngot:  %#v\n", tc.WantDump, string(dump))
			}
		})
	}

	if dg := runtime.NumGoroutine() - numg0; dg > 4 {
		buf := make([]byte, 4096)
		buf = buf[:runtime.Stack(buf, true)]
		t.Errorf("Unexpectedly large number of new goroutines: %d new: %s", dg, buf)
	}
}

func chunk(s string) string {
	return fmt.Sprintf("%x\r\n%s\r\n", len(s), s)
}

// TODO: better to not depend on interwebz...
func TestFetch(t *testing.T) {
	cases := []struct {
		in, want, wantErr string
	}{
		{"http://example.com", "<html>", ""},
		{"http://fairly-certain-this-doesnt-exist-asdasd12g1ghdfddd.com", "", "cannot download"},
		{"http://httpbin.org/status/400", "", "400"},
		{"http://httpbin.org/status/500", "", "500"},
		// Make sure we return the body as well.
		{"http://httpbin.org/status/418", "teapot", "418"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc.in), func(t *testing.T) {
			out, err := Fetch(tc.in)
			if !test.ErrorContains(err, tc.wantErr) {
				t.Errorf("wrong error\nout:  %#v\nwant: %#v\n", err, tc.wantErr)
			}
			if !strings.Contains(string(out), tc.want) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", string(out), tc.want)
			}
		})
	}
}

// TODO: Add test.
func TestSave(t *testing.T) {
}
