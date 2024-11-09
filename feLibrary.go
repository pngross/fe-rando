package main

import (
	random "math/rand/v2"
	"strings"
)

// Check if a class matches a character
func matchClass(class feClass, unit feChar, settings randomizerSettings) bool {
	result := false
	if class.classSet == "P" {
		return strings.ReplaceAll(class.personal, unit.name, "") == class.personal
	} else if settings.game == "FE16" {
		if ((class.classSet == "N") || (class.classSet == unit.classSet)) && ((class.personal == "") || (class.personal == unit.name)) {
			if class.name == "Dancer" {
				if (unit.name != "Byleth") && (unit.specialProperty != "fac") && ((settings.route != "SS") || (unit.specialProperty != "hilda")) {
					result = true
				}
			} else {
				result = true
			}
		}
	} else if settings.game == "FE11" || settings.game == "FE12" {
		// handles FE12's special cases: female-exclusive Falcoknight (class set "F") and special case dual-gender General (class set "D")
		// in FE11, male crossover is deactivated and classes always have A, B or P as a class set - the return formula was unified so it works either way
		return (unit.classSet == "F" && class.classSet != "B") || unit.classSet == class.classSet || (class.classSet != "F" && unit.classSet != "F" && settings.useMaleCrossover == Yes)
	}
	return result
}

// check if there are any compatible classes for the character
func checkForValidClasses(listOfClasses []feClass, unit feChar, settings randomizerSettings) bool {
	for i := 0; i < len(listOfClasses); i++ {
		if matchClass(listOfClasses[i], unit, settings) && (listOfClasses[i].amountLeft > 0) {
			return true
		}
	}
	return false
}

func randomizeList(inputList []feChar, amount int) []feChar {
	shufflingOrder := random.Perm(len(inputList))
	outputList := []feChar{}
	for i := 0; i < len(inputList) && i < amount; i++ {
		outputList = append(outputList, inputList[shufflingOrder[i]])
	}
	return outputList
}
