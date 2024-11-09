package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func readAllUnits(settings randomizerSettings) ([]feChar, []feChar) {
	forcedChars, allChars := []feChar{}, []feChar{}

	f, err := os.Open("data/units_" + settings.game + ".csv")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.Comma = ';'
	inputData, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println(err.Error())
	}

	var unitRoutes string
	currentChar := feChar{}

	for i, line := range inputData {
		if i > 0 {
			for j, field := range line {
				if j == 0 {
					currentChar.name = field
				} else if j == 1 {
					currentChar.classSet = field
				} else if j == 2 {
					unitRoutes = field
				} else if j == 3 {
					currentChar.specialProperty = field
				}
			}
			if settings.game == "FE16" {
				if ((settings.route == "CF") && ((unitRoutes == "all") || (unitRoutes == "cf"))) || ((settings.route == "VW") && ((unitRoutes == "all") || (unitRoutes == "notcf") || (unitRoutes == "vw"))) || ((settings.route == "AM") && ((unitRoutes == "all") || (unitRoutes == "notcf") || (unitRoutes == "am"))) || ((settings.route == "SS") && ((unitRoutes == "all") || (unitRoutes == "notcf"))) {
					if currentChar.specialProperty == "lord" {
						forcedChars = append(forcedChars, currentChar)
					} else {
						allChars = append(allChars, currentChar)
					}
				}
			} else {
				allChars = append(allChars, currentChar)
			}
		}
	}

	return forcedChars, allChars
}

func readAllClasses(amount int, settings randomizerSettings) []feClass {
	f, err := os.Open("data/classes_" + settings.game + ".csv")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.Comma = ';'
	inputData, err := csvReader.ReadAll()
	if (err != nil) && (inputData == nil) {
		fmt.Println(err.Error())
	}
	allClasses := []feClass{}
	currentClass := feClass{}

	for i, line := range inputData {
		if i > 0 {
			for j, field := range line {
				if j == 0 {
					currentClass.name = field
				} else if j == 1 {
					currentClass.classSet = field
				} else if j == 2 {
					currentClass.personal = field
				}
			}
			currentClass.amountLeft = amount
			allClasses = append(allClasses, currentClass)
		}
	}

	return allClasses
}
