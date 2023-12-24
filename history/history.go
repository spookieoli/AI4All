package history

// History object is used to store the history of the chat
type History struct {
	chat []*Element
}

// His is a singleton and will be used to store the history of the chat
var His History

// Init will instantiate the history object
func init() {
	His = History{}
}

// AddElement will add a new element to the history
func (h *History) AddElement(from, text string) {
	// Delete the oldest element if the history is too long (max 6 elements)
	if len(h.chat) > 5 {
		h.chat = h.chat[1:]
	}
	h.chat = append(h.chat, CreateElement(from, text))
}

// GetHistory will return the Text in the format "from: text"
func (h *History) GetHistory() string {
	history := ""
	for _, element := range h.chat {
		history = history + element.from + ": " + element.text + "\n"
	}
	return history
}

// ClearHistory will clear the history
func (h *History) ClearHistory() {
	h.chat = []*Element{}
}
