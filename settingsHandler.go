package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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
