package ioutilx

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/teamwork/test"
)

func TestMain(m *testing.M) {
	fifo := "test/fifo"
	err := syscall.Mkfifo(fifo, 644)
	if err != nil {
		panic(err)
	}

	e := m.Run()
	if err := os.Remove(fifo); err != nil {
		panic(err)
	}
	os.Exit(e)
}

func TestIsSymLink(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"test/file1", false},
		{"test/dir1", false},
		{"test/link1", true},
		{"test/link2", true},
		{"test/link3", true},
	}

	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			st, err := os.Lstat(tc.in)
			if err != nil {
				t.Fatal(err)
			}
			out := IsSymlink(st)
			if out != tc.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}

func TestIsSameFile(t *testing.T) {
	cases := []struct {
		src, dst string
		want     string
	}{
		{"test/file1", "test/link1", "are the same file"},
		{"test/file1", "test/link2", "are the same file"},
		{"test/file1", "test/dir1", ""},
		{"test/file1", "nonexistent", ""},
		{"nonexistent", "test/file1", ""},
		{"nonexistent", "nonexistent", ""},
	}

	for _, tc := range cases {
		t.Run(tc.src+":"+tc.dst, func(t *testing.T) {
			out := IsSameFile(tc.src, tc.dst)
			if !test.ErrorContains(out, tc.want) {
				t.Fatalf("\nwant: %s\nout:  %+v", tc.want, out)
			}
		})
	}
}

func TestIsSpecialFile(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"test/file1", ""},
		{"test/dir1", ""},
		{"test/link1", ""},
		{"test/fifo", "named pipe"},
		{"/dev/null", "device file"},
	}

	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {

			st, err := os.Lstat(tc.in)
			if err != nil {
				t.Fatal(err)
			}

			out := IsSpecialFile(st)
			if !test.ErrorContains(out, tc.want) {
				t.Fatalf("\nwant: %s\nout:  %+v", tc.want, out)
			}
		})
	}
}

func TestCopyData(t *testing.T) {
	cases := []struct {
		src, dst, want string
	}{
		{"test/file1", "test/file1", "same file"},
		{"nonexistent", "test/copydst", "no such file"},
		{"test/file1", "test/file2", "already exists"},
		{"test/fifo", "test/newfile", "named pipe"},
		{"test/link1/asd", "test/dst1", "not a directory"},
		{"test/file1", "/cantwritehere", "permission denied"},
		{"test/file1", "test/dst1", ""},
		{"test/link1", "test/dst1", ""},
	}

	for _, tc := range cases {
		t.Run(tc.src+":"+tc.dst, func(t *testing.T) {
			if tc.want == "" {
				defer clean(t, tc.dst)
			}

			out := CopyData(tc.src, tc.dst)
			if !test.ErrorContains(out, tc.want) {
				t.Fatalf("\nwant: %s\nout:  %+v", tc.want, out)
			}

			if tc.want == "" {
				filesMatch(t, tc.src, tc.dst)
			}
		})
	}
}

func TestCopyMode(t *testing.T) {
	cases := []struct {
		src, dst string
		mode     Modes
		want     string
	}{
		{"test/file1", "test/file1", Modes{}, "same file"},
		{"nonexistent", "test/copydst", Modes{}, "no such file"},
		{"test/fifo", "test/newfile", Modes{}, "named pipe"},
		{"test/link1/asd", "test/dst1", Modes{}, "not a directory"},
		{"test/file1", "/cantwritehere", Modes{}, "no such file or directory"},

		{"test/exec", "test/dst1", Modes{Permissions: true, Owner: true, Mtime: true}, ""},
	}

	for _, tc := range cases {
		t.Run(tc.src+":"+tc.dst, func(t *testing.T) {
			if tc.want == "" {
				touch(t, tc.dst)
				defer clean(t, tc.dst)
			}

			out := CopyMode(tc.src, tc.dst, tc.mode)
			if !test.ErrorContains(out, tc.want) {
				t.Fatalf("\nwant: %s\nout:  %+v", tc.want, out)
			}

			if tc.want == "" {
				mode, err := os.Stat(tc.dst)
				if err != nil {
					t.Fatal(err)
				}

				if mode.Mode().String() != "-rwxr-xr-x" {
					t.Fatalf("wrong mode: %s", mode.Mode())
				}
			}
		})
	}
}

