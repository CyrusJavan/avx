package color

type Color string

const (
	Reset Color = "\033[0m"
	Red   Color = "\033[31m"
	Green Color = "\033[32m"
)

func Sprint(s string, c Color) string {
	return string(c) + s + string(Reset)
}