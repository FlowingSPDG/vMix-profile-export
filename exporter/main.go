package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/sqweek/dialog"

	"github.com/FlowingSPDG/vMix-profile-export/models"
)

const (
	// InputMovie InputType:0 = movie
	InputMovie string = "0"

	// InputImage InputType:1 = image
	InputImage string = "1"

	// InputCaptureDevices InputType:5 = Capture devices
	InputCaptureDevices string = "5"

	// InputAudioDevices InputType:7 = Audio devices
	InputAudioDevices string = "7"

	// InputBlank InputType:12 = Blank
	InputBlank string = "12"

	// InputBrowser InputType:5000 = browser
	InputBrowser string = "5000"

	// InputTitle InputType:9000 = title
	InputTitle string = "9000"
)

func init(){
	os.Setenv("FYNE_FONT","C:\\Windows\\Fonts\\meiryo.ttc")
}
func main() {
	log.Println("STARTING...")

	// Declare variables
	profile := &models.Profile{}

	// Init GUI
	app := app.New()
	app.Settings().SetTheme(theme.DarkTheme())
   
	window := app.NewWindow("vMix Profile Exporter")
	window.Resize(fyne.NewSize(300, 300))
   
	// Message box
	msgbox := widget.NewEntry()
	msgbox.ReadOnly = true

	// Tool bars
	toolbar := widget.NewToolbar(
	  widget.NewToolbarAction(theme.FolderOpenIcon(), func() {
		// Open vmix profile...
		fname, err := dialog.File().Filter("vMix profile(.vmix)", "vmix").Title("Open profile").Load()
		if err != nil {
			log.Println("Failed to load profile :",err)
			msgbox.SetText(err.Error())
			return
		}
		log.Println("File load ok...")

		// Parse XML
		prof, err := ioutil.ReadFile(fname)
		if err := xml.Unmarshal(prof,profile); err != nil {
			log.Println("Failed to unmarshal profile XML :",err)
			msgbox.SetText(err.Error())
			return
		}
		log.Println("Marshal ok...")

		t := fmt.Sprintf("vMix profile version : %s\n%d Inputs found\n",profile.Version,len(profile.Input))
		for i:=0;i<len(profile.Input);i++{
			// InputNumber / Name / Type / Path
			t += fmt.Sprintf("%d / %s / %s / %s\n",i+1,profile.Input[i].OriginalTitle,profile.Input[i].Type,profile.Input[i].Text)
		}
		msgbox.SetText(t)
		
		// Replace paths...
		for _,v := range profile.Input {
			if v.Type == InputMovie || v.Type == InputImage || v.Type == InputTitle {
				log.Printf("Input Type : %s : Path : %s\n",v.Type,v.Text)
			}
		}

		// Save replaced profile and assets
	  }),
	)
   
	box := fyne.NewContainerWithLayout(
	  layout.NewBorderLayout(toolbar, nil, nil, nil),
	  toolbar, msgbox,
	)
	window.SetContent(box)
   
	window.ShowAndRun()
}
