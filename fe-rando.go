package main

import (
	"fmt"
	random "math/rand/v2"
)

func main() {
	application := InitUI()
	application.window.ShowAndRun()
}

func RandomizeTeam(settings randomizerSettings) string {

	output := fmt.Sprintf("Game: %s\n", settings.game)
	if (settings.game == "FE16") || (settings.game == "FE14") { // added FE14 scenario because I might add an implementation for FE14
		output += fmt.Sprintf("Route: %s\n", settings.route)
	}

	dancer := feClass{"Dancer", "N", "", 1}
	forcedChars, allChars, freeChars := readAllUnits(settings)

	if settings.game == "FE12" || settings.game == "FE16" {
		avatar := generateAvatarUnit(settings)
		forcedChars = append(forcedChars, avatar)
		if avatar.classSet == "F" {
			output += fmt.Sprintf("Female %s was chosen!\n", avatar.name)
		} else {
			output += fmt.Sprintf("Male %s was chosen!\n", avatar.name)
		}
	}

	outputList := randomizeList(allChars, settings.numberOfUnits-len(forcedChars))
	outputList = append(forcedChars, outputList...)

	listOfClasses := readAllClasses(settings.numberPerClass, settings)
	classSlotsLeft := len(listOfClasses) * settings.numberPerClass

	randomNumber, zaehler, classIndex := 0, 0, 0

	// check the force dancer setting (FE16 only) -> reroll for Byleth, Silver Snow!Hilda or faculty members
	if settings.game == "FE16" {
		if settings.forceDancer {
			dancerFound := false
			for !dancerFound {
				randomNumber = random.IntN(len(outputList))
				if matchClass(dancer, outputList[randomNumber], settings) {
					outputList[randomNumber].className = "Dancer"
					dancerFound = true
				}
			}
		} else {
			listOfClasses = append(listOfClasses, dancer)
			classSlotsLeft += 1
		}
	}

	// randomly assign classes to output list
	for i := 0; i < len(outputList); i++ {
		for outputList[i].className == "" {
			if (classSlotsLeft == 0) || !checkForValidClasses(listOfClasses, outputList[i], settings) {
				classIndex = random.IntN(len(listOfClasses))
			} else {
				randomNumber = random.IntN(classSlotsLeft)
				zaehler, classIndex = 0, -1
				for j := 0; j < len(listOfClasses); j++ {
					if (zaehler >= randomNumber) && (listOfClasses[j].amountLeft != 0) {
						classIndex = j
						break
					} else {
						zaehler += listOfClasses[j].amountLeft
					}

				}
			}

			if classIndex >= 0 {
				if matchClass(listOfClasses[classIndex], outputList[i], settings) {
					outputList[i].className = listOfClasses[classIndex].name
					if listOfClasses[classIndex].amountLeft > 0 {
						listOfClasses[classIndex].amountLeft--
					}
				}
			}
		}
	}

	// print result to console
	for i := 0; i < len(outputList); i++ {
		output += fmt.Sprintf("%s!%s\n", outputList[i].className, outputList[i].name)
	}

	if len(freeChars) > 0 {
		output += "\nFree Units: "
		for i := 0; i < len(freeChars); i++ {
			output += freeChars[i].name + " "
		}
	}
	return output
}
