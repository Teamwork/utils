// Package ioutilx implements some I/O utility functions.
package ioutilx // import "github.com/teamwork/utils/ioutilx"

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

// DumpReader reads all of b to memory and then returns two equivalent
// ReadClosers which will yield the same bytes.
//
// This is useful if you want to read data from an io.Reader more than once.
//
// It returns an error if the initial reading of all bytes fails. It does not
// attempt to make the returned ReadClosers have identical error-matching
// behavior.
//
// This is based on httputil.DumpRequest, see
// github.com/teamwork/ioutilx.DumpBody() for an example usage.
//
// Copyright 2009 The Go Authors. All rights reserved. Use of this source code
// is governed by a BSD-style license that can be found in the LICENSE file:
// https://golang.org/LICENSE
func DumpReader(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if b == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return http.NoBody, http.NoBody, nil
	}

	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}

	if err = b.Close(); err != nil {
		return nil, b, err
	}

	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}
