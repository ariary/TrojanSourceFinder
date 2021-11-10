package utils

import "fmt"

var (
	Info = Teal
	Warn = Yellow
	Evil = Red
	Good = Green
	Code = Cyan
)

var (
	Black         = Color("\033[1;30m%s\033[0m")
	Red           = Color("\033[1;31m%s\033[0m")
	Green         = Color("\033[1;32m%s\033[0m")
	Yellow        = Color("\033[1;33m%s\033[0m")
	Purple        = Color("\033[1;34m%s\033[0m")
	Magenta       = Color("\033[1;35m%s\033[0m")
	Teal          = Color("\033[1;36m%s\033[0m")
	White         = Color("\033[1;37m%s\033[0m")
	Cyan          = Color("\033[1;96m%s\033[0m")
	Underlined    = Color("\033[4m%s\033[24m")
	Bold          = Color("\033[1m%s\033[0m")
	Italic        = Color("\033[3m%s\033[0m")
	RedForeground = Color("\033[1;41m%s\033[0m")
)

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}
