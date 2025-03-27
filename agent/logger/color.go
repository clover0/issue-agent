package logger

type Color string

const (
	reset   Color = "\033[0m"
	green   Color = "\033[32m"
	yellow  Color = "\033[33m"
	blue    Color = "\033[34m"
	red     Color = "\033[31m"
	magenta Color = "\033[35m"
	cyan    Color = "\033[36m"
	gray    Color = "\033[37m"
	white   Color = "\033[97m"
)

var noColor = false

func SetNoColor() {
	noColor = true
}

func (c Color) String() string {
	return string(c)
}

func Green(str string) string {
	if noColor {
		return str
	}
	return green.String() + str + reset.String()
}

func Yellow(str string) string {
	if noColor {
		return str
	}
	return yellow.String() + str + reset.String()
}

func Blue(str string) string {
	if noColor {
		return str
	}
	return blue.String() + str + reset.String()
}

func Red(str string) string {
	if noColor {
		return str
	}
	return red.String() + str + red.String()
}
