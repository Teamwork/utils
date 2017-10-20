package httputilx

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"testing"
)

type dumpTest struct {
	Req  http.Request
	Body interface{} // optional []byte or func() io.ReadCloser to populate Req.Body

	WantDump string
	NoBody   bool // if true, set DumpRequest{,Out} body to false
}

var dumpTests = []dumpTest{

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
	},
}

func TestDumpRequest(t *testing.T) {
	numg0 := runtime.NumGoroutine()
	for i, tt := range dumpTests {
		setBody := func() {
			if tt.Body == nil {
				return
			}
			switch b := tt.Body.(type) {
			case []byte:
				tt.Req.Body = ioutil.NopCloser(bytes.NewReader(b))
			case func() io.ReadCloser:
				tt.Req.Body = b()
			default:
				t.Fatalf("Test %d: unsupported Body of %T", i, tt.Body)
			}
		}
		setBody()
		if tt.Req.Header == nil {
			tt.Req.Header = make(http.Header)
		}

		setBody()
		dump, err := DumpBody(&tt.Req)
		if err != nil {
			t.Errorf("DumpBody #%d: %s", i, err)
			continue
		}
		if string(dump) != tt.WantDump {
			t.Errorf("DumpBody %d, expecting:\n%s\nGot:\n%s\n", i, tt.WantDump, string(dump))
			continue
		}
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

func mustParseURL(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(fmt.Sprintf("Error parsing URL %q: %v", s, err))
	}
	return u
}

func mustNewRequest(method, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(fmt.Sprintf("NewRequest(%q, %q, %p) err = %v", method, url, body, err))
	}
	return req
}

func mustReadRequest(s string) *http.Request {
	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(s)))
	if err != nil {
		panic(err)
	}
	return req
}
