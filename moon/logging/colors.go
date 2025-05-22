package logging

import "fmt"

const (
	ColorBlack        = 30
	ColorRed          = 31
	ColorGreen        = 32
	ColorYellow       = 33
	ColorBlue         = 34
	ColorMagenta      = 35
	ColorCyan         = 36
	ColorLightGray    = 37
	ColorGray         = 90
	ColorLightRed     = 91
	ColorLightGreen   = 92
	ColorLightYellow  = 93
	ColorLightBlue    = 94
	ColorLightMagenta = 95
	ColorLightCyan    = 96
	ColorWhite        = 97
)

func Color(color int, str string) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", color, str)
}

func Black(str string) string        { return Color(ColorBlack, str) }
func Red(str string) string          { return Color(ColorRed, str) }
func Green(str string) string        { return Color(ColorGreen, str) }
func Yellow(str string) string       { return Color(ColorYellow, str) }
func Blue(str string) string         { return Color(ColorBlue, str) }
func Magenta(str string) string      { return Color(ColorMagenta, str) }
func Cyan(str string) string         { return Color(ColorCyan, str) }
func LightGray(str string) string    { return Color(ColorLightGray, str) }
func Gray(str string) string         { return Color(ColorGray, str) }
func LightRed(str string) string     { return Color(ColorLightRed, str) }
func LightGreen(str string) string   { return Color(ColorLightGreen, str) }
func LightYellow(str string) string  { return Color(ColorLightYellow, str) }
func LightBlue(str string) string    { return Color(ColorLightBlue, str) }
func LightMagenta(str string) string { return Color(ColorLightMagenta, str) }
func LightCyan(str string) string    { return Color(ColorLightCyan, str) }
func White(str string) string        { return Color(ColorWhite, str) }
