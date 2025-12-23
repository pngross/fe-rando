package main

import (
	"fmt"
	"slices"
	"strconv"

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
	writeFunc := r.settings.MakeUIBinding(menu)

	// TODO INSERT SETTINGS HERE
	// additional idea: allow default settings (one per game) to be saved and loaded via a json file

	menu.Add(widget.NewButton("Randomize!", func() {
		writeFunc()
		if !r.settings.Validate() {
			return
		}
		r.Randomize()
	}))
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

func (s *randomizerSettings) Validate() bool {
	if !slices.Contains([]string{"FE11", "FE12", "FE16"}, s.game) {
		return false
	}

	if s.game == "FE16" && !slices.Contains([]string{"VW", "AM", "SS", "CF"}, s.route) {
		return false
	}

	if s.numberOfUnits < 0 || s.numberPerClass < 0 {
		return false
	}

	return true
}

func (settings *randomizerSettings) MakeUIBinding(menu *fyne.Container) func() {
	unitNumBox := widget.NewEntry()
	unitNumBox.Text = fmt.Sprintf("%d", settings.numberOfUnits)
	perclassNumBox := widget.NewEntry()
	perclassNumBox.Text = fmt.Sprintf("%d", settings.numberPerClass)

	routeDropdown := widget.Select{
		Options: []string{"AM", "VW", "SS", "CF"},
	}

	menu.Add(unitNumBox)
	menu.Add(perclassNumBox)

	forceJagen := widget.Check{
		Text:    "Force Jagen",
		Checked: settings.forceJagen,
	}

	forceDancer := widget.Check{
		Text:    "Force Dancer",
		Checked: settings.forceDancer,
	}

	maleCrossover := widget.Check{
		Text:    "Male Crossover Classes",
		Checked: settings.useMaleCrossover,
	}

	useGaidens := widget.Check{
		Text:    "Use Gaidens",
		Checked: settings.useGaidens,
	}

	if settings.game == "FE11" {
		menu.Add(&useGaidens)
	} else {
		menu.Add(&forceDancer)
	}

	if settings.game == "FE16" {
		menu.Add(&routeDropdown)
	} else {
		menu.Add(&forceJagen)
	}

	if settings.game == "FE12" {
		menu.Add(&maleCrossover)
	}

	writeFunc := func() {
		settings.numberOfUnits = parseNumWithDefault(unitNumBox.Text, settings.numberOfUnits)
		settings.numberPerClass = parseNumWithDefault(perclassNumBox.Text, settings.numberPerClass)
		settings.forceDancer = forceDancer.Checked
		settings.forceJagen = forceJagen.Checked
		settings.useGaidens = useGaidens.Checked
		settings.route = routeDropdown.Selected
		settings.useMaleCrossover = maleCrossover.Checked
	}
	return writeFunc
}

func parseNumWithDefault(str string, deflt int) int {
	if num, err := strconv.Atoi(str); err == nil {
		return num
	}
	return deflt
}
