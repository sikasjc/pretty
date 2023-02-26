package pretty

import (
	"bytes"
	"github.com/stretchr/testify/assert"
)
import "testing"

func TestMap(t *testing.T) {
	m := map[string]interface{}{
		"a":     123,
		"b":     "string",
		"z":     1.2,
		"c":     "xyz",
		"foo":   "bar",
		"empty": map[string]interface{}{},
		"map": map[int]interface{}{
			1: "one",
			2: 2,
			3: 3.0,
		},
	}
	t.Run("not compact", func(t *testing.T) {
		want := `{
  "a": 123,
  "b": "string",
  "c": "xyz",
  "empty": {},
  "foo": "bar",
  "map": 
    {
      1: "one",
      2: 2,
      3: 3
    },
  "z": 1.2
}`
		var out bytes.Buffer
		p := Printer{SortMapKey: ASC, Indent: defaultIndent}
		p.Print(&out, m)
		assert.Equal(t, want, out.String())
	})
	t.Run("compact", func(t *testing.T) {
		want := "{\"a\": 123, \"b\": \"string\", \"c\": \"xyz\", \"empty\": {}, \"foo\": \"bar\", \"map\": {1: \"one\", 2: 2, 3: 3}, \"z\": 1.2}"
		var out bytes.Buffer
		p := Printer{SortMapKey: ASC, CompactMap: true}
		p.Print(&out, m)
		assert.Equal(t, want, out.String())
	})
}

func TestStruct(t *testing.T) {
	type Struct struct {
		A string
		B int
		C float64
		D struct {
			e map[string]string
			f []int
		}
	}
	s := Struct{
		A: "a",
		B: 2,
		C: 3.0,
		D: struct {
			e map[string]string
			f []int
		}{
			e: map[string]string{"e": "e"},
			f: []int{1, 2, 3, 4}},
	}
	want := `{
  A: "a",
  B: 2,
  C: 3,
  D: {
    e: {"e": "e"},
    f: [1, 2, 3, 4]
  }
}`
	var out bytes.Buffer
	p := Printer{SortMapKey: ASC, Indent: defaultIndent, CompactMap: true, CompactArray: true}
	p.Print(&out, s)
	assert.Equal(t, want, out.String())
}

func TestSortMapKey(t *testing.T) {
	m := map[string]interface{}{
		"a":   123,
		"b":   "string",
		"z":   1.2,
		"c":   "xyz",
		"foo": "bar",
	}
	t.Run("asc", func(t *testing.T) {
		want := "{\"a\": 123, \"b\": \"string\", \"c\": \"xyz\", \"foo\": \"bar\", \"z\": 1.2}"
		var out bytes.Buffer
		p := Printer{SortMapKey: ASC, CompactMap: true}
		p.Print(&out, m)
		assert.Equal(t, want, out.String())
	})

	t.Run("desc", func(t *testing.T) {
		want := "{\"z\": 1.2, \"foo\": \"bar\", \"c\": \"xyz\", \"b\": \"string\", \"a\": 123}"
		var out bytes.Buffer
		p := Printer{SortMapKey: DESC, CompactMap: true}
		p.Print(&out, m)
		assert.Equal(t, want, out.String())
	})
}

func TestHexadecimal(t *testing.T) {
	t.Run("compact", func(t *testing.T) {
		arr := []interface{}{
			123,
			[]string{"b", "B"},
			[]byte{1, 2, 3},
		}
		want := "[123, [\"b\", \"B\"], [0x1, 0x2, 0x3]]"
		var out bytes.Buffer
		p := Printer{Hexadecimal: true, CompactArray: true}
		p.Print(&out, arr)
		assert.Equal(t, want, out.String())
	})
	t.Run("hex dump array", func(t *testing.T) {
		arr := []interface{}{
			123,
			[]string{"b", "B"},
			[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
		}
		want := `[
  123,
  [
    "b",
    "B"
  ],
  [
    0000 01  02  03  04  05  06  07  08  09  0a  0b  0c  0d  0e  0f  10    '................'
    0016 11                                                                '.'
  ]
]`
		var out bytes.Buffer
		p := Printer{Indent: defaultIndent, Hexadecimal: true}
		p.Print(&out, arr)
		assert.Equal(t, want, out.String())
	})
	t.Run("hex dump map", func(t *testing.T) {
		want := `[
  123,
  [
    "b",
    "B"
  ],
  [
    1,
    2,
    3,
    4,
    5
  ],
  [
    [
      1,
      2,
      3,
      4,
      5
    ]
  ],
  {
    "map": 
      [
        0000 01  02  03  04  05  06  07  08  09  0a  0b  0c  0d  0e  0f  10    '................'
        0016 11                                                                '.'
      ]
  },
  {
    "A": "B",
    "C": "D"
  }
]`
		arr := []interface{}{
			123,
			[]string{"b", "B"},
			[]int{
				1, 2, 3, 4, 5,
			},
			[][]int{
				{1, 2, 3, 4, 5},
			},
			map[string][]byte{
				"map": {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
			},
			map[string]string{
				"A": "B",
				"C": "D",
			},
		}
		var out bytes.Buffer
		p := Printer{Indent: defaultIndent, Hexadecimal: true, SortMapKey: ASC}
		p.Print(&out, arr)
		assert.Equal(t, want, out.String())
	})
}
