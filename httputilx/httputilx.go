// Package httputilx provides HTTP utility functions.
package httputilx // import "github.com/teamwork/utils/v2/httputilx"

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/teamwork/utils/v2/ioutilx"
)

// DumpBody reads the body of a HTTP request without consuming it, so it can be
// read again later.
// It will read at most maxSize of bytes. Use -1 to read everything.
//
// It's based on httputil.DumpRequest.
//
// Copyright 2009 The Go Authors. All rights reserved. Use of this source code
// is governed by a BSD-style license that can be found in the LICENSE file:
// https://golang.org/LICENSE
func DumpBody(r *http.Request, maxSize int64) ([]byte, error) {
	if r.Body == nil {
		return nil, nil
	}

	save, body, err := ioutilx.DumpReader(r.Body)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	var dest io.Writer = &b

	chunked := len(r.TransferEncoding) > 0 && r.TransferEncoding[0] == "chunked"
	if chunked {
		dest = httputil.NewChunkedWriter(dest)
	}

	if maxSize < 0 {
		_, err = io.Copy(dest, body)
	} else {
		_, err = io.CopyN(dest, body, maxSize)
		if err == io.EOF {
			err = nil
		}
	}
	if err != nil {
		return nil, err
	}
	if chunked {
		_ = dest.(io.Closer).Close()
		_, _ = io.WriteString(&b, "\r\n")
	}

	r.Body = save
	return b.Bytes(), nil
}

// ErrNotOK is used when the status code is not 200 OK.
type ErrNotOK struct {
	URL string
	Err string
}

func (e ErrNotOK) Error() string {
	return fmt.Sprintf("code %v while downloading %v", e.Err, e.URL)
}

// Fetch the contents of an HTTP URL.
//
// This is not intended to cover all possible use cases  for fetching files,
// only the most common ones. Use the net/http package for more advanced usage.
func Fetch(url string) ([]byte, error) {
	client := http.Client{Timeout: 60 * time.Second}
	response, err := client.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot download %v", url)
	}
	defer response.Body.Close() // nolint: errcheck

	// TODO: Maybe add sanity check to bail out of the Content-Length is very
	// large?
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read body of %v", url)
	}

	if response.StatusCode != http.StatusOK {
		return data, ErrNotOK{
			URL: url,
			Err: fmt.Sprintf("%v %v", response.StatusCode, response.Status),
		}
	}

	return data, nil
}

// Save an HTTP URL to the directory dir with the filename. The filename can be
// generated from the URL if empty.
//
// It will return the full path to the save file. Note that it may create both a
// file *and* return an error (e.g. in cases of non-200 status codes).
//
// This is not intended to cover all possible use cases  for fetching files,
// only the most common ones. Use the net/http package for more advanced usage.
func Save(url string, dir string, filename string) (string, error) {
	// Use last path of url if filename is empty
	if filename == "" {
		tokens := strings.Split(url, "/")
		filename = tokens[len(tokens)-1]
	}
	path := filepath.FromSlash(dir + "/" + filename)

	client := http.Client{Timeout: 60 * time.Second}
	response, err := client.Get(url)
	if err != nil {
		return "", errors.Wrapf(err, "cannot download %v", url)
	}
	defer response.Body.Close() // nolint: errcheck

	output, err := os.Create(path)
	if err != nil {
		return "", errors.Wrapf(err, "cannot create %v", path)
	}
	defer output.Close() // nolint: errcheck

	_, err = io.Copy(output, response.Body)
	if err != nil {
		return path, errors.Wrapf(err, "cannot read body of %v in to %v", url, path)
	}

	if response.StatusCode != http.StatusOK {
		return path, ErrNotOK{
			URL: url,
			Err: fmt.Sprintf("%v %v", response.StatusCode, response.Status),
		}
	}

	return path, nil
}

// ExponentialBackoffOptions contains options for the exponential backoff retry
// mechanism.
type ExponentialBackoffOptions struct {
	client            *http.Client
	maxRetries        int
	initialBackoff    time.Duration
	maxBackoff        time.Duration
	backoffMultiplier float64
	shouldRetry       func(resp *http.Response, err error) bool
	logger            *slog.Logger
}

