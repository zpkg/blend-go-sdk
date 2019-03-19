package ansi

// Black applies a given color.
func Black(text string) string {
	return ApplyColor(ColorBlack, text)
}

// Red applies a given color.
func Red(text string) string {
	return ApplyColor(ColorRed, text)
}

// Green applies a given color.
func Green(text string) string {
	return ApplyColor(ColorGreen, text)
}

// Yellow applies a given color.
func Yellow(text string) string {
	return ApplyColor(ColorYellow, text)
}

// Blue applies a given color.
func Blue(text string) string {
	return ApplyColor(ColorBlue, text)
}

// Purple applies a given color.
func Purple(text string) string {
	return ApplyColor(ColorPurple, text)
}

// Cyan applies a given color.
func Cyan(text string) string {
	return ApplyColor(ColorCyan, text)
}

// White applies a given color.
func White(text string) string {
	return ApplyColor(ColorWhite, text)
}

// LightBlack applies a given color.
func LightBlack(text string) string {
	return ApplyColor(ColorLightBlack, text)
}

// LightRed applies a given color.
func LightRed(text string) string {
	return ApplyColor(ColorLightRed, text)
}

// LightGreen applies a given color.
func LightGreen(text string) string {
	return ApplyColor(ColorLightGreen, text)
}

// LightYellow applies a given color.
func LightYellow(text string) string {
	return ApplyColor(ColorLightYellow, text)
}

// LightBlue applies a given color.
func LightBlue(text string) string {
	return ApplyColor(ColorLightBlue, text)
}

// LightPurple applies a given color.
func LightPurple(text string) string {
	return ApplyColor(ColorLightPurple, text)
}

// LightCyan applies a given color.
func LightCyan(text string) string {
	return ApplyColor(ColorLightCyan, text)
}

// LightWhite applies a given color.
func LightWhite(text string) string {
	return ApplyColor(ColorLightWhite, text)
}

// Color represents an ansi color code fragment.
type Color string

// Escaped escapes the color for use in the terminal.
func (c Color) Escaped() string {
	return "\033[" + string(c)
}

// Apply applies a color to a given string.
func (c Color) Apply(text string) string {
	return ApplyColor(c, text)
}

// ApplyColor applies a given color.
func ApplyColor(colorCode Color, text string) string {
	return colorCode.Escaped() + text + ColorReset.Escaped()
}

// Color codes
const (
	ColorBlack       Color = "30m"
	ColorRed         Color = "31m"
	ColorGreen       Color = "32m"
	ColorYellow      Color = "33m"
	ColorBlue        Color = "34m"
	ColorPurple      Color = "35m"
	ColorCyan        Color = "36m"
	ColorWhite       Color = "37m"
	ColorLightBlack  Color = "90m"
	ColorLightRed    Color = "91m"
	ColorLightGreen  Color = "92m"
	ColorLightYellow Color = "93m"
	ColorLightBlue   Color = "94m"
	ColorLightPurple Color = "95m"
	ColorLightCyan   Color = "96m"
	ColorLightWhite  Color = "97m"
	ColorGray        Color = ColorLightBlack
	ColorReset       Color = "0m"
)
