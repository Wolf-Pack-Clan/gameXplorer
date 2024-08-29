package main

import (
	"fmt"
	"gameXplorer/utils"
	"os/exec"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"

	//"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Initialize the Fyne app
	myApp := app.New()
	myWindow := myApp.NewWindow("gameXplorer")

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Name")
	pathEntry := widget.NewEntry()
	pathEntry.SetPlaceHolder("Select game path")

	browseButton := widget.NewButton("Browse...", func() {
		dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				pathEntry.SetText(reader.URI().Path())
			}
		}, myWindow).Show()
	})

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			popup_content := container.NewVBox(
				widget.NewLabel("Custom Popup Content"),
				widget.NewForm(
					widget.NewFormItem("Name", nameEntry),
					widget.NewFormItem("Path", container.NewHBox(pathEntry, browseButton)),
				),
				widget.NewButton("Save", func() {
					fmt.Println("Name:", nameEntry.Text)
					fmt.Println("Game Path:", pathEntry.Text)
					myWindow.Close()
				}),
				widget.NewButton("Close", func() {
					myWindow.Close()
				}),
			)
			widget.NewModalPopUp(popup_content, myWindow.Canvas()).Show()
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {}),
		/*widget.NewToolbarAction(theme.ContentCopyIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentPasteIcon(), func() {}),*/
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.HelpIcon(), func() {}),
	)

	// Create the game cards
	appCard := createGameCard("Call of Duty", "Best Shooter", "Play")
	cardGrid := container.NewVScroll(container.NewVBox(container.NewAdaptiveGrid(4, appCard)))
	//thanks to https://github.com/fyne-io/fyne/issues/2825#issuecomment-1060405595

	ct := container.NewBorder(toolbar, nil, cardGrid, nil, nil)
	myWindow.SetContent(ct)

	myWindow.Resize(fyne.NewSize(600, 400))
	myWindow.ShowAndRun()

	if utils.IsWineInstalled() {
		fmt.Println("1")
	} else {
		fmt.Println("0")
	}
}

func createGameCard(title string, desc string, buttonText string) fyne.CanvasObject {

	button := widget.NewButton(buttonText, func() {
		fyne.CurrentApp().SendNotification(fyne.NewNotification("gameXplorer", "Launching "+title+"..."))
		cmd := exec.Command("./run.sh")
		cmd.Dir = "/mnt/localdrive/cod/"

		output, err := cmd.CombinedOutput()

		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}

		fmt.Printf("Output: %s\n", output)
	})

	cardLayout := container.NewVBox(
		button,
	)
	card := widget.NewCard(title, desc, cardLayout)

	return card
}
