[![Build Status](https://travis-ci.org/Teamwork/utils.svg?branch=master)](https://travis-ci.org/Teamwork/utils)
[![codecov](https://codecov.io/gh/Teamwork/utils/branch/master/graph/badge.svg?token=n0k8YjbQOL)](https://codecov.io/gh/Teamwork/utils)
[![GoDoc](https://godoc.org/github.com/Teamwork/utils?status.svg)](https://godoc.org/github.com/Teamwork/utils)

`utils` is a collection of small – and sometimes not so small – extensions to
Go's standard library. There are no external dependencies.

The naming scheme is `[type]util` or `[pkgname]util`. If there already is a
`*util` packge in stdlib it's named `utilx` (e.g. `ioutilx`).

Other useful packages:

- [`github.com/teamwork/toutf8`](https://github.com/Teamwork/toutf8) – Convert
  strings to UTF-8.
- [`github.com/jinzhu/now`](https://github.com/jinzhu/now) – Time toolkit.
