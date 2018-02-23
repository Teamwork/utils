package ioutilx

// Note that these functions may not be portable to all systems. Specifically,
// none of this is tested on Windows.
//
// This code is based on: https://github.com/termie/go-shutil

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/teamwork/utils/sliceutil"
)

// ErrSameFile is used when the source and destination file are the same file.
type ErrSameFile struct {
	Src string
	Dst string
}

func (e ErrSameFile) Error() string {
	return fmt.Sprintf("%v and %v are the same file", e.Src, e.Dst)
}

// ErrExists is used when the destination already exists.
type ErrExists struct {
	Dst string
}

func (e ErrExists) Error() string {
	return fmt.Sprintf("%s already exists", e.Dst)
}

// ErrNotDir is used when attempting to copy a directory tree that is
// not a directory.
type ErrNotDir struct {
	Src string
}

func (e ErrNotDir) Error() string {
	return fmt.Sprintf("%s is not a directory", e.Src)
}

// ErrSpecialFile is used when the source or destination file is a special
// file, and not something we can operate on.
type ErrSpecialFile struct {
	FileInfo os.FileInfo
}

func (e ErrSpecialFile) Error() string {
	var mode string
	switch {
	case (e.FileInfo.Mode() & os.ModeDevice) == os.ModeDevice:
		mode = "device file"
	case (e.FileInfo.Mode() & os.ModeNamedPipe) == os.ModeNamedPipe:
		mode = "named pipe"
	case (e.FileInfo.Mode() & os.ModeSocket) == os.ModeSocket:
		mode = "domain socket"
	case (e.FileInfo.Mode() & os.ModeCharDevice) == os.ModeCharDevice:
		mode = "character device"
	default:
		panic("this should never happen")
	}

	return fmt.Sprintf("%v is not a regular file but a %v",
		e.FileInfo.Name(), mode)
}

// IsSameFile reports if two files refer to the same file object. It will return
// a ErrSameFile if it is.
//
// It is not an error if one of the two files doesn't exist.
func IsSameFile(src string, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return errors.WithStack(err)
	}

	dstInfo, err := os.Stat(dst)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return errors.WithStack(err)
	}

	if os.SameFile(srcInfo, dstInfo) {
		return &ErrSameFile{Src: src, Dst: dst}
	}
	return nil
}

// IsSpecialFile reports if this file is a special file such as a named pipe,
// device file, or socket. If so it will return a ErrSpecialFile.
func IsSpecialFile(fi os.FileInfo) error {
	if (fi.Mode()&os.ModeDevice) == os.ModeDevice ||
		(fi.Mode()&os.ModeNamedPipe) == os.ModeNamedPipe ||
		(fi.Mode()&os.ModeSocket) == os.ModeSocket ||
		(fi.Mode()&os.ModeCharDevice) == os.ModeCharDevice {

		return &ErrSpecialFile{FileInfo: fi}
	}

	return nil
}

// IsSymlink reports if this file is a symbolic link.
func IsSymlink(fi os.FileInfo) bool {
	return (fi.Mode() & os.ModeSymlink) == os.ModeSymlink
}

// CopyData copies the file data from the file in src to the path in dst.
//
// Note that this only copies data; permissions and other special file bits may
// get lost.
func CopyData(src, dst string) error {
	src, srcStat, err := copyCheck(src, dst)
	if err != nil {
		return errors.WithStack(err)
	}

	if _, exists := os.Stat(dst); exists == nil {
		return &ErrExists{dst}
	}

	// Do the actual copy.
	fsrc, err := os.Open(src)
	if err != nil {
		return errors.WithStack(err)
	}
	defer fsrc.Close() // nolint: errcheck

	fdst, err := os.Create(dst)
	if err != nil {
		return errors.Wrap(err, "create failed")
	}
	defer fdst.Close() // nolint: errcheck

	size, err := io.Copy(fdst, fsrc)
	if err != nil {
		return errors.Wrap(err, "copy failed")
	}

	if size != srcStat.Size() {
		return fmt.Errorf("%s: %d/%d copied", src, size, srcStat.Size())
	}

	return errors.Wrap(fdst.Close(), "close failed")
}

// Some sanity checks for CopyData() and CopyMode().
func copyCheck(src, dst string) (string, os.FileInfo, error) {
	if err := IsSameFile(src, dst); err != nil {
		return "", nil, errors.WithStack(err)
	}

	// Make sure src exists and neither are special files.
	srcStat, err := os.Lstat(src)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}
	if err := IsSpecialFile(srcStat); err != nil {
		return "", nil, errors.WithStack(err)
	}

	// Follow symlinks for the source file.
	if IsSymlink(srcStat) {
		dir := filepath.Dir(src)
		src, err = os.Readlink(src)
		if err != nil {
			return "", nil, errors.WithStack(err)
		}
		src, err = filepath.Abs(filepath.Join(dir, src))
		if err != nil {
			return "", nil, errors.WithStack(err)
		}

		srcStat, err = os.Stat(src)
		if err != nil {
			return "", nil, errors.WithStack(err)
		}
	}

	return src, srcStat, nil
}

// Modes to copy with CopyMode() and Copy().
type Modes struct {
	Permissions bool // chmod
	Owner       bool // chown
	Mtime       bool // mtime
}

