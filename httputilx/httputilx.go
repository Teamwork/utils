// Package httputilx provides HTTP utility functions.
package httputilx // import "github.com/teamwork/utils/httputilx"

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/teamwork/utils/ioutilx"
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
	data, err := ioutil.ReadAll(response.Body)
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
