package pretty

import "github.com/fatih/color"

type Theme struct {
	Nil     *color.Color
	Float   *color.Color
	Integer *color.Color
	String  *color.Color
	Bool    *color.Color
}
