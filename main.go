package main

import (
	"encoding/csv"
	"fmt"
	random "math/rand/v2"
	"os"
)

type feChar struct {
	name            string
	classSet        string
	specialProperty string
	class           string
}

type feClass struct {
	name       string
	classSet   string
	personal   string
	amountLeft int
}

var fe16ClassesMale = []string{}
var fe16ClassesFemale = []string{}

var game = "FE16"
var route = "CF"
var maleCrossoverClasses = true // FE12 only

func main() {
	fmt.Println("Hello randomizer")

	amount := 15
	allChars := []feChar{}
	forcedChars := []feChar{}
	forceDancer := true
	numberPerClass := 2

	dancer := feClass{"Dancer", "N", "", 1}

	forcedChars, allChars = readAllUnits()

	if random.IntN(2) == 0 {
		forcedChars = append(forcedChars, feChar{"Byleth", "F", "", ""})
		fmt.Println("Female Byleth was chosen")
	} else {
		forcedChars = append(forcedChars, feChar{"Byleth", "M", "", ""})
		fmt.Println("Male Byleth was chosen")
	}

	outputList := randomizeList(allChars, amount-len(forcedChars))
	outputList = append(forcedChars, outputList...)

	listOfClasses := readAllClasses(numberPerClass)
	classSlotsLeft := len(listOfClasses) * numberPerClass

	randomNumber, zaehler, classIndex := 0, 0, 0

	// force dancer -> reroll for Byleth or faculty members
	if forceDancer {
		dancerFound := false
		for !dancerFound {
			randomNumber = random.IntN(len(outputList))
			if matchClass(dancer, outputList[randomNumber]) {
				outputList[randomNumber].class = "Dancer"
				dancerFound = true
			}
		}
	} else {
		listOfClasses = append(listOfClasses, dancer)
		classSlotsLeft += 1
	}

	// randomly assign classes to output list
	for i := 0; i < len(outputList); i++ {
		for outputList[i].class == "" {
			if (classSlotsLeft == 0) || !checkForValidClasses(listOfClasses, outputList[i]) {
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
				if matchClass(listOfClasses[classIndex], outputList[i]) {
					outputList[i].class = listOfClasses[classIndex].name
					if listOfClasses[classIndex].amountLeft > 0 {
						listOfClasses[classIndex].amountLeft--
					}
				}
			}
		}
	}

	if (game == "FE16") || (game == "FE14") {
		fmt.Println("Game: " + game)
		fmt.Println("Route: " + route)
	} else {
		fmt.Println("Game: " + game)
	}
	// print result to console
	for i := 0; i < len(outputList); i++ {
		fmt.Println(outputList[i].class + "!" + outputList[i].name)
	}
}

// Check if a class matches a character
func matchClass(class feClass, unit feChar) bool {
	result := false
	if game == "FE16" {
		if ((class.classSet == "N") || (class.classSet == unit.classSet)) && ((class.personal == "") || (class.personal == unit.name)) {
			if class.name == "Dancer" {
				if (unit.name != "Byleth") && (unit.specialProperty != "fac") && ((route != "SS") || (unit.specialProperty != "hilda")) {
					result = true
				}
			} else {
				result = true
			}
		}
	} else if game == "FE12" {
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
	} else if game == "FE11" {
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
func checkForValidClasses(listOfClasses []feClass, unit feChar) bool {
	for i := 0; i < len(listOfClasses); i++ {
		if matchClass(listOfClasses[i], unit) && (listOfClasses[i].amountLeft > 0) {
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

func readAllUnits() ([]feChar, []feChar) {
	forcedChars, allChars := []feChar{}, []feChar{}

	f, err := os.Open("units_" + game + ".csv")
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
			if game == "FE16" {
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

func readAllClasses(amount int) []feClass {
	f, err := os.Open("classes_" + game + ".csv")
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
