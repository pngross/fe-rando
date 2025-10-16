package main

import (
	"fmt"
	random "math/rand/v2"
	"strings"
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
	forceJagen       TernaryBool // FE11/FE12 only
}

type TernaryBool int

const (
	Unclear TernaryBool = iota
	Yes
	No
)

var supportedGames = []string{"FE11", "FE12", "FE16"}

// Check if a class matches a character
func matchClass(class feClass, unit feChar, settings randomizerSettings) bool {
	result := false
	if class.classSet == "P" {
		return strings.ReplaceAll(class.personal, unit.name, "") != class.personal
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
		// in FE11, male crossover is deactivated and classes always have A, B or P as a class set
		// the return formula was unified so it works either way
		return (unit.classSet == "F" && class.classSet != "B") || (unit.classSet == class.classSet && (class.classSet != "F" || unit.specialProperty != "no-falco")) || (class.classSet != "F" && unit.classSet != "F" && settings.useMaleCrossover == Yes)
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

func generateAvatarUnit(settings randomizerSettings) feChar {
	if settings.game == "FE16" {
		if random.IntN(2) == 0 {
			fmt.Println("Female Byleth was chosen")
			return feChar{"Byleth", "F", "", ""}
		} else {
			fmt.Println("Male Byleth was chosen")
			return feChar{"Byleth", "M", "", ""}
		}
	} else if settings.game == "FE12" {
		if random.IntN(2) == 0 {
			fmt.Println("Female Kris was chosen")
			krisClasses := []string{"Pegasus Knight->Dracoknight/Falcoknight", "Cav->Paladin", "Archer->Sniper", "Myrmidon->Swordmaster", "Mage->Sage"}
			return feChar{"Kris", "F", "", krisClasses[random.IntN(5)]}
		} else {
			fmt.Println("Male Kris was chosen")
			mKrisClasses := []string{"Knight->General", "Cav->Paladin", "Fighter->Warrior", "Archer->Sniper", "Mercenary->Hero", "Mage->Sage"}
			return feChar{"Kris", "M", "", mKrisClasses[random.IntN(6)]}
		}
	} else {
		return feChar{}
	}
}
