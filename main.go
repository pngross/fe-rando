package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	random "math/rand/v2"
	"os"
	"strconv"
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
}

type TernaryBool int

const (
	Unclear TernaryBool = iota
	Yes
	No
)

var gaidens bool              // FE11 only
var maleCrossoverClasses bool // FE12 only
var route string              // FE16 only

func main() {
	fmt.Println("Hello randomizer")

	settings, err := readSettings()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = validateSettings()
	if err != nil {

	}
	if settings.game == "" {
		return
	}

	allChars := []feChar{}
	forcedChars := []feChar{}

	dancer := feClass{"Dancer", "N", "", 1}

	forcedChars, allChars = readAllUnits(settings)

	if random.IntN(2) == 0 {
		forcedChars = append(forcedChars, feChar{"Byleth", "F", "", ""})
		fmt.Println("Female Byleth was chosen")
	} else {
		forcedChars = append(forcedChars, feChar{"Byleth", "M", "", ""})
		fmt.Println("Male Byleth was chosen")
	}

	outputList := randomizeList(allChars, settings.numberOfUnits-len(forcedChars))
	outputList = append(forcedChars, outputList...)

	listOfClasses := readAllClasses(settings.numberPerClass, settings)
	classSlotsLeft := len(listOfClasses) * settings.numberPerClass

	randomNumber, zaehler, classIndex := 0, 0, 0

	// force dancer -> reroll for Byleth, Silver Snow!Hilda or faculty members
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

	if (settings.game == "FE16") || (settings.game == "FE14") {
		fmt.Println("Game: " + settings.game)
		fmt.Println("Route: " + settings.route)
	} else {
		fmt.Println("Game: " + settings.game)
	}
	// print result to console
	for i := 0; i < len(outputList); i++ {
		fmt.Println(outputList[i].className + "!" + outputList[i].name)
	}
}

func readSettings() (randomizerSettings, error) {
	f, err := os.Open("settings.txt")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.Comma = ':'
	inputData, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println(err.Error())
	}

	settings := randomizerSettings{}

	for lineNo, line := range inputData {
		if (len(line) < 0) || (len(line) > 2) {
			return settings, errors.New("line " + string(lineNo+1) + " in settings.txt is too long, must contain ':' exactly once")
		}
		optionName := strings.Trim(line[0], " ")
		x := strings.Trim(line[1], " ")

		if optionName == "game" {
			if (x != "FE11") && (x != "FE12") && (x != "FE16") {
				return settings, errors.New("Game " + x + " is unknown")
			} else {
				settings.game = x
			}
		} else if optionName == "route" {
			settings.route = x
		} else if optionName == "male_crossover" {
			if (x == "yes") || (x == "true") {
				settings.useMaleCrossover = Yes
			} else if (x == "no") || (x == "false") {
				settings.useMaleCrossover = No
			}
		} else if optionName == "gaidens" {
			if (x == "yes") || (x == "true") {
				settings.useGaidens = Yes
			} else if (x == "no") || (x == "false") {
				settings.useGaidens = No
			}
		} else if optionName == "force_dancer" {
			if (x == "yes") || (x == "true") {
				settings.forceDancer = Yes
			} else if (x == "no") || (x == "false") {
				settings.forceDancer = No
			}
		} else if optionName == "units" {
			intChecked, err := strconv.Atoi(x)
			if err != nil {
				return settings, errors.New("invalid int check, option name 'units': " + x)
			} else {
				settings.numberOfUnits = intChecked
			}
		} else if optionName == "same_class_limit" {
			intChecked, err := strconv.Atoi(x)
			if err != nil {
				return settings, errors.New("invalid int check, option name 'same_class_limit': " + x)
			} else {
				settings.numberPerClass = intChecked
			}
		} else {
			return settings, errors.New("Invalid option name in line " + string(lineNo) + " - check the readme file for valid option names.")
		}
	}
	return settings, nil
}

func validateSettings() error {
	return nil
}

// Check if a class matches a character
func matchClass(class feClass, unit feChar, settings randomizerSettings) bool {
	result := false
	if settings.game == "FE16" {
		if ((class.classSet == "N") || (class.classSet == unit.classSet)) && ((class.personal == "") || (class.personal == unit.name)) {
			if class.name == "Dancer" {
				if (unit.name != "Byleth") && (unit.specialProperty != "fac") && ((route != "SS") || (unit.specialProperty != "hilda")) {
					result = true
				}
			} else {
				result = true
			}
		}
	} else if settings.game == "FE12" {
		if class.classSet == "A" {
			if (unit.classSet == "F") || (unit.classSet == "A") || ((unit.classSet == "B") && maleCrossoverClasses) {
				result = true
			}
		} else if class.classSet == "B" {
			if (unit.classSet == "B") || ((unit.classSet == "A") && maleCrossoverClasses) {
				result = true
			}
		} else if class.classSet == "D" { // special scenario in FE12 - female General
			if (unit.classSet == "F") || (unit.classSet == "B") || ((unit.classSet == "A") && maleCrossoverClasses) {
				result = true
			}
		}
	} else if settings.game == "FE11" {
		if class.classSet == "A" {
			if (unit.classSet == "F") || (unit.classSet == "A") { // must be male with A class-set or female
				result = true
			}
		} else if (class.classSet == "B") && (unit.classSet == "B") {
			result = true
		}
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
				if ((route == "CF") && ((unitRoutes == "all") || (unitRoutes == "cf"))) || ((route == "VW") && ((unitRoutes == "all") || (unitRoutes == "notcf") || (unitRoutes == "vw"))) || ((route == "AM") && ((unitRoutes == "all") || (unitRoutes == "notcf") || (unitRoutes == "am"))) || ((route == "SS") && ((unitRoutes == "all") || (unitRoutes == "notcf"))) {
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
