package history

// Element is a struct that will store one history element

type Element struct {
	from, text string
}

func CreateElement(from, text string) *Element {
	return &Element{from, text}
}
