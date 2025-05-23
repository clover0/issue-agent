package logger

type Color string

const (
	Reset   Color = "\033[0m"
	Green   Color = "\033[32m"
	Yellow  Color = "\033[33m"
	Blue    Color = "\033[34m"
	Red     Color = "\033[31m"
	Magenta Color = "\033[35m"
	Cyan    Color = "\033[36m"
	Gray    Color = "\033[37m"
	White   Color = "\033[97m"
)

type ColorFunc func(string) string

func (c Color) String() string {
	return string(c)
}

func GetColorize(color Color) ColorFunc {
	return func(str string) string {
		return color.String() + str + Reset.String()
	}
}
