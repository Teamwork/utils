// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package ioutilx

import (
	"os"
	"syscall"

	"github.com/pkg/errors"
)

func setOwner(srcStat os.FileInfo, dst string) error {
	statT, ok := srcStat.Sys().(*syscall.Stat_t)
	if !ok {
		return errors.New("could not get file owner: type assertion to syscall.Stat_t failed")
	}
	err := os.Chown(dst, int(statT.Uid), int(statT.Gid))
	if err != nil {
		return errors.Wrap(err, "could not chown")
	}

	return nil
}
