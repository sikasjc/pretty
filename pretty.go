package pretty

import (
	"bytes"
	"fmt"
	"github.com/rogpeppe/go-internal/fmtsort"
	"io"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	defaultIndent    = "  " // two space
	defaultNilString = "nil"
)

func NoColor() {
	DefaultPrinter.NoColor()
}

var (
	DefaultPrinter = &Printer{
		Theme:       DefaultTheme,
		Indent:      defaultIndent,
		NilString:   defaultNilString,
		SortMapKey:  ASC,
		Hexadecimal: true,
	}

	DefaultOut = os.Stdout
)

type Mod = int

const (
	_ Mod = iota
	ASC
	DESC
)

type HandleUnsupportedType func(reflect.Value) string

// Printer The context for printing
type Printer struct {
	*Theme
	// Indent indent string
	Indent string
	// NilString string for nil
	NilString string
	// CompactArray
	CompactArray bool
	// CompactMap
	CompactMap bool
	// MaxLevel Maximum nesting level
	MaxLevel int
	// SortMapKey if map key is string or number, sort and print
	SortMapKey Mod
	// Hexadecimal use hexadecimal format `byte`
	Hexadecimal bool
	// HandleUnsupportedType ...
	HandleUnsupportedType HandleUnsupportedType
}

// Print pretty print the input value (to stdout)
func Print(i interface{}) {
	PrintTo(DefaultOut, i)
}

// Println pretty print the input value (to stdout)
func Println(i interface{}) {
	PrintlnTo(DefaultOut, i)
}

// Format pretty print the input value (to a string)
func Format(i interface{}) string {
	var out bytes.Buffer
	PrintTo(&out, i)
	return out.String()
}

// PrintTo pretty print the input value (to specified writer)
func PrintTo(w io.Writer, i interface{}) {
	DefaultPrinter.Print(w, i)
}

// PrintlnTo pretty print the input value (to specified writer)
func PrintlnTo(w io.Writer, i interface{}) {
	DefaultPrinter.Println(w, i)
}

// Print pretty print the input value (no newline)
func (p *Printer) Print(w io.Writer, i interface{}) {
	p.PrintValue(w, reflect.ValueOf(i), 0)
}

// Println pretty print the input value (newline)
func (p *Printer) Println(w io.Writer, i interface{}) {
	p.PrintValue(w, reflect.ValueOf(i), 0)
	WriteString(w, "\n")
}

func (p *Printer) printMap(w io.Writer, val reflect.Value, level int) {
	l := val.Len()
	if l == 0 {
		WriteString(w, "{}")
		return
	}
	cur := strings.Repeat(p.Indent, level)
	next := strings.Repeat(p.Indent, level+1)
	nl, inner := "\n", "\n"
	if p.CompactMap {
		nl = ""
		inner = " "
	}
	WriteString(w, "{"+nl)
	keys := val.MapKeys()
	switch p.SortMapKey {
	case ASC:
		sortedMap := fmtsort.Sort(val)
		keys = sortedMap.Key
	case DESC:
		sortedMap := fmtsort.Sort(val)
		sort.Sort(sort.Reverse(sortedMap))
		keys = sortedMap.Key
	}
	for i, k := range keys {
		if !p.CompactMap {
			WriteString(w, next)
		}
		p.PrintValue(w, k, level)
		WriteString(w, ": ")
		v := val.MapIndex(k)
		if !IsPrimitive(v) && !IsEmpty(v) {
			WriteString(w, nl)
			WriteString(w, next+p.Indent)
			p.PrintValue(w, v, level+2)
		} else {
			p.PrintValue(w, v, level+1)
		}
		if i < l-1 {
			WriteString(w, ","+inner)
		} else {
			WriteString(w, nl)
		}
	}
	if !p.CompactMap {
		WriteString(w, cur)
	}
	WriteString(w, "}")
}

func (p *Printer) printArray(w io.Writer, val reflect.Value, level int) {
	l := val.Len()
	if l == 0 {
		WriteString(w, "[]")
		return
	}
	cur := strings.Repeat(p.Indent, level)
	next := strings.Repeat(p.Indent, level+1)
	start, end, inner := "\n", "\n", "\n"
	if p.CompactArray {
		start, end, inner = "", "", " "
	}
	WriteString(w, "["+start)
	if val.Index(0).Kind() == reflect.Uint8 && !p.CompactArray {
		HexDump(w, val.Interface().([]uint8), 16, next)
	} else {
		for i := 0; i < l; i++ {
			if !p.CompactArray {
				WriteString(w, next)
			}
			p.PrintValue(w, val.Index(i), level+1)
			if i < l-1 {
				WriteString(w, ","+inner)
			} else {
				WriteString(w, end)
			}
		}
	}
	if !p.CompactArray {
		WriteString(w, cur)
	}
	WriteString(w, "]")
}

