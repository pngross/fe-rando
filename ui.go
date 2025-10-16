package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Randomizer struct {
	app      fyne.App
	window   fyne.Window
	settings randomizerSettings
}

func InitUI() Randomizer {
	application := app.New()
	w := application.NewWindow("FE Team Randomizer")
	rando := Randomizer{settings: randomizerSettings{}, app: application, window: w}
	rando.ToMainMenu()
	return rando
}

func (r *Randomizer) ToMainMenu() {
	menu := container.NewVBox()
	menu.Add(widget.NewLabel("Randomize your FE team!"))
	for _, game := range supportedGames {
		btn := widget.NewButton(game, func() { r.SelectGame(game) })
		menu.Add(btn)
	}
	r.window.SetContent(menu)
}

func (r *Randomizer) SelectGame(game string) {
	menu := container.NewVBox()
	menu.Add(widget.NewLabel(game + " - Settings:"))

	r.settings, _ = readDefaultSettings(game)

	// TODO INSERT SETTINGS HERE
	// additional idea: allow default settings (one per game) to be saved and loaded via a json file

	menu.Add(widget.NewButton("Randomize!", func() { r.Randomize() }))
	menu.Add(r.ReturnBtn())
	r.window.SetContent(menu)
}

func (r *Randomizer) Randomize() {
	randomizedTeam := RandomizeTeam(r.settings)

	resultsPage := container.NewVBox()

	textbox := widget.NewLabel(randomizedTeam)
	resultsPage.Add(textbox)

	resultsPage.Add(widget.NewLabel("Save to: "))
	filenameEntry := widget.NewEntry()
	filenameEntry.Text = fmt.Sprintf(`out\%s_team.txt`, r.settings.game)
	resultsPage.Add(filenameEntry)

	buttons := container.NewHBox()
	buttons.Add(widget.NewButton("Save", func() { SaveToFile(randomizedTeam, filenameEntry.Text) }))
	buttons.Add(widget.NewButton("Reroll", func() { r.Randomize() }))
	buttons.Add(r.ReturnBtn())

	resultsPage.Add(buttons)
	r.window.SetContent(resultsPage)
}

func (r *Randomizer) ReturnBtn() *widget.Button {
	return widget.NewButton("Back to Menu", func() { r.ToMainMenu() })
}
