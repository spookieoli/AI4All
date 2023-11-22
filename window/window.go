package window

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

// ChatWindow
type ChatWindow struct {
	Title    string
	App      fyne.App
	Win      fyne.Window
	Output   *ExtendedEntry
	Input    *widget.Entry
	ChatText string
}

type ExtendedEntry struct {
	widget.Entry
	input *widget.Entry
}

// CustomTheme will overwrite TextColor
type CustomTheme struct {
	fyne.Theme
}

func (c CustomTheme) TextColor() color.Color {
	return color.White
}

// NewWindow creates a new window - TODO: Make it better
func NewWindow(title string) *ChatWindow {
	cw := &ChatWindow{
		Title: title,
	}
	cw.App = app.New()
	cw.App.Settings().SetTheme(&CustomTheme{Theme: cw.App.Settings().Theme()})
	cw.Win = cw.App.NewWindow(title)
	// Fix the size of the window
	cw.Win.SetFixedSize(true)
	// Create a ChatInput Entry Widget
	cw.Input = widget.NewMultiLineEntry()
	// Set the placeholder text
	cw.Input.SetPlaceHolder("Enter your message here...")
	// Set text wrapping to true
	cw.Input.Wrapping = fyne.TextWrapWord
	// Create a Container for the ChatInput
	input := container.NewVScroll(cw.Input)
	// Create a ChatOutpt Entry Widget
	cw.Output = &ExtendedEntry{input: cw.Input}
	cw.Output.OnChanged = cw.Changed
	// Set the Textcolor to white
	cw.Output.TextStyle.Bold = true
	cw.Output.TextStyle.Monospace = true
	// Create a Container for the ChatOutput
	output := container.NewVScroll(cw.Output)
	output.SetMinSize(fyne.NewSize(0, 350))
	// Create a horizontal container
	horizontal := container.NewHBox()
	// Create a Send Button
	sendButton := widget.NewButton("Send", cw.Send)
	sendButton.Resize(fyne.NewSize(50, 100))
	// Create a Clear Button
	clearButton := widget.NewButton("Clear", cw.Clear)
	// Set the size of the ChatInput
	input.SetMinSize(fyne.NewSize(750, 100))
	// Add the input and the button to the horizontal container
	horizontal.Add(input)
	horizontal.Add(sendButton)
	horizontal.Add(clearButton)
	// Create a Grid Row for the ChatOutput and ChatInput
	content := container.NewVBox(output, horizontal)
	// Set the content of the window
	cw.Win.SetContent(content)
	// Create a new Menu
	cw.CreateMainMenu()
	// Resize to 800x600
	cw.Win.Resize(fyne.NewSize(800, 400))
	// Set to fixed size
	cw.Win.SetFixedSize(true)
	return cw
}

// Send will send the text from the ChatInput to the ChatOutput
func (w *ChatWindow) Send() {
	if w.Input.Text != "" {
		w.ChatText = w.ChatText + "\n" + "You: \n" + w.Input.Text + "\n"
		w.Input.SetText("")
		w.Output.SetText(w.ChatText)
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

// Tapped Function of the ExtendedEntry
func (e *ExtendedEntry) Tapped(*fyne.PointEvent) {
	e.input.Enable()
}

// CreateMainMenu creates a new Menu
func (cw *ChatWindow) CreateMainMenu() {
	// The Menueitems
	loadItem := fyne.NewMenuItem("Load Model", func() {
		d := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				cw.App.Quit()
			}
			// TODO: Add Model Loading
			fmt.Println("Model Loading")
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
