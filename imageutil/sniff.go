// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package imageutil

import (
	"bytes"
	"io"
)

// The algorithm uses at most sniffLen bytes to make its decision.
const sniffLen = 14

// DetectImage detects the image type of the data.
//
// It's a copy of http.DetectContentType with some unnecessary parts removed,
// making it about 10 times faster.
//
// It returns an empty string if the image type can't be determined.
func DetectImage(data []byte) string {
	if len(data) > sniffLen {
		data = data[:sniffLen]
	}

	// Index of the first non-whitespace byte in data.
	firstNonWS := 0
	//for ; firstNonWS < len(data) && isWS(data[firstNonWS]); firstNonWS++ {
	//}

	for _, sig := range sniffSignatures {
		if ct := sig.match(data, firstNonWS); ct != "" {
			return ct
		}
	}

	return ""
}

// DetectImageStream works like DetectImage, but operates on a file stream. It
// reads the minimal amount of data needed and Seek() back to where the offset
// was.
func DetectImageStream(fp io.ReadSeeker) (string, error) {
	head := make([]byte, sniffLen)
	n, err := fp.Read(head)
	if err != nil {
		return "", err
	}

	ct := DetectImage(head)

	_, err = fp.Seek(int64(-n), io.SeekCurrent)
	if err != nil {
		return "", err
	}

	return ct, nil
}

type sniffSig interface {
	// match returns the MIME type of the data, or "" if unknown.
	match(data []byte, firstNonWS int) string
}

// Data matching the table in section 6.
var sniffSignatures = []sniffSig{
	&exactSig{[]byte("GIF87a"), "image/gif"},
	&exactSig{[]byte("GIF89a"), "image/gif"},
	&exactSig{[]byte("\x89\x50\x4E\x47\x0D\x0A\x1A\x0A"), "image/png"},
	&exactSig{[]byte("\xFF\xD8\xFF"), "image/jpeg"},
	&exactSig{[]byte("BM"), "image/bmp"},
	&exactSig{[]byte("\x00\x00\x01\x00"), "image/vnd.microsoft.icon"},
	&maskedSig{
		mask: []byte("\xFF\xFF\xFF\xFF\x00\x00\x00\x00\xFF\xFF\xFF\xFF\xFF\xFF"),
		pat:  []byte("RIFF\x00\x00\x00\x00WEBPVP"),
		ct:   "image/webp",
	},
}

type exactSig struct {
	sig []byte
	ct  string
}

func (e *exactSig) match(data []byte, firstNonWS int) string {
	if bytes.HasPrefix(data, e.sig) {
		return e.ct
	}
	return ""
}

type maskedSig struct {
	mask, pat []byte
	//skipWS    bool
	ct string
}

func (m *maskedSig) match(data []byte, firstNonWS int) string {
	// pattern matching algorithm section 6
	// https://mimesniff.spec.whatwg.org/#pattern-matching-algorithm

	//if m.skipWS {
	//	data = data[firstNonWS:]
	//}
	if len(m.pat) != len(m.mask) {
		return ""
	}
	if len(data) < len(m.mask) {
		return ""
	}
	for i, mask := range m.mask {
		db := data[i] & mask
		if db != m.pat[i] {
			return ""
		}
	}
	return m.ct
}
