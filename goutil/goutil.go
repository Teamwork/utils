// Package goutil provides functions to work with Go source files.
package goutil // import "github.com/teamwork/utils/goutil"

import (
	"errors"
	"go/build"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Expand a list of package and/or directory names to Go package names.
//
//  - "./example" is expanded to "full/package/path/example".
//  - "/absolute/src/package/path" is abbreviated to "package/path".
//  - "full/package" is kept-as is.
//  - "package/path/..." will include "package/path" and all subpackages.
//
// The packages will be sorted with duplicate packages removed. The /vendor/
// directory is automatically ignored.
func Expand(paths []string, mode build.ImportMode) ([]*build.Package, error) {
	var out []*build.Package
	seen := make(map[string]struct{})
	for _, p := range paths {
		if strings.HasSuffix(p, "/...") {
			subPkgs, err := ResolveWildcard(p, mode)
			if err != nil {
				return nil, err
			}
			for _, sub := range subPkgs {
				if _, ok := seen[sub.ImportPath]; !ok {
					out = append(out, sub)
					seen[sub.ImportPath] = struct{}{}
				}
			}
			continue
		}

		pkg, err := ResolvePackage(p, mode)
		if err != nil {
			return nil, err
		}
		if _, ok := seen[pkg.ImportPath]; !ok {
			out = append(out, pkg)
			seen[pkg.ImportPath] = struct{}{}
		}
	}

	sort.Slice(out, func(i, j int) bool { return out[i].ImportPath < out[j].ImportPath })

	return out, nil
}

// ResolvePackage resolves a package path, which can either be a local directory
// relative to the current dir (e.g. "./example"), a full path (e.g.
// ~/go/src/example"), or a package path (e.g. "example").
func ResolvePackage(path string, mode build.ImportMode) (pkg *build.Package, err error) {
	if len(path) == 0 {
		// TODO: maybe resolve like '.'? Dunno what makes more sense.
		return nil, errors.New("cannot resolve empty string")
	}

	switch path[0] {
	case '/':
		pkg, err = build.ImportDir(path, mode)
	case '.':
		path, err = filepath.Abs(path)
		if err != nil {
			return nil, err
		}
		pkg, err = build.ImportDir(path, mode)
	default:
		pkg, err = build.Import(path, ".", mode)
	}
	if err != nil {
		return nil, err
	}

	return pkg, err
}

// ResolveWildcard finds all subpackages in the "example/..." format. The
// "/vendor/" directory will be ignored.
func ResolveWildcard(path string, mode build.ImportMode) ([]*build.Package, error) {
	root, err := ResolvePackage(path[:len(path)-4], mode)
	if err != nil {
		return nil, err
	}

	// Gather a list of directories with *.go files.
	goDirs := make(map[string]struct{})
	err = filepath.Walk(root.Dir, func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".go") || info.IsDir() || strings.Contains(path, "/vendor/") {
			return nil
		}

		goDirs[filepath.Dir(path)] = struct{}{}
		return nil
	})
	if err != nil {
		return nil, err
	}

	var out []*build.Package
	for d := range goDirs {
		pkg, err := ResolvePackage(d, mode)
		if err != nil {
			return nil, err
		}
		out = append(out, pkg)
	}

	return out, nil
}
