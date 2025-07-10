package output

type Output interface {
	Type(text string)
	Backspace(n int)
}
