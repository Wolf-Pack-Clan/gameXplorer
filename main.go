package main

import (
	"fmt"
	"gameXplorer/utils"
	"os/exec"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"

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
				widget.NewLabel("Select Game Executable"),
				widget.NewForm(
					widget.NewFormItem("Name", nameEntry),
					widget.NewFormItem("Path", container.NewGridWithRows(1, pathEntry, browseButton)),
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
			popup := widget.NewPopUp(popup_content, myWindow.Canvas())
			popup.Resize(fyne.NewSize(350, 200))
			popup.ShowAtPosition(fyne.NewPos(100, 100))
			//popup.Show()
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {}),
		/*widget.NewToolbarAction(theme.ContentCopyIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentPasteIcon(), func() {}),*/
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.HelpIcon(), func() {}),
	)

	banner := canvas.NewImageFromFile("gameXplorer_c.png")
	banner.FillMode = canvas.ImageFillContain
	banner.SetMinSize(fyne.NewSize(384.5/1.5, 117/1.5))
	top_layout := container.NewVBox(toolbar, widget.NewSeparator(), banner, widget.NewSeparator())

	// List games from .desktop files
	games, err := utils.ListGames()
	if err != nil {
		fmt.Println("Error listing games:", err)
		return
	}
	// Create game cards
	var cards []fyne.CanvasObject
	for _, game := range games {
		cards = append(cards, createGameCard(game.Name, game.Desc, game.Exec, game.ExecDir, game.Icon))
	}
	cardGrid := container.NewVScroll(container.NewVBox(container.NewAdaptiveGrid(3, cards...)))
	//thanks to https://github.com/fyne-io/fyne/issues/2825#issuecomment-1060405595

	ct := container.NewBorder(top_layout, nil, cardGrid, nil, nil)
	myWindow.SetContent(ct)

	myWindow.Resize(fyne.NewSize(500, 600))
	myWindow.ShowAndRun()

	if utils.IsWineInstalled() {
		fmt.Println("1")
	} else {
		fmt.Println("0")
	}
}

func createGameCard(title string, desc string, execCommand string, dir string, icon string) fyne.CanvasObject {

	button := widget.NewButton("        Play        ", func() {
		fyne.CurrentApp().SendNotification(fyne.NewNotification("gameXplorer", "Launching "+title+"..."))
		cmd := exec.Command("sh", "-c", execCommand)
		cmd.Dir = dir

		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		fmt.Printf("Output: %s\n", output)
	})
	options := widget.NewButtonWithIcon("", theme.MenuIcon(), func() {})

	_icon := canvas.NewImageFromFile(icon)
	_icon.FillMode = canvas.ImageFillOriginal
	_title := widget.NewLabelWithStyle(title, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	_desc := widget.NewLabel(desc)
	_desc.Wrapping = fyne.TextWrapBreak
	//_desc.TextStyle = fyne.TextStyle{}
	cardLayout := container.NewVBox(_icon, _title, _desc, container.NewHBox(button, options))
	return cardLayout
}
