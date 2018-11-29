package goutil

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/token"
	"reflect"
	"sort"
	"strings"
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
				"github.com/teamwork/utils/errorutil",
				"github.com/teamwork/utils/goutil",
				"github.com/teamwork/utils/httputilx",
				"github.com/teamwork/utils/httputilx/header",
				"github.com/teamwork/utils/imageutil",
				"github.com/teamwork/utils/ioutilx",
				"github.com/teamwork/utils/jsonutil",
				"github.com/teamwork/utils/maputil",
				"github.com/teamwork/utils/mathutil",
				"github.com/teamwork/utils/netutil",
				"github.com/teamwork/utils/raceutil",
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

func TestResolveImport(t *testing.T) {
	cases := []struct {
		inFile, inPkg, want, wantErr string
	}{
		// Twice to test it works from cache
		{"package main\nimport \"net/http\"\n", "http", "net/http", ""},
		{"package main\nimport \"os\"\n", "os", "os", ""},
		{"package main\nimport xxx \"net/http\"\n", "xxx", "net/http", ""},
		{"package main\nimport \"net/http\"\n", "httpx", "", ""},

		// Make sure it works from vendor
		{"package main\n import \"github.com/teamwork/test\"\n", "test", "github.com/teamwork/test", ""},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			f, clean := test.TempFile(t, tc.inFile)
			defer clean()

			out, err := ResolveImport(f, tc.inPkg)
			if !test.ErrorContains(err, tc.wantErr) {
				t.Fatalf("wrong err: %v", err)
			}
			if out != tc.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}

	t.Run("cache", func(t *testing.T) {
		f, clean := test.TempFile(t, "package main\nimport \"net/http\"\n")
		defer clean()

		importsCache = make(map[string]map[string]string)
		out, err := ResolveImport(f, "http")
		if err != nil {
			t.Fatal(err)
		}
		if out != "net/http" {
			t.Fatalf("out wrong: %v", out)
		}

		// Second time
		out, err = ResolveImport(f, "http")
		if err != nil {
			t.Fatal(err)
		}
		if out != "net/http" {
			t.Fatalf("out wrong: %v", out)
		}

		if len(importsCache) != 1 {
			t.Error(importsCache)
		}
	})
}

func TestTagName(t *testing.T) {
	cases := []struct {
		in, inName, want string
	}{
		{`json:"w00t"`, "json", "w00t"},
		{`yaml:"w00t"`, "json", "Original"},
		{`json:"w00t" yaml:"xxx""`, "yaml", "xxx"},
		{`JSON:"w00t"`, "json", "Original"},
		{`JSON: "w00t"`, "json", "Original"},
		{`json:"w00t,omitempty"`, "json", "w00t"},
		{`json:"w00t,"`, "json", "w00t"},
		{`json:"-"`, "json", "-"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc.in), func(t *testing.T) {
			f := &ast.Field{
				Names: []*ast.Ident{&ast.Ident{Name: "Original"}},
				Tag:   &ast.BasicLit{Value: fmt.Sprintf("`%v`", tc.in)}}

			out := TagName(f, tc.inName)
			if out != tc.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}

	t.Run("nil", func(t *testing.T) {
		f := &ast.Field{
			Names: []*ast.Ident{&ast.Ident{Name: "Original"}},
		}

		out := TagName(f, "json")
		if out != "Original" {
			t.Errorf("\nout:  %#v\nwant: %#v\n", out, "Original")
		}
	})

	t.Run("nil", func(t *testing.T) {
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("didn't panic")
			}
			if !strings.HasPrefix(r.(string), "cannot use TagName on struct with more than one name: ") {
				t.Errorf("wrong message: %#v", r)
			}
		}()

		f := &ast.Field{
			Names: []*ast.Ident{&ast.Ident{Name: "Original"},
				&ast.Ident{Name: "Second"}}}
		_ = TagName(f, "json")
	})

	t.Run("embed", func(t *testing.T) {
		cases := []struct {
			name string
			in   *ast.Field
			want string
		}{
			{
				"notag",
				&ast.Field{
					Tag:  &ast.BasicLit{Value: "`b:\"Bar\"`"},
					Type: &ast.Ident{Name: "Foo"},
				},
				"Foo",
			},
			{
				"ident",
				&ast.Field{Type: &ast.Ident{Name: "Foo"}},
				"Foo",
			},
			{
				"pointer",
				&ast.Field{Type: &ast.StarExpr{X: &ast.Ident{Name: "Foo"}}},
				"Foo",
			},
			{
				"pkg",
				&ast.Field{Type: &ast.SelectorExpr{Sel: &ast.Ident{Name: "Foo"}}},
				"Foo",
			},
			{
				"pkg-pointer",
				&ast.Field{
					Type: &ast.StarExpr{
						X: &ast.SelectorExpr{Sel: &ast.Ident{Name: "Foo"}},
					},
				},
				"Foo",
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				out := TagName(tc.in, "a")
				if out != tc.want {
					t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
				}
			})
		}

	})
}
