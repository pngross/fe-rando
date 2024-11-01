package main

import (
	"encoding/csv"
	"fmt"
	random "math/rand/v2"
	"os"
)

type feChar struct {
	name            string
	gender          string
	specialProperty string
	class           string
}

type feClass struct {
	name       string
	gender     string
	personal   string
	amountLeft int
}

var fe16ClassesMale = []string{}
var fe16ClassesFemale = []string{}

var game = "FE16"
var route = "CF"

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

	// dancer forcieren -> reroll bei Byleth oder Fakultätsmitgliedern.
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

	// Klassen randomisiert den Charakteren zuweisen
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
	// Ergebnis ausgeben
	for i := 0; i < len(outputList); i++ {
		fmt.Println(outputList[i].class + "!" + outputList[i].name)
	}
}

// Prüft, ob eine Klasse einem Charakter zugewiesen werden kann - FE16-Version, erweiterbar für FE11/12
func matchClass(class feClass, unit feChar) bool {
	result := false
	if ((class.gender == "N") || (class.gender == unit.gender)) && ((class.personal == "") || (class.personal == unit.name)) {
		if (game == "FE16") && (class.name == "Dancer") {
			if (unit.name != "Byleth") && (unit.specialProperty != "fac") && ((route != "SS") || (unit.specialProperty != "hilda")) {
				result = true
			}
		} else {
			result = true
		}
	}

	return result
}

// gibt es kompatible Klassen mit >0 Restanzahl
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
					currentChar.gender = field
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
					currentClass.gender = field
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