func TestCopy(t *testing.T) {
	cases := []struct {
		src, dst string
		mode     Modes
		want     string
	}{
		{"test/file1", "test/file1", Modes{}, "same file"},
		{"nonexistent", "test/copydst", Modes{}, "no such file"},
		{"test/fifo", "test/newfile", Modes{}, "named pipe"},
		{"test/link1/asd", "test/dst1", Modes{}, "not a directory"},
		{"test/file1", "/cantwritehere", Modes{}, "permission denied"},

		{"test/exec", "test/dst1", Modes{Permissions: true, Owner: true, Mtime: true}, ""},
		{"test/exec", "test/dir1", Modes{Permissions: true, Owner: true, Mtime: true}, ""},
		{"test/exec", "test/dir1/", Modes{Permissions: true, Owner: true, Mtime: true}, ""},
	}

	for _, tc := range cases {
		t.Run(tc.src+":"+tc.dst, func(t *testing.T) {
			c := tc.dst
			if strings.HasPrefix(c, "test/dir1") {
				c = filepath.Join(c, "exec")
			}

			if tc.want == "" {
				defer clean(t, c)
			}

			out := Copy(tc.src, tc.dst, tc.mode)
			if !test.ErrorContains(out, tc.want) {
				t.Fatalf("\nwant: %s\nout:  %+v", tc.want, out)
			}

			if tc.want == "" {
				mode, err := os.Stat(c)
				if err != nil {
					t.Fatal(err)
				}

				if mode.Mode().String() != "-rwxr-xr-x" {
					t.Fatalf("wrong mode: %s", mode.Mode())
				}
			}
		})
	}
}

func clean(t *testing.T, n string) {
	err := os.Remove(n)
	if err != nil {
		t.Fatalf("could not cleanup %v: %v", n, err)
	}
}

func filesMatch(t *testing.T, src, dst string) {
	srcContents, err := ioutil.ReadFile(src)
	if err != nil {
		t.Fatal(err)
	}

	dstContents, err := ioutil.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(srcContents, dstContents) {
		t.Errorf("%v and %v are not identical\nout:  %s\nwant: %s\n",
			src, dst, srcContents, dstContents)
	}
}

func touch(t *testing.T, n string) {
	fp, err := os.Create(n)
	if err != nil {
		t.Fatal(err)
	}
	err = fp.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCopyTree(t *testing.T) {
	t.Run("nonexistent", func(t *testing.T) {
		err := CopyTree("nonexistent", "test_copytree", nil)
		if !test.ErrorContains(err, "no such file or directory") {
			t.Error(err)
		}
	})
	t.Run("dst-exists", func(t *testing.T) {
		err := CopyTree("test", "test", nil)
		if !test.ErrorContains(err, "already exists") {
			t.Error(err)
		}
	})
	t.Run("dst-nodir", func(t *testing.T) {
		err := CopyTree("test/file1", "test", nil)
		if !test.ErrorContains(err, "not a directory") {
			t.Error(err)
		}
	})
	t.Run("permission", func(t *testing.T) {
		err := CopyTree("test", "/cant/write/here", nil)
		if !test.ErrorContains(err, "permission denied") {
			t.Error(err)
		}
	})

	defer func() {
		err := os.RemoveAll("test_copytree")
		if err != nil {
			t.Fatalf("could not clean: %v", err)
		}
	}()

	err := CopyTree("test", "test_copytree", &CopyTreeOptions{
		Symlinks: false,
		Ignore: func(path string, fi []os.FileInfo) []string {
			return []string{"fifo"}
		},
		CopyFunction:           Copy,
		IgnoreDanglingSymlinks: false,
	})
	if err != nil {
		t.Error(err)
		return
	}

	filesMatch(t, "test/file1", "test_copytree/file1")
}
