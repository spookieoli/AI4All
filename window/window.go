package window

import (
	"AI4All/modelloader"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"strings"
)

// ChatWindow
type ChatWindow struct {
	Title      string
	App        fyne.App
	Win        fyne.Window
	Output     *ExtendedEntry
	Input      *InputEntry
	SendButton *widget.Button
	ChatText   string
	Model      *modelloader.ModelLoader
	Scroller   *container.Scroll
	prompt_num bool
}

/****************************************************************************/
/* Types for building the Window											*/
/****************************************************************************/

// ExtendedEntry will overwrite the Tapped Function TODO: Add NewExtendedEntry
type ExtendedEntry struct {
	widget.Entry
	input *InputEntry
}

// InputEntry will overwrite the Tapped Function
type InputEntry struct {
	widget.Entry
	cw *ChatWindow
}

// CustomTheme will overwrite TextColor
type CustomTheme struct {
	fyne.Theme
}

// NewWindow creates a new window - TODO: Make it beautify
func NewWindow(title string) *ChatWindow {
	// Create a new ChatWindow
	cw := &ChatWindow{
		Title:      title,
		Model:      nil,
		prompt_num: true,
	}
	// Create a new App
	cw.App = app.New()
	cw.App.Settings().SetTheme(&CustomTheme{Theme: cw.App.Settings().Theme()})
	cw.Win = cw.App.NewWindow(title)
	// Fix the size of the window
	cw.Win.SetFixedSize(true)
	// Create a ChatInput Entry Widget
	cw.Input = NewInputEntry(cw)
	// Create a Container for the ChatInput
	input := container.NewVScroll(cw.Input)
	// Create a ChatOutpt Entry Widget
	cw.Output = NewExtendedEntry(cw.Input)
	cw.Output.OnChanged = cw.Changed
	// Create a Container for the ChatOutput
	cw.Scroller = container.NewVScroll(cw.Output)
	cw.Scroller.SetMinSize(fyne.NewSize(0, 350))
	// Create a Send Button
	sendButton := widget.NewButton("Send", cw.Send)
	sendButton.Resize(fyne.NewSize(50, 100))
	cw.SendButton = sendButton
	// Create a Clear Button
	clearButton := widget.NewButton("Clear", cw.Clear)
	// Set the size of the ChatInput
	input.SetMinSize(fyne.NewSize(750, 100))
	// Create a horizontal container
	horizontal := container.NewHBox(input, sendButton, clearButton)
	// Create a Grid Row for the ChatOutput and ChatInput
	content := container.NewVBox(cw.Scroller, horizontal)
	// Set the content of the window
	cw.Win.SetContent(content)
	// Create a new Menu
	cw.CreateMainMenu()
	// add key shortcuts
	shiftEnter := &desktop.CustomShortcut{KeyName: fyne.KeyReturn, Modifier: fyne.KeyModifierControl}
	cw.Win.Canvas().AddShortcut(shiftEnter, func(shortcut fyne.Shortcut) {
		sendButton.Tapped(nil)
	})
	// Resize to 800x600
	cw.Win.Resize(fyne.NewSize(800, 400))
	// Set to fixed size
	cw.Win.SetFixedSize(true)
	return cw
}

// Send will send the text from the ChatInput to the ChatOutput
func (w *ChatWindow) Send() {
	// TODO: Should be loaded from config
	system := "System: Your role is to be a helpful and friendly chat assistant. " +
		"Your task is to provide accurate, short and understandable answers to a wide range of questions. "
	inputtext := w.Input.Text
	if inputtext != "" {
		// Update the ChatText
		w.ChatText = w.ChatText + "\n" + "You: \n" + w.Input.Text + "\n"
		w.Input.SetText("")
		w.Output.SetText(w.ChatText)
		w.Output.ToBottom()
		// Send the Text to the Model
		if w.Model != nil {
			answer, err := w.Model.Predictor(system+"User: "+inputtext+" Assistant: ", []string{"User:"}, 512, 47, 1, 0.1)
			if err != nil {
				fmt.Println("Error while predicting:", err.Error())
				return
			}
			// Update the ChatText
			answer = strings.Trim(answer, " ")
			answer = strings.Trim(answer, "\n")
			w.ChatText = w.ChatText + "AI: \n" + answer + "\n"
			w.Output.SetText(w.ChatText)
			w.Output.Refresh()
			w.Output.ToBottom()
		}
	}
}

