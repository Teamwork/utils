package header

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/teamwork/test"
)

func TestSetContentDisposition(t *testing.T) {
	cases := []struct {
		args          DispositionArgs
		want, wantErr string
	}{
		{DispositionArgs{}, "", "mandatory"},
		{DispositionArgs{Type: "ASD"}, "", "must be"},

		{DispositionArgs{Type: TypeInline},
			"inline", ""},
		{DispositionArgs{Type: TypeInline, Filename: "hello.pdf"},
			`inline; filename="hello.pdf"`, ""},
		{DispositionArgs{Type: TypeAttachment, Filename: "hello.pdf"},
			`attachment; filename="hello.pdf"`, ""},
		{DispositionArgs{Type: TypeInline, Filename: `hello, "world".pdf`},
			`inline; filename="hello, \"world\".pdf"`, ""},
		{DispositionArgs{Type: TypeInline, Filename: `hello, "world".pdf/%20\`},
			`inline; filename="hello, \"world\".pdf/20"`, ""},
		{DispositionArgs{Type: TypeInline, Filename: `h€llo.pdf`},
			`inline; filename="hllo.pdf"; filename*=UTF-8''h%E2%82%ACllo.pdf`, ""},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			h := http.Header{}
			err := SetContentDisposition(h, tc.args)
			if !test.ErrorContains(err, tc.wantErr) {
				t.Errorf("wrong error\nout:  %s\nwant: %s\n", err, tc.wantErr)
			}

			out := h.Get("Content-Disposition")
			if out != tc.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}
