pretty
======
A Pretty-printer for Go data structures.

## Features
1. Support to print maps ordered by key using [fmtsort](https://github.com/rogpeppe/go-internal/fmtsort)
2. Support to print byte with hexadecimal
3. Expose a handler to handle Unsupported Type
4. Support to print compact array or slice

## Installation
```
go get github.com/sikasjc/pretty
```

## Example
```Go
package main

import "github.com/sikasjc/pretty"

func main() {
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

    pretty.Print(m)
}

/* output
{
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
}
 */
```
