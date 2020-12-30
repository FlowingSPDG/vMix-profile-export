package main

import (
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	// "github.com/FlowingSPDG/vMix-profile-export/models"
)

 
func main(){
	a := app.New()
	w := a.NewWindow("vMix Profile Exporter")

	hello := widget.NewLabel("Hello Fyne!")
	w.SetContent(widget.NewVBox(
		hello,
		widget.NewButton("Hi!", func() {
			hello.SetText("Welcome :)")
		}),
	))

	w.ShowAndRun()
}