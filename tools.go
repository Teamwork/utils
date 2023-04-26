//go:build tools
// +build tools

// Package tools is used for declaring tool dependencies, see:
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
package tools

import (
	_ "github.com/mattn/goveralls"
)
