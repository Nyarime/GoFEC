//go:build arm64

package gf256

import "golang.org/x/sys/cpu"

var hasSVE = cpu.ARM64.HasSVE
