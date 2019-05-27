package main

import (
	"flag"
	"testing"
)

func TestParse(t *testing.T) {
	flag.Parse()
	tests := []struct {
		path         string
		result       []string
		includeTests bool
	}{
		{
			path: "./fixtures/01",
			result: []string{
				`fixtures/01/example1.go:13:2:TODO(fix): something (Line 13)`,
				`fixtures/01/example1.go:20:2:todo compare apples to oranges (Line 20)`,
				`fixtures/01/example1.go:24:1:TODO: Multiline C1 (Line 24)`,
				`fixtures/01/example1.go:25:1:TODO: Multiline C2 (Line 25)`,
				`fixtures/01/example1.go:26:1:FIX: Your attitude (Line 26)`,
				`fixtures/01/example2.go:4:12:TODO: Add JSON tag (Line 4)`,
				`fixtures/01/example2.go:5:2:toDO add more fields (Line 5)`,
				`fixtures/01/example2.go:11:1:TODO: multiline todo 1 (Line 11)`,
				`fixtures/01/example2.go:15:1:TOdo multiline todo 2 (Line 15)`,
			},
		},
		{
			path: "./fixtures/02",
			result: []string{
				`fixtures/02/example3.go:3:1:TODO: remove foo (Line 3)`,
				`fixtures/02/example3.go:7:14:TODO: Rename field (Line 7)`,
				`fixtures/02/example3.go:10:1:TODO: get cat food (Line 10)`,
				`fixtures/02/example3.go:15:3:todo  : todo comment (Line 15)`,
				`fixtures/02/example3_test.go:8:2:TODO write test`,
			},
			includeTests: true,
		},
		{
			path: "./fixtures/03",
			result: []string{
				`fixtures/03/main.go:1:1:TODO: Add package documentation`,
				`fixtures/03/main.go:2:1:TODO: Write an actual application`,
				`fixtures/03/main.go:8:2:FIX: Spelling`,
				`fixtures/03/main.go:13:1:TODO: Multi line 1`,
				`fixtures/03/main.go:14:1:TODO: Multi line 2`,
				`fixtures/03/main.go:15:1:FIX: Mutli line 3`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			messages, err := godox(tt.path, tt.includeTests)

			if err != nil {
				t.Error("should not return error")
			}

			if len(messages) > len(tt.result) {
				t.Error("should return less messages")
			}

			if len(messages) < len(tt.result) {
				t.Error("should return more messages")
			}

			for i := range tt.result {
				if tt.result[i] != messages[i] {
					t.Errorf("not equal\nexpected: %s\nactual: %s", tt.result[i], messages[i])
				}
			}
		})
	}
}
