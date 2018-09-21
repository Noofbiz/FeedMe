package main

import (
	"log"
	"os"

	"github.com/Noofbiz/FeedMe/model"
	"github.com/Noofbiz/FeedMe/server"

	"github.com/asticode/go-astilectron"
)

func main() {
	model.TheFeedData.Open()
	defer model.TheFeedData.Close()
	port := server.StartServer()

	p, _ := os.Getwd()
	var a *astilectron.Astilectron
	var err error
	if a, err = astilectron.New(astilectron.Options{
		AppName:           "Feed Me!",
		BaseDirectoryPath: p,
	}); err != nil {
		log.Fatalf("Failed to create new astillectron. Error: %v", err.Error())
	}
	defer a.Close()
	a.HandleSignals()

	if err = a.Start(); err != nil {
		log.Fatalf("Failed to start. Error: %v", err.Error())
	}

	var w *astilectron.Window
	if w, err = a.NewWindow("http://localhost:"+port, &astilectron.WindowOptions{
		Center:          astilectron.PtrBool(true),
		Width:           astilectron.PtrInt(740),
		Height:          astilectron.PtrInt(960),
		BackgroundColor: astilectron.PtrStr("#2B3E50"),
	}); err != nil {
		log.Fatalf("Failed to create new window. Error: %v", err.Error())
	}
	if err = w.Create(); err != nil {
		log.Fatalf("Failed at window.Create(). Error: %v", err.Error())
	}

	w.OpenDevTools()

	a.Wait()
}
