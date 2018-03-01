package goutil

import (
	"fmt"
	"go/build"
	"go/token"
	"reflect"
	"sort"
	"testing"

	"github.com/teamwork/test"
)

// This also tests ResolvePackage() and ResolveWildcard().
func TestExpand(t *testing.T) {
	cases := []struct {
		in      []string
		want    []string
		wantErr string
	}{
		{
			[]string{"fmt"},
			[]string{"fmt"},
			"",
		},
		{
			[]string{"fmt", "fmt"},
			[]string{"fmt"},
			"",
		},
		{
			[]string{"fmt", "net/http"},
			[]string{"fmt", "net/http"},
			"",
		},
		{
			[]string{"net/..."},
			[]string{"net", "net/http", "net/http/cgi", "net/http/cookiejar",
				"net/http/fcgi", "net/http/httptest", "net/http/httptrace",
				"net/http/httputil", "net/http/internal", "net/http/pprof",
				"net/internal/socktest", "net/mail", "net/rpc", "net/rpc/jsonrpc",
				"net/smtp", "net/textproto", "net/url",
			},
			"",
		},
		{
			[]string{"github.com/teamwork/utils"},
			[]string{"github.com/teamwork/utils"},
			"",
		},
		{
			[]string{"."},
			[]string{"github.com/teamwork/utils/goutil"},
			"",
		},
		{
			[]string{".."},
			[]string{"github.com/teamwork/utils"},
			"",
		},
		{
			[]string{"../..."},
			[]string{
				"github.com/teamwork/utils",
				"github.com/teamwork/utils/aesutil",
				"github.com/teamwork/utils/byteutil",
				"github.com/teamwork/utils/goutil",
				"github.com/teamwork/utils/httputilx",
				"github.com/teamwork/utils/httputilx/header",
				"github.com/teamwork/utils/ioutilx",
				"github.com/teamwork/utils/jsonutil",
				"github.com/teamwork/utils/maputil",
				"github.com/teamwork/utils/mathutil",
				"github.com/teamwork/utils/netutil",
				"github.com/teamwork/utils/sliceutil",
				"github.com/teamwork/utils/sqlutil",
				"github.com/teamwork/utils/stringutil",
				"github.com/teamwork/utils/syncutil",
				"github.com/teamwork/utils/timeutil",
			},
			"",
		},

		// Errors
		{
			[]string{""},
			nil,
			"cannot resolve empty string",
		},
		{
			[]string{"this/will/never/exist"},
			nil,
			`cannot find package "this/will/never/exist"`,
		},
		{
			[]string{"this/will/never/exist/..."},
			nil,
			`cannot find package "this/will/never/exist"`,
		},
		{
			[]string{"./doesnt/exist"},
			nil,
			"cannot find package",
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			out, err := Expand(tc.in, build.FindOnly)
			if !test.ErrorContains(err, tc.wantErr) {
				t.Fatal(err)
			}

			sort.Strings(tc.want)
			var outPkgs []string
			for _, p := range out {
				outPkgs = append(outPkgs, p.ImportPath)
			}

			if !reflect.DeepEqual(tc.want, outPkgs) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", outPkgs, tc.want)
			}
		})
	}
}

func TestParseFiles(t *testing.T) {
	pkg, err := ResolvePackage("net/http", 0)
	if err != nil {
		t.Fatal(err)
	}

	fset := token.NewFileSet()
	out, err := ParseFiles(fset, pkg.Dir, pkg.GoFiles, 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(out) != 1 {
		t.Fatalf("len(out) == %v", len(out))
	}

	for _, pkg := range out {
		if pkg.Name != "http" {
			t.Errorf("name == %v", pkg.Name)
		}

		if len(pkg.Files) < 10 {
			t.Errorf("len(pkg.Files) == %v", len(pkg.Files))
		}
	}
}
