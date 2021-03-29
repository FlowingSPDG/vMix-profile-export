package main

import (
	"encoding/xml"
	"flag"
	"io"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/FlowingSPDG/vMix-profile-export/models"
	"github.com/sirupsen/logrus"
	"github.com/sqweek/dialog"
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

var (
	// ImportProfilePath vMix profile path to import
	ImportProfilePath string

	// exportPath path to export
	exportPath string
)

func init() {
	// Parse args
	flag.Parse()

	// setup import file path
	ImportProfilePath = flag.Arg(0)
	if ImportProfilePath == "" {
		// Open vmix profile...
		fname, err := dialog.File().Filter("vMix profile(.vmix)", "vmix").Title("Open profile").Load()
		if err != nil {
			panic(err)
		}
		logrus.Debugln("File load ok...")
		ImportProfilePath = fname
	}

	// setup export path
	exportPath = flag.Arg(1)
	if exportPath == "" {
		// Specify destination...
		directory, err := dialog.Directory().Title("Directory to save assets").Browse()
		if err != nil {
			panic(err)
		}
		logrus.Debugln("File load ok...")
		exportPath = directory
	}

	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	logrus.Infoln("STARTING...")

	// Declare variables
	profile := &models.Profile{}

	// Parse XML
	prof, err := os.ReadFile(ImportProfilePath)
	if err != nil {
		logrus.Fatalln("Failed to read profile file :", err)
	}
	logrus.Debugln("Read ok...")

	if err := xml.Unmarshal(prof, profile); err != nil {
		logrus.Fatalln("Failed to unmarshal profile XML :", err)
		panic(err)
	}
	logrus.Debugln("Marshal ok...")

	logrus.Infof("vMix profile version : %s\n", profile.Version)
	logrus.Infof("%d Inputs found\n", len(profile.Input))
	for i := 0; i < len(profile.Input); i++ {
		// InputNumber / Name / Type / Path
		logrus.Debugf("%d / %s / %s / %s\n", i+1, profile.Input[i].OriginalTitle, profile.Input[i].Type, profile.Input[i].Text)
	}

	// goroutine setup
	wg := sync.WaitGroup{}
	for i := 0; i < len(profile.Input); i++ {
		if profile.Input[i].Type == InputMovie || profile.Input[i].Type == InputImage {
			logrus.Debugf("Input Type : %s : Path : %s\n", profile.Input[i].Type, profile.Input[i].Text)
			wg.Add(1)
			go func(v *models.Input) {
				logrus.Debugf("v.Text: %#v", v.Text)
				defer wg.Done()

				// Resolve file paths
				absPath := strings.ReplaceAll(v.Text, "\\", "/")
				filename := path.Base(absPath)
				logrus.Debugf("Name: %s Abs path:%s, filename:%s\n", v.Title, absPath, filename)

				// Copy/Save all static assets...
				source, err := os.Open(absPath)
				if err != nil {
					logrus.Errorln("Failed to read file :", err)
					return
				}
				defer source.Close()

				// make destination folder and profile
				destpath := path.Join(exportPath, filename)
				logrus.Debugln("Creating file:", destpath)
				destination, err := os.Create(destpath)
				if err != nil {
					logrus.Errorln("Failed to make file :", err)
					return
				}
				defer destination.Close()

				// copy assets?
				logrus.Debugf("Start copy. src:%s dest:%s\n:", absPath, destpath)
				_, err = io.Copy(destination, source)
				if err != nil {
					logrus.Errorln("Failed to copy file :", err)
					return
				}

				// Replace paths(relative)
				v.Text = destpath

				// TODO: Save CRC checks
			}(&profile.Input[i])
		}
	}
	wg.Wait()

	// Write profile itself
	profPath := path.Join(exportPath, "profile.vmix")
	b, _ := xml.Marshal(profile)
	os.WriteFile(profPath, b, 0777)
	logrus.Infoln("DONE")
}