func (p *Printer) printStruct(w io.Writer, val reflect.Value, level int) {
	if val.CanInterface() {
		cur := strings.Repeat(p.Indent, level)
		next := strings.Repeat(p.Indent, level+1)
		nl := "\n"

		i := val.Interface()
		switch v := i.(type) {
		case time.Time:
			p.writeTime(w, v)
			return
		case fmt.Stringer:
			WriteString(w, v.String())
			return
		}

		l := val.NumField()

		sOpen := "{"

		if l == 0 {
			WriteString(w, "{}")
		} else {
			WriteString(w, sOpen+nl)
			for i := 0; i < l; i++ {
				WriteString(w, next)
				WriteString(w, val.Type().Field(i).Name)
				WriteString(w, ": ")
				p.PrintValue(w, val.Field(i), level+1)
				if i < l-1 {
					WriteString(w, ","+nl)
				} else {
					WriteString(w, nl)
				}
			}
			WriteString(w, cur)
			WriteString(w, "}")
		}
	} else {
		WriteString(w, "protected")
	}
}

// PrintValue recursively print the input value
func (p *Printer) PrintValue(w io.Writer, val reflect.Value, level int) {
	if !val.IsValid() {
		p.writeNil(w)
		return
	}

	if p.MaxLevel > 0 && level >= p.MaxLevel {
		WriteString(w, val.String())
		return
	}

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p.writeInteger(w, val.Int())

	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		p.writeUInteger(w, val.Uint())

	case reflect.Uint8: // byte
		format := "%d"
		if p.Hexadecimal {
			format = "0x%x"
		}
		WriteString(w, fmt.Sprintf(format, val.Uint()))

	case reflect.Float32, reflect.Float64:
		p.writeFloat(w, val.Float())

	case reflect.String:
		p.writeString(w, val.String())

	case reflect.Bool:
		p.writeBool(w, val.Bool())

	case reflect.Map:
		p.printMap(w, val, level)

	case reflect.Array, reflect.Slice:
		p.printArray(w, val, level)

	case reflect.Interface, reflect.Ptr:
		p.PrintValue(w, val.Elem(), level)

	case reflect.Struct:
		p.printStruct(w, val, level)

	default:
		if p.HandleUnsupportedType != nil {
			output := p.HandleUnsupportedType(val)
			WriteString(w, output)
			return
		}
		WriteString(w, "unsupported:")
		WriteString(w, val.String())
	}
}

// GetTheme returns the theme used, or the default theme if not set
func (p *Printer) GetTheme() *Theme {
	theme := p.Theme
	if theme == nil {
		theme = DefaultTheme
	}
	return theme
}

func (p *Printer) writeNil(w io.Writer) {
	_, _ = p.GetTheme().Nil.Fprint(w, p.NilString)
}

func (p *Printer) writeInteger(w io.Writer, val int64) {
	_, _ = p.GetTheme().Integer.Fprint(w, strconv.FormatInt(val, 10))
}

func (p *Printer) writeUInteger(w io.Writer, val uint64) {
	_, _ = p.GetTheme().Integer.Fprint(w, strconv.FormatUint(val, 10))
}

func (p *Printer) writeFloat(w io.Writer, val float64) {
	_, _ = p.GetTheme().Float.Fprint(w, strconv.FormatFloat(val, 'f', -1, 64))
}

func (p *Printer) writeString(w io.Writer, val string) {
	_, _ = p.GetTheme().String.Fprint(w, strconv.Quote(val))
}

func (p *Printer) writeBool(w io.Writer, val bool) {
	_, _ = p.GetTheme().Bool.Fprint(w, strconv.FormatBool(val))
}

func (p *Printer) writeTime(w io.Writer, val time.Time) {
	theme := p.GetTheme()
	_, _ = theme.Time.Fprint(w, val.Format(theme.TimeLayout))
}

func IsPrimitive(val reflect.Value) bool {
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String,
		reflect.Bool:
		return true
	case reflect.Interface, reflect.Ptr:
		return IsPrimitive(val.Elem())
	}
	return false
}

func IsEmpty(val reflect.Value) bool {
	switch val.Kind() {
	case reflect.Map, reflect.Array, reflect.Slice:
		return val.Len() == 0
	case reflect.Interface, reflect.Ptr:
		return IsEmpty(val.Elem())
	}
	return false
}

func WriteString(w io.Writer, s string) {
	_, _ = io.WriteString(w, s)
}