// CopyMode copies the given modes from src to dst.
func CopyMode(src, dst string, modes Modes) error {
	_, srcStat, err := copyCheck(src, dst)
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = os.Stat(dst)
	if err != nil {
		return errors.WithStack(err)
	}

	if modes.Permissions {
		err := os.Chmod(dst, srcStat.Mode())
		if err != nil {
			return errors.Wrap(err, "could not chmod")
		}
	}

	if modes.Owner {
		statT, ok := srcStat.Sys().(*syscall.Stat_t)
		if !ok {
			return errors.New("could not get file owner: type assertion to syscall.Stat_t failed")
		}
		err := os.Chown(dst, int(statT.Uid), int(statT.Gid))
		if err != nil {
			return errors.Wrap(err, "could not chown")
		}
	}

	if modes.Mtime {
		err := os.Chtimes(dst, time.Now(), srcStat.ModTime())
		if err != nil {
			return errors.Wrap(err, "could not chown")
		}
	}

	return nil
}

// Copy data and the given mode bits; this is the same as a CopyData() followed
// by a CopyMode().
//
// The destination may be a directory, in which case the file name of "src" is
// used in that directory (as with "cp").
//
// TODO: we can make this a bit more efficient by not duplicating all the checks
// in CopyData() and CopyMode().
func Copy(src, dst string, modes Modes) error {
	dstInfo, err := os.Stat(dst)
	if err == nil && dstInfo.Mode().IsDir() {
		dst = filepath.Join(dst, filepath.Base(src))
	}

	if err != nil && !os.IsNotExist(err) {
		return errors.WithStack(err)
	}

	if err = CopyData(src, dst); err != nil {
		return errors.WithStack(err)
	}
	return errors.WithStack(CopyMode(src, dst, modes))
}

// CopyTreeOptions are flags for the CopyTree function.
type CopyTreeOptions struct {
	Symlinks               bool
	IgnoreDanglingSymlinks bool
	CopyFunction           func(string, string, Modes) error
	Ignore                 func(string, []os.FileInfo) []string
}

// DefaultCopyTreeOptions is used when the options to CopyTree() is nil.
var DefaultCopyTreeOptions = &CopyTreeOptions{
	Symlinks:               false,
	Ignore:                 nil,
	CopyFunction:           Copy,
	IgnoreDanglingSymlinks: false,
}

// CopyTree recursively copies a directory tree.
//
// The destination directory must not already exist.
//
// If the optional Symlinks flag is true, symbolic links in the source tree
// result in symbolic links in the destination tree; if it is false, the
// contents of the files pointed to by symbolic links are copied. If the file
// pointed by the symlink doesn't exist, an error will be returned.
//
// You can set the optional IgnoreDanglingSymlinks flag to true if you want to
// silence this error. Notice that this has no effect on platforms that don't
// support os.Symlink.
//
// The optional ignore argument is a callable. If given, it is called with the
// `src` parameter, which is the directory being visited by CopyTree(), and
// `names` which is the list of `src` contents, as returned by ioutil.ReadDir():
//
//   callable(src, entries) -> ignoredNames
//
// Since CopyTree() is called recursively, the callable will be called once for
// each directory that is copied. It returns a list of names relative to the
// `src` directory that should not be copied.
//
// The optional copyFunction argument is a callable that will be used to copy
// each file. It will be called with the source path and the destination path as
// arguments. By default, Copy() is used, but any function that supports the
// same signature (like Copy2() when it exists) can be used.
func CopyTree(src, dst string, options *CopyTreeOptions) error {
	if options == nil {
		options = DefaultCopyTreeOptions
	}

	// Sanity checks.
	srcFileInfo, err := os.Stat(src)
	if err != nil {
		return errors.WithStack(err)
	}
	if !srcFileInfo.IsDir() {
		return &ErrNotDir{src}
	}
	_, err = os.Open(dst)
	if !os.IsNotExist(err) {
		return &ErrExists{dst}
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return errors.Wrapf(err, "could not read %v", src)
	}

	// Create dst.
	if err = os.MkdirAll(dst, srcFileInfo.Mode()); err != nil {
		return errors.Wrapf(err, "could not create %v", dst)
	}

	ignoredNames := []string{}
	if options.Ignore != nil {
		ignoredNames = options.Ignore(src, entries)
	}

	for _, entry := range entries {
		if sliceutil.InStringSlice(ignoredNames, entry.Name()) {
			continue
		}

		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		entryFileInfo, err := os.Lstat(srcPath)
		if err != nil {
			return errors.WithStack(err)
		}

		switch {

		// Symlinks
		case IsSymlink(entryFileInfo):
			linkTo, err := os.Readlink(srcPath)
			if err != nil {
				return errors.WithStack(err)
			}
			dir := filepath.Dir(srcPath)
			linkTo, err = filepath.Abs(filepath.Join(dir, linkTo))
			if err != nil {
				return errors.WithStack(err)
			}

			if options.Symlinks {
				err = os.Symlink(linkTo, dstPath)
				if err != nil {
					return errors.WithStack(err)
				}
				// CopyStat(srcPath, dstPath, false)
			} else {
				// ignore dangling symlink if flag is on
				linkToStat, err := os.Stat(linkTo)
				if os.IsNotExist(err) && options.IgnoreDanglingSymlinks {
					continue
				}

				if linkToStat.IsDir() {
					err = CopyTree(srcPath, dstPath, options)
				} else {
					err = options.CopyFunction(srcPath, dstPath, Modes{})
				}
				if err != nil {
					return errors.WithStack(err)
				}
			}

		// Dir
		case entryFileInfo.IsDir():
			err = CopyTree(srcPath, dstPath, options)
			if err != nil {
				return errors.WithStack(err)
			}

		// Anything else.
		default:
			err = options.CopyFunction(srcPath, dstPath, Modes{})
			if err != nil {
				return errors.WithStack(err)
			}
		}
	}
	return nil
}
