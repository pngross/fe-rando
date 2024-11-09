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

func validateSettings(settings randomizerSettings) error {
	errmsg := ""
	if settings.game == "" {
		return errors.New("No FE game was named. Add the following line to the settings.txt file and specify one of the games:\ngame: <FE11|FE12|FE16>\n")
	} else {
		if settings.numberOfUnits == 0 {
			errmsg = "Number of team members is missing or zero. Make sure the settings.txt file contains the following line:\nunits: <number>\n"
		}
		if settings.game == "FE16" {
			settings.useGaidens = Unclear
			settings.useMaleCrossover = Unclear
			if settings.route == "" {
				errmsg += "FE16: You didn't specify a route. Make sure the settings.txt file contains the following line:\nroute: <CF|AM|VW|SS>\n"
			}
			if settings.forceDancer == Unclear {
				errmsg += "FE16: Settings don't specify if you want to force a dancer. Make sure the settings.txt file contains the following line:\nforce_dancer: <yes|no>\n"
			}
		} else if settings.game == "FE12" {
			settings.route = ""
			settings.useGaidens = Unclear
			if settings.useMaleCrossover == Unclear {
				errmsg += "FE12: Settings don't specify if you want male crossover class sets. Make sure the settings.txt file contains the following line:\nmale_crossover: <yes|no>\n"
			}
			if settings.forceDancer == Unclear {
				errmsg += "FE12: Settings don't specify if you want to force a dancer. Make sure the settings.txt file contains the following line:\nforce_dancer: <yes|no>\n"
			}
		} else if settings.game == "FE11" {
			settings.route = ""
			settings.useMaleCrossover = Unclear
			settings.forceDancer = Unclear
			if settings.useGaidens == Unclear {
				return errors.New("FE11: Settings don't specify if you want gaiden characters. Make sure the settings.txt file contains the following line:\ngaidens: <yes|no>\n")
			}
		}
	}
	if errmsg != "" {
		return errors.New(errmsg)
	}

	if settings.numberPerClass == 0 {
		fmt.Println("Number of units per class is missing or set to zero -> unlimited")
	}
	return nil
}
