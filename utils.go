package rr

import (
	"context"
	"os/exec"
	"path/filepath"
	"strings"
)

func isHidden(path string) bool {
	return len(path) > 1 && strings.HasPrefix(filepath.Base(path), ".") || len(path) > 1 && strings.HasPrefix(filepath.Dir(path), ".")
}

func StringInSlice(s string, ss []string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

func InExcludeDir(s string, exclude []string) bool {
	s = filepath.Clean(filepath.Dir(s))
	return StringInSlice(s, exclude)
}

func NewCommand(ctx context.Context, s string) *exec.Cmd {
	ss := strings.Fields(s)
	if len(s) == 1 {
		return exec.CommandContext(ctx, s)
	}
	return exec.CommandContext(ctx, ss[0], ss[1:]...)
}

func NoneFunc() {

}
