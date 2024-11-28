package logger

type Color string

const (
	reset   Color = "\033[0m"
	red     Color = "\033[31m"
	green   Color = "\033[32m"
	yellow  Color = "\033[33m"
	blue    Color = "\033[34m"
	magenta Color = "\033[35m"
	cyan    Color = "\033[36m"
	gray    Color = "\033[37m"
	white   Color = "\033[97m"
)

func (c Color) String() string {
	return string(c)
}

func Green(str string) string {
	return green.String() + str + reset.String()
}