// Changed will update the ChatOutput
func (w *ChatWindow) Changed(text string) {
	w.Output.SetText(w.ChatText)
}

// Clear will clear the ChatOutput
func (w *ChatWindow) Clear() {
	w.ChatText = ""
	w.Output.SetText(w.ChatText)
}

// MouseDown will overwrite the MouseDown in Chatoutput Function
func (c *ExtendedEntry) MouseDown(e *desktop.MouseEvent) {
	c.input.Entry.MouseDown(e)
}

// CreateMainMenu creates a new Menu
func (cw *ChatWindow) CreateMainMenu() {
	// The Menueitems
	loadItem := fyne.NewMenuItem("Load Model", func() {
		d := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, cw.Win)
				return
			}
			cw.LoadModel(reader.URI().Path())
		}, cw.Win)
		// Show Only gguf Data
		d.SetFilter(storage.NewExtensionFileFilter([]string{".gguf"}))
		d.Show()
	})
	configItem := fyne.NewMenuItem("Config", func() {
		// TODO: Add Config
		dialog.ShowInformation("About", "We will add a config soon", cw.Win)
	})
	quitItem := fyne.NewMenuItem("Quit", func() {
		cw.App.Quit()
	})
	aboutItem := fyne.NewMenuItem("About", func() {
		// Opens up a Dialog with informations about the program
		dialog.ShowInformation("About", "AI4All\nVersion: 0.1.0\nAuthor: AI4All", cw.Win)
	})
	// Create the Main Menu
	mainMenu := fyne.NewMainMenu(
		fyne.NewMenu("File", loadItem, configItem, quitItem),
		fyne.NewMenu("About", aboutItem),
	)
	cw.Win.SetMainMenu(mainMenu)
}

// Change the TextColor to white
func (c CustomTheme) TextColor() color.Color {
	return color.White
}

// NewInputEntry creates a new InputEntry
func NewInputEntry(cw *ChatWindow) *InputEntry {
	entry := &InputEntry{cw: cw}
	entry.ExtendBaseWidget(entry) // Initialisiert das Basiselement
	entry.SetPlaceHolder("Enter your message here...")
	return entry
}

// NewExtendedEntry creates a new ExtendedEntry
func NewExtendedEntry(input *InputEntry) *ExtendedEntry {
	entry := &ExtendedEntry{input: input}
	entry.MultiLine = true
	entry.ExtendBaseWidget(entry)
	entry.Wrapping = fyne.TextWrapWord
	entry.TextStyle.Bold = true
	entry.TextStyle.Monospace = true
	return entry
}

// Will overwrite the TypedShortcut Function
func (m *InputEntry) TypedShortcut(s fyne.Shortcut) {
	if _, ok := s.(*desktop.CustomShortcut); !ok {
		m.Entry.TypedShortcut(s)
		return
	}
	m.cw.SendButton.Tapped(nil)
}

// LoadModel loads a model - TODO: Add config from configuration window
func (cw *ChatWindow) LoadModel(path string) {
	// Create a new ModelLoader
	ml, err := modelloader.NewModelLoader(path, 4, 512, 0)
	if err != nil {
		fmt.Println("Loading the model failed:", err.Error())
		return
	}
	// Set the Model
	cw.Model = ml
}

// Scroll to bottom to show all the new Text
func (e *ExtendedEntry) ToBottom() {
	e.CursorRow = len(e.Text) - 1
}