// ExponentialBackoffOption is a function that configures
// ExponentialBackoffOptions.
type ExponentialBackoffOption func(*ExponentialBackoffOptions)

// ExponentialBackoffWithClient sets the HTTP client to be used when sending the
// API requests. By default, http.DefaultClient is used.
func ExponentialBackoffWithClient(client *http.Client) ExponentialBackoffOption {
	return func(o *ExponentialBackoffOptions) {
		o.client = client
	}
}

// ExponentialBackoffWithConfig sets the configuration for the exponential
// backoff retry mechanism. By default, it will retry up to 3 times, starting
// with a 100ms backoff, doubling each time up to a maximum of 5s.
func ExponentialBackoffWithConfig(
	maxRetries int,
	initialBackoff, maxBackoff time.Duration,
	backoffMultiplier float64,
) ExponentialBackoffOption {
	return func(o *ExponentialBackoffOptions) {
		o.maxRetries = maxRetries
		o.initialBackoff = initialBackoff
		o.maxBackoff = maxBackoff
		o.backoffMultiplier = backoffMultiplier
	}
}

// ExponentialBackoffWithShouldRetry sets the function to determine whether a
// request should be retried based on the response and error. By default, it
// retries on any error, as well as on HTTP 5xx and 429 status codes.
func ExponentialBackoffWithShouldRetry(
	shouldRetry func(resp *http.Response, err error) bool,
) ExponentialBackoffOption {
	return func(o *ExponentialBackoffOptions) {
		o.shouldRetry = shouldRetry
	}
}

// ExponentialBackoffWithLogger sets the logger to be used for logging retry
// attempts. By default, a no-op logger is used.
func ExponentialBackoffWithLogger(logger *slog.Logger) ExponentialBackoffOption {
	return func(o *ExponentialBackoffOptions) {
		o.logger = logger
	}
}

// DoExponentialBackoff will send an API request using exponential backoff until
// it either succeeds or the maximum number of retries is reached.
func DoExponentialBackoff(req *http.Request, options ...ExponentialBackoffOption) (*http.Response, error) {
	o := ExponentialBackoffOptions{
		client:            http.DefaultClient,
		maxRetries:        3,
		initialBackoff:    100 * time.Millisecond,
		maxBackoff:        5 * time.Second,
		backoffMultiplier: 2.0,
		shouldRetry: func(resp *http.Response, err error) bool {
			if err != nil {
				return true
			}
			if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
				return true
			}
			return false
		},
		logger: slog.New(slog.DiscardHandler),
	}
	for _, option := range options {
		option(&o)
	}

	backoff := o.initialBackoff

	for attempt := 0; attempt <= o.maxRetries; attempt++ {
		reqClone := req.Clone(req.Context())
		if req.Body != nil {
			if seeker, ok := req.Body.(interface {
				Seek(int64, int) (int64, error)
			}); ok {
				_, _ = seeker.Seek(0, 0)
			}
			reqClone.Body = req.Body
		}

		resp, err := o.client.Do(reqClone)
		if !o.shouldRetry(resp, err) || attempt >= o.maxRetries {
			return resp, nil
		}

		logArgs := []any{
			slog.Int("attempt", attempt+1),
			slog.Duration("backoff", backoff),
		}
		if err != nil {
			logArgs = append(logArgs, slog.String("error", err.Error()))
		}
		if resp != nil {
			if err := resp.Body.Close(); err != nil {
				o.logger.Error("failed to close response body",
					slog.Int("attempt", attempt+1),
					slog.String("error", err.Error()),
				)
			}
			logArgs = append(logArgs, slog.Int("status_code", resp.StatusCode))
		}

		o.logger.Debug("request failed", logArgs...)
		time.Sleep(backoff)
		backoff = min(time.Duration(float64(backoff)*o.backoffMultiplier), o.maxBackoff)
	}

	return nil, fmt.Errorf("request failed after %d attempts", o.maxRetries+1)
}
