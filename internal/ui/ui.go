// Package ui provides the fyne.io GUI for go-t3.
package ui

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

// Run starts the fyne application window.
func Run() {
	a := app.New()
	w := a.NewWindow("go-t3: Tic Tac Toe")
	w.SetContent(widget.NewLabel("Tic Tac Toe — coming soon"))
	w.ShowAndRun()
}
