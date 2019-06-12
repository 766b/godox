package main

import (
	"bytes"
	"flag"
	"path/filepath"
	"testing"
)

func TestParse(t *testing.T) {
	flag.Parse()
	tests := []struct {
		path   string
		result string
	}{
		{
			path: "./fixtures/01",
			result: func() (r string) {
				r += filepath.Join("fixtures", "01", "example1.go:13:2:TODO(fix): something (Line 13)\n")
				r += filepath.Join("fixtures", "01", "example1.go:20:2:todo compare apples to oranges (Line 20)\n")
				r += filepath.Join("fixtures", "01", "example1.go:24:1:TODO: Multiline C1 (Line 24)\n")
				r += filepath.Join("fixtures", "01", "example1.go:25:1:TODO: Multiline C2 (Line 25)\n")
				r += filepath.Join("fixtures", "01", "example1.go:26:1:FIX: Your attitude (Line 26)\n")
				r += filepath.Join("fixtures", "01", "example2.go:4:12:TODO: Add JSON tag (Line 4)\n")
				r += filepath.Join("fixtures", "01", "example2.go:5:2:toDO add more fields (Line 5)\n")
				r += filepath.Join("fixtures", "01", "example2.go:11:1:TODO: multiline todo 1 (Line 11)\n")
				r += filepath.Join("fixtures", "01", "example2.go:15:1:TOdo multiline todo 2 (Line 15)\n")
				return
			}(),
		},
		{
			path: "./fixtures/02",
			result: func() (r string) {
				r += filepath.Join("fixtures", "02", "example3.go:3:1:TODO: remove foo (Line 3)\n")
				r += filepath.Join("fixtures", "02", "example3.go:7:14:TODO: Rename field (Line 7)\n")
				r += filepath.Join("fixtures", "02", "example3.go:10:1:TODO: get cat food (Line 10)\n")
				r += filepath.Join("fixtures", "02", "example3.go:15:3:todo  : todo comment (Line 15)\n")
				return
			}(),
		},
		{
			path: "./fixtures/03",
			result: func() (r string) {
				r += filepath.Join("fixtures", "03", "main.go:1:1:TODO: Add package documentation\n")
				r += filepath.Join("fixtures", "03", "main.go:2:1:TODO: Write an actual application\n")
				r += filepath.Join("fixtures", "03", "main.go:8:2:FIX: Spelling\n")
				r += filepath.Join("fixtures", "03", "main.go:13:1:TODO: Multi line 1\n")
				r += filepath.Join("fixtures", "03", "main.go:14:1:TODO: Multi line 2\n")
				r += filepath.Join("fixtures", "03", "main.go:15:1:FIX: Mutli line 3\n")
				return
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			output := bytes.NewBuffer(nil)
			new(tt.path, output).parse()

			if tt.result != output.String() {
				t.Errorf("not equal\nexpected: %s\nactual: %s", tt.result, output.String())
			}
		})
	}
}
