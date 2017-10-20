// Package httputilx provides HTTP utility functions.
package httputilx // import "github.com/teamwork/utils/httputilx"

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httputil"

	"github.com/teamwork/utils/ioutilx"
)

// DumpBody reads the body of a HTTP request without consuming it, so it can be
// read again later.
//
// It's based on httputil.DumpRequest.
//
// Copyright 2009 The Go Authors. All rights reserved. Use of this source code
// is governed by a BSD-style license that can be found in the LICENSE file:
// https://golang.org/LICENSE
func DumpBody(r *http.Request) ([]byte, error) {
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
	_, err = io.Copy(dest, body)
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
