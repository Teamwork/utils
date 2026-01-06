package httputilx

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"runtime"
	"strings"
	"testing"
	"time"

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
					tc.Req.Body = io.NopCloser(bytes.NewReader(b))
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
		{"http://example.com", "<html", ""},
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

func TestDoExponentialBackoff(t *testing.T) {
	tests := []struct {
		name         string
		options      []ExponentialBackoffOption
		handler      http.HandlerFunc
		requestBody  io.Reader
		wantBody     string
		wantErr      string
		wantAttempts int
	}{
		{
			name: "Success",
			options: []ExponentialBackoffOption{
				ExponentialBackoffWithConfig(3, 100*time.Millisecond, 5*time.Second, 2.0),
			},
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("ok"))
			},
			wantBody:     "ok",
			wantErr:      "",
			wantAttempts: 1,
		},
		{
			name: "RetryOnError",
			options: []ExponentialBackoffOption{
				ExponentialBackoffWithConfig(3, 100*time.Millisecond, 5*time.Second, 2.0),
			},
			handler: func() http.HandlerFunc {
				attempts := 0
				return func(w http.ResponseWriter, _ *http.Request) {
					attempts++
					if attempts < 3 {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte("success"))
				}
			}(),
			wantBody:     "success",
			wantErr:      "",
			wantAttempts: 3,
		},
		{
			name: "TooManyRequests",
			options: []ExponentialBackoffOption{
				ExponentialBackoffWithConfig(3, 100*time.Millisecond, 5*time.Second, 2.0),
			},
			handler: func() http.HandlerFunc {
				attempts := 0
				return func(w http.ResponseWriter, _ *http.Request) {
					attempts++
					if attempts < 2 {
						w.WriteHeader(http.StatusTooManyRequests)
						return
					}
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte("done"))
				}
			}(),
			wantBody:     "done",
			wantErr:      "",
			wantAttempts: 2,
		},
		{
			name: "MaxRetriesExceeded",
			options: []ExponentialBackoffOption{
				ExponentialBackoffWithConfig(2, 100*time.Millisecond, 5*time.Second, 2.0),
			},
			handler: func() http.HandlerFunc {
				return func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}(),
			wantBody:     "",
			wantErr:      "",
			wantAttempts: 3,
		},
		{
			name: "CustomShouldRetry",
			options: []ExponentialBackoffOption{
				ExponentialBackoffWithConfig(2, 100*time.Millisecond, 5*time.Second, 2.0),
				ExponentialBackoffWithShouldRetry(func(resp *http.Response, _ error) bool {
					if resp != nil && resp.StatusCode == http.StatusBadRequest {
						return true
					}
					return false
				}),
			},
			handler: func() http.HandlerFunc {
				return func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusBadRequest)
				}
			}(),
			wantBody:     "",
			wantErr:      "",
			wantAttempts: 3,
		},
		{
			name: "RequestBodyCopiedOnRetry",
			options: []ExponentialBackoffOption{
				ExponentialBackoffWithConfig(4, 100*time.Millisecond, 5*time.Second, 2.0),
			},
			handler: func() http.HandlerFunc {
				initialBody := "request body content"

				attempts := 0
				return func(w http.ResponseWriter, r *http.Request) {
					attempts++

					if r.ContentLength != int64(len(initialBody)) {
						w.WriteHeader(http.StatusInternalServerError)
						_, _ = fmt.Fprintf(w, "wrong content-length: got %d, want %d", r.ContentLength, len(initialBody))
						return
					}

					body, err := io.ReadAll(r.Body)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					if len(body) != len(initialBody) {
						w.WriteHeader(http.StatusInternalServerError)
						_, _ = fmt.Fprintf(w, "content-length mismatch: header=%d actual=%d", r.ContentLength, len(body))
						return
					}

					// Verify body is correctly sent on all attempts
					if string(body) != initialBody {
						w.WriteHeader(http.StatusInternalServerError)
						_, _ = fmt.Fprintf(w, "incorrect body: %q", string(body))
						return
					}
					if attempts < 3 {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte("body received correctly"))
				}
			}(),
			requestBody:  io.NopCloser(bytes.NewBuffer([]byte("request body content"))),
			wantBody:     "body received correctly",
			wantErr:      "",
			wantAttempts: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attempts := 0
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				attempts++
				tt.handler(w, r)
			}))
			defer ts.Close()

			req, err := http.NewRequest(http.MethodGet, ts.URL, tt.requestBody)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			resp, err := DoExponentialBackoff(req, tt.options...)
			if err == nil {
				defer resp.Body.Close() //nolint:errcheck
			}
			if tt.wantErr != "" {
				if err == nil {
					t.Errorf("expected error %q, got nil", tt.wantErr)
				} else if !test.ErrorContains(err, tt.wantErr) {
					t.Errorf("expected error %q, got %q", tt.wantErr, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %s", err.Error())
				} else {
					body, _ := io.ReadAll(resp.Body)
					if tt.wantBody != string(body) {
						t.Errorf("expected body %q, got %q", tt.wantBody, string(body))
					}
				}
			}

			if attempts != tt.wantAttempts {
				t.Errorf("expected %d attempts, got %d", tt.wantAttempts, attempts)
			}
		})
	}
}
