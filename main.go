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
	myApp := app.NewWithID("org.kazam.gameXplorer")
	myWindow := myApp.NewWindow("gameXplorer")

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Name")
	descEntry := widget.NewEntry()
	descEntry.SetPlaceHolder("Description")
	pathEntry := widget.NewEntry()
	pathEntry.SetPlaceHolder("Select game path")

	browseButton := widget.NewButton("Browse...", func() {
		pathEntry.FocusGained()
		dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				pathEntry.SetText(reader.URI().Path())
			}
		}, myWindow).Show()
	})
	var shared_game bool
	shared_radio := widget.NewCheck("(must be running as root)", func(value bool) {
		if value {
			shared_game = true
		} else {
			shared_game = false
		}
	})

	var NewGamePopUp *widget.PopUp
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			popup_content := container.NewVBox(
				widget.NewLabel("Add New Game"),
				widget.NewForm(
					widget.NewFormItem("Name", nameEntry),
					widget.NewFormItem("Description", descEntry),
					widget.NewFormItem("Path", container.NewGridWithRows(1, pathEntry, browseButton)),
					widget.NewFormItem("Shared", shared_radio),
				),
				widget.NewButton("Save", func() {
					fmt.Println("Name:", nameEntry.Text)
					fmt.Println("Description:", descEntry.Text)
					fmt.Println("Game Path:", pathEntry.Text)
					utils.SaveGame(nameEntry.Text, descEntry.Text, pathEntry.Text, shared_game)
					nameEntry.SetText("")
					descEntry.SetText("")
					pathEntry.SetText("")
					shared_radio.Checked = false
					NewGamePopUp.Hide()
				}),
				widget.NewButton("Close", func() {
					nameEntry.SetText("")
					descEntry.SetText("")
					pathEntry.SetText("")
					shared_radio.Checked = false
					NewGamePopUp.Hide()
				}),
			)
			NewGamePopUp = widget.NewModalPopUp(popup_content, myWindow.Canvas())
			NewGamePopUp.Resize(fyne.NewSize(350, 200))
			NewGamePopUp.Show()
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {}),
		/*widget.NewToolbarAction(theme.ContentCopyIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentPasteIcon(), func() {}),*/
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.HelpIcon(), func() {}),
	)

	banner := canvas.NewImageFromFile("gameXplorer_c_s.png")
	banner.FillMode = canvas.ImageFillOriginal
	//banner.SetMinSize(fyne.NewSize(387.5, 117.5))
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

	/*if utils.IsWineInstalled() {
		fmt.Println("1")
	} else {
		fmt.Println("0")
	}*/
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
	utils.ExtractEXEIcon("/mnt/localdrive/cod/c1cx.exe")
	options := widget.NewButtonWithIcon("", theme.MenuIcon(), func() {})

	_icon := canvas.NewImageFromFile(icon)
	_icon.FillMode = canvas.ImageFillOriginal
	_title := widget.NewLabelWithStyle(title, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	_desc := widget.NewLabel(desc)
	_desc.Wrapping = fyne.TextWrapBreak
	cardLayout := container.NewVBox(_icon, _title, _desc, container.NewHBox(button, options))
	return cardLayout
}
