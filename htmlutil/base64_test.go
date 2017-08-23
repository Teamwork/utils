package htmlutil

import (
	"fmt"
	"strings"
	"testing"

	"github.com/teamwork/test"
)

func TestStripBase64Images(t *testing.T) {
	var cases = []struct {
		inFile   string
		expected string
	}{
		{"one", `With an image <img src="http://testing.com/image.png"  />`},
		{"two", `With an image and single quotes <img src='http://testing.com/image.png'  />`},
		{"three", `With an image and multiple attrs <img src='http://testing.com/image.png'  alt='hey there' />`},
		{"four", `With image and no data <img src="http://testing.com/image.png" />`},
		{"five", `With no image`},
		{"six", `data-src without base64 <img src="http://testing.com/image.png" data-src="github.com" />`},
		{"seven", string(test.Read(t, "./base64_test", "seven_expected.txt"))},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			out := strings.TrimSpace(StripBase64Images(string(test.Read(t, "./base64_test", tc.inFile+".txt"))))
			tc.expected = strings.TrimSpace(tc.expected)
			if out != tc.expected {
				t.Errorf("\nout:      %#v\nexpected: %#v\n", out, tc.expected)
			}
		})
	}
}
