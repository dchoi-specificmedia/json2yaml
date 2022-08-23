package json2yaml_test

import (
	"strings"
	"testing"

	"github.com/itchyny/json2yaml"
)

func TestConvert(t *testing.T) {
	testCases := []struct {
		name string
		src  string
		want string
		err  string
	}{
		{
			name: "null",
			src:  "null",
			want: `null
`,
		},
		{
			name: "boolean",
			src:  "false true",
			want: `false
---
true
`,
		},
		{
			name: "number",
			src:  "0 128 -320 3.14 -6.63e-34",
			want: `0
---
128
---
-320
---
3.14
---
-6.63e-34
`,
		},
		{
			name: "string",
			src:  `"" "foo" "128" "hello, world" "\b\f\n\r\n\t"`,
			want: `""
---
"foo"
---
"128"
---
"hello, world"
---
"\u0008\u000c\n\r\n\t"
`,
		},
		{
			name: "empty object",
			src:  "{}",
			want: `{}
`,
		},
		{
			name: "simple object",
			src:  `{"foo": 128, "bar": null, "baz": false}`,
			want: `"foo": 128
"bar": null
"baz": false
`,
		},
		{
			name: "nested object",
			src: `{
				"foo": {"bar": {"baz": 128, "bar": null}, "baz": 0},
				"bar": {"foo": {}, "bar": {"bar": {}}, "baz": {}},
				"baz": {}
			}`,
			want: `"foo":
  "bar":
    "baz": 128
    "bar": null
  "baz": 0
"bar":
  "foo": {}
  "bar":
    "bar": {}
  "baz": {}
"baz": {}
`,
		},
		{
			name: "multiple objects",
			src:  `{}{"foo":128}{}`,
			want: `{}
---
"foo": 128
---
{}
`,
		},
		{
			name: "unclosed object",
			src:  "{",
			err:  "unexpected EOF",
		},
		{
			name: "empty array",
			src:  "[]",
			want: `[]
`,
		},
		{
			name: "simple array",
			src:  `[null,false,true,-128,12345678901234567890,"foo bar baz"]`,
			want: `- null
- false
- true
- -128
- 12345678901234567890
- "foo bar baz"
`,
		},
		{
			name: "nested array",
			src:  "[0,[1],[2,3],[4,[5,[6,[],7],[]],[8]],[],9]",
			want: `- 0
- - 1
- - 2
  - 3
- - 4
  - - 5
    - - 6
      - []
      - 7
    - []
  - - 8
- []
- 9
`,
		},
		{
			name: "nested object and array",
			src:  `{"foo":[0,{"bar":[],"foo":{}},[{"foo":[{"foo":[]}]}],[[[{}]]]],"bar":[{}]}`,
			want: `"foo":
  - 0
  - "bar": []
    "foo": {}
  - - "foo":
        - "foo": []
  - - - - {}
"bar":
  - {}
`,
		},
		{
			name: "multiple arrays",
			src:  `[][{"foo":128}][]`,
			want: `[]
---
- "foo": 128
---
[]
`,
		},
		{
			name: "deeply nested object",
			src:  `{"x":{"x":{"x":{"x":{"x":{"x":{"x":{"x":{"x":{"x":{"x":{"x":{"x":{"x":{"x":{"x":{"x":{"x":{"x":{"x":{}}}}}}}}}}}}}}}}}}}}}`,
			want: `"x":
  "x":
    "x":
      "x":
        "x":
          "x":
            "x":
              "x":
                "x":
                  "x":
                    "x":
                      "x":
                        "x":
                          "x":
                            "x":
                              "x":
                                "x":
                                  "x":
                                    "x":
                                      "x": {}
`,
		},
		{
			name: "unclosed array",
			src:  "[",
			err:  "unexpected EOF",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var sb strings.Builder
			err := json2yaml.Convert(&sb, strings.NewReader(tc.src))
			if tc.err == "" {
				if err != nil {
					t.Fatalf("should not raise an error but got: %s", err)
				}
				if got, want := diff(sb.String(), tc.want); got != want {
					t.Fatalf("should write\n  %q\nbut got\n  %q\nwhen source is\n  %q", want, got, tc.src)
				}
			} else {
				if err == nil {
					t.Fatalf("should raise an error %s but got no error", tc.err)
				}
				if !strings.Contains(err.Error(), tc.err) {
					t.Fatalf("should raise an error %s but got error %s", tc.err, err)
				}
			}
		})
	}
}

func diff(xs, ys string) (string, string) {
	if xs == ys {
		return "", ""
	}
	for {
		i := strings.IndexByte(xs, '\n')
		j := strings.IndexByte(ys, '\n')
		if i < 0 || j < 0 || xs[:i] != ys[:j] {
			break
		}
		xs, ys = xs[i+1:], ys[j+1:]
	}
	for {
		i := strings.LastIndexByte(xs, '\n')
		j := strings.LastIndexByte(ys, '\n')
		if i < 0 || j < 0 || xs[i:] != ys[j:] {
			break
		}
		xs, ys = xs[:i], ys[:j]
	}
	return xs, ys
}
