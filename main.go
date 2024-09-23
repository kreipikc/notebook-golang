package main

import (
	"fmt"
	"io"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Notebook")
	w.Resize(fyne.NewSize(950, 750))

	tabs := AddTabs(w)
	entries := AddFirstTab(tabs)
	AddMenu(w, a, tabs, entries)

	w.Show()
	a.Run()
}

func NewTextEntry() *widget.Entry {
	textInput := widget.NewMultiLineEntry()
	textInput.SetPlaceHolder("Enter text here...")

	return textInput
}

func AddTabs(w fyne.Window) *container.AppTabs {
	tabs := container.NewAppTabs()
	w.SetContent(tabs)

	return tabs
}

func AddFirstTab(tabs *container.AppTabs) []*widget.Entry {
	entries := []*widget.Entry{}
	firstEntry := NewTextEntry()
	tabs.Append(container.NewTabItem("untitled", firstEntry))

	return append(entries, firstEntry)
}

func AddMenu(w fyne.Window, a fyne.App, tabs *container.AppTabs, entries []*widget.Entry) {
	new := fyne.NewMenuItem("New", func() {
		entry := NewTextEntry()
		entries = append(entries, entry)
		tabs.Append(container.NewTabItem("untitled", entry))
	})

	save := fyne.NewMenuItem("Save", func() {
		tName := tabs.Selected().Text

		if tName == "untitled" {
			SaveFile(w, tabs, entries)
		} else {
			file, _ := os.Create(tName)
			defer file.Close()

			file.WriteString(entries[tabs.SelectedIndex()].Text)
		}
	})

	saveAs := fyne.NewMenuItem("Save As...", func() { SaveFile(w, tabs, entries) })

	open := fyne.NewMenuItem("Open...", func() { OpenFile(w, tabs, &entries) })

	close := fyne.NewMenuItem("Close", func() { CloseFile(w, tabs, &entries) })

	menu := fyne.NewMenu("File", new, save, saveAs, open, close)

	theme := fyne.NewMenuItem("Theme", nil)
	theme.ChildMenu = fyne.NewMenu(
		"Theme",
		fyne.NewMenuItem("Dark", func() { setDarkTheme(a) }),
		fyne.NewMenuItem("Light", func() { setLightTheme(a) }),
	)

	settings := fyne.NewMenu("Settings", theme)

	w.SetMainMenu(fyne.NewMainMenu(menu, settings))
}

func setDarkTheme(a fyne.App) {
	a.Settings().SetTheme(theme.DarkTheme())
}

func setLightTheme(a fyne.App) {
	a.Settings().SetTheme(theme.LightTheme())
}

func SaveFile(w fyne.Window, tabs *container.AppTabs, entries []*widget.Entry) {
	dialog.ShowFileSave(func(uc fyne.URIWriteCloser, err error) {
		if uc != nil {
			io.WriteString(uc, entries[tabs.SelectedIndex()].Text)
			tabs.Selected().Text = uc.URI().Path()
			tabs.Refresh()
		}
	}, w)
}

func OpenFile(w fyne.Window, tabs *container.AppTabs, entries *[]*widget.Entry) {
	dialog.ShowFileOpen(func(uc fyne.URIReadCloser, err error) {
		if uc != nil {
			data, _ := io.ReadAll(uc)
			entry := NewTextEntry()
			entry.Text = string(data)
			*entries = append(*entries, entry)
			tabs.Append(container.NewTabItem(uc.URI().Path(), entry))
			fmt.Println(entries)
		}
	}, w)
}

func CloseFile(w fyne.Window, tabs *container.AppTabs, entries *[]*widget.Entry) {
	index := tabs.SelectedIndex()
	*entries = append((*entries)[:index], (*entries)[index+1:]...)
	tabs.Remove(tabs.Selected())
}
