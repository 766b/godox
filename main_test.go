package main

import (
	"bytes"
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	flag.Parse()
	tests := []struct {
		path   string
		result string
	}{
		{
			path: "./test",
			result: `test\example1.go:13:2:TODO(fix): something (Line 13)
test\example1.go:20:2:todo compare apples to oranges (Line 20)
test\example1.go:24:1:TODO: Multiline C1 (Line 24)
test\example1.go:25:1:TODO: Multiline C2 (Line 25)
test\example1.go:26:1:FIX: Your attitude (Line 26)
test\example2.go:4:12:TODO: Add JSON tag (Line 4)
test\example2.go:5:2:toDO add more fields (Line 5)
test\example2.go:11:1:TODO: multiline todo 1 (Line 11)
test\example2.go:15:1:TOdo multiline todo 2 (Line 15)
`,
		},
		{
			path: "./test/addon",
			result: `test\addon\example3.go:3:1:TODO: remove foo (Line 3)
test\addon\example3.go:7:14:TODO: Rename field (Line 7)
test\addon\example3.go:10:1:TODO: get cat food (Line 10)
test\addon\example3.go:15:3:todo  : todo comment (Line 15)
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			output := bytes.NewBuffer(nil)
			New(tt.path, output).parse()
			assert.Equal(t, tt.result, output.String())
		})
	}
}
