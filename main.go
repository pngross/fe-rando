package main

import (
	"fmt"
	random "math/rand/v2"
)

type feChar struct {
	name            string
	classSet        string
	specialProperty string
	className       string
}

type feClass struct {
	name       string
	classSet   string
	personal   string
	amountLeft int
}

// Controlled by settings file
type randomizerSettings struct {
	game             string
	numberOfUnits    int
	numberPerClass   int
	useGaidens       TernaryBool // FE11 only
	useMaleCrossover TernaryBool // FE12 only
	route            string      // FE16 only
	forceDancer      TernaryBool // FE12/FE16 only
}

type TernaryBool int

const (
	Unclear TernaryBool = iota
	Yes
	No
)

func main() {
	fmt.Println("Hello randomizer")

	settings, err := readSettings()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = validateSettings(settings)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	allChars := []feChar{}
	forcedChars := []feChar{}

	dancer := feClass{"Dancer", "N", "", 1}

	forcedChars, allChars = readAllUnits(settings)

	if settings.game == "FE12" || settings.game == "FE16" {
		forcedChars = append(forcedChars, generateAvatarUnit(settings))
	}

	outputList := randomizeList(allChars, settings.numberOfUnits-len(forcedChars))
	outputList = append(forcedChars, outputList...)

	listOfClasses := readAllClasses(settings.numberPerClass, settings)
	classSlotsLeft := len(listOfClasses) * settings.numberPerClass

	randomNumber, zaehler, classIndex := 0, 0, 0

	// check the force dancer setting (FE16 only) -> reroll for Byleth, Silver Snow!Hilda or faculty members
	if settings.game == "FE16" {
		if settings.forceDancer == Yes {
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

	fmt.Println("Game: " + settings.game)
	if (settings.game == "FE16") || (settings.game == "FE14") { // added FE14 scenario because I might add an implementation for FE14
		fmt.Println("Route: " + settings.route)
	}
	// print result to console
	for i := 0; i < len(outputList); i++ {
		fmt.Println(outputList[i].className + "!" + outputList[i].name)
	}
}
