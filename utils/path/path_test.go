package path

import (
	"strings"
	"testing"
)

const (
	pathSep = "$"
)

var (
	modifiers = []string{modSlashNix, modSlashWin}
)

func TestIsAbs(t *testing.T) {
	fixtures := []struct {
		path  string
		isAbs bool
	}{
		{
			path:  "test",
			isAbs: false,
		},
		{
			path:  ".test",
			isAbs: false,
		},
		{
			path:  "$.test",
			isAbs: true,
		},
		{
			path:  "$$",
			isAbs: true,
		},
		{
			path:  "$$.test",
			isAbs: true,
		},
		{
			path:  ".$",
			isAbs: false,
		},
		{
			path:  ".$.test",
			isAbs: false,
		},
		{
			path:  "$.$.test",
			isAbs: true,
		},
		{
			path:  "..$",
			isAbs: false,
		},
		{
			path:  "..$.test",
			isAbs: false,
		},
		{
			path:  "..$$.test",
			isAbs: false,
		},
		{
			path:  "..$$.test",
			isAbs: false,
		},
	}

	for _, m := range modifiers {
		for i, f := range fixtures {
			p := strings.Replace(f.path, pathSep, m, -1)
			if IsAbs(p) != f.isAbs {
				t.Fatalf("#%d: %s failed the test", i, f.path)
			}
		}
	}
}

func TestIsRel(t *testing.T) {
	fixtures := []struct {
		path  string
		isRel bool
	}{
		{
			path:  "test",
			isRel: false,
		},
		{
			path:  ".test",
			isRel: false,
		},
		{
			path:  "$.test",
			isRel: false,
		},
		{
			path:  "$$",
			isRel: false,
		},
		{
			path:  "$$.test",
			isRel: false,
		},
		{
			path:  ".$",
			isRel: true,
		},
		{
			path:  ".$.test",
			isRel: true,
		},
		{
			path:  "$.$.test",
			isRel: false,
		},
		{
			path:  "..$",
			isRel: true,
		},
		{
			path:  "..$.test",
			isRel: true,
		},
		{
			path:  "..$$.test",
			isRel: true,
		},
		{
			path:  "..$$.test",
			isRel: true,
		},
	}
	for _, m := range modifiers {
		for i, f := range fixtures {
			p := strings.Replace(f.path, pathSep, m, -1)
			if IsRel(p) != f.isRel {
				t.Fatalf("#%d: %s failed the test", i, f.path)
			}
		}
	}
}

func TestTrimModPrefix(t *testing.T) {
	fixtures := []struct {
		path    string
		trimmed string
	}{
		{
			path:    "test",
			trimmed: "test",
		},
		{
			path:    ".test",
			trimmed: ".test",
		},
		{
			path:    "$.test",
			trimmed: ".test",
		},
		{
			path:    "$$",
			trimmed: "",
		},
		{
			path:    "$$.test",
			trimmed: ".test",
		},
		{
			path:    ".$",
			trimmed: "",
		},
		{
			path:    ".$.test",
			trimmed: ".test",
		},
		{
			path:    "$.$.test",
			trimmed: ".test",
		},
		{
			path:    "..$",
			trimmed: "",
		},
		{
			path:    "..$.test",
			trimmed: ".test",
		},
		{
			path:    "..$$.test",
			trimmed: ".test",
		},
		{
			path:    "..$$.test",
			trimmed: ".test",
		},
	}
	for _, m := range modifiers {
		for i, f := range fixtures {
			p := strings.Replace(f.path, pathSep, m, -1)
			got := TrimModPrefix(p)
			want := strings.Replace(f.trimmed, pathSep, m, -1)
			if got != want {
				t.Fatalf("#%d: %s failed the test, got: '%s', want: '%s'", i, f.path, got, want)
			}
		}
	}
}
