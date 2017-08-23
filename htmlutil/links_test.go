package htmlutil

import (
	"fmt"
	"testing"
)

func TestFindURLs(t *testing.T) {
	testCases := []struct {
		in  string
		out []string
	}{
		{
			in:  "XD",
			out: nil,
		},
		{
			in:  "<a href='https://google.com'>Google</a>",
			out: []string{"https://google.com"},
		},
		{
			in:  "<a href=''>Google</a>",
			out: []string{},
		},
		{
			in: `<a href="http://test.com">Test</a> a bunch of different stuff.
			Even with another line<p>Paragraphs and stuff</p>
			<a href="http://teamwork.com">Teamwork!</a>`,
			out: []string{"http://test.com", "http://teamwork.com"},
		},
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			result, _ := FindURLs(tc.in)

			if len(result) != len(tc.out) {
				t.Fatalf("Length of results didn't match %v != %v", len(result), len(tc.out))
			}

			for idx, url := range result {
				if url != tc.out[idx] {
					t.Fatalf("Expected %v got %v", url, tc.out[idx])
				}
			}
		})
	}
}
