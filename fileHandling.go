package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func SaveToFile(text string, filepath string) error {
	if !strings.HasSuffix(filepath, ".txt") {
		return errors.New("Output file must end with '.txt'")
	}
	os.Remove(filepath)
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	file.Write([]byte(text))
	return nil
}

func readDefaultSettings(game string) (randomizerSettings, error) {
	settings, err := readSettings(fmt.Sprintf("settings_%s.txt", game))

	if err != nil {
	}

	if !slices.Contains(supportedGames, game) {
		return settings, errors.New("Game " + game + " is unknown")
	} else {
		settings.game = game
	}

	return settings, err
}

func readSettings(filepath string) (randomizerSettings, error) {
	f, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer f.Close()

	settings := randomizerSettings{forceDancer: true,
		forceJagen:       true,
		useMaleCrossover: false,
		useGaidens:       false,
		numberOfUnits:    12,
		numberPerClass:   2}

	csvReader := csv.NewReader(f)
	csvReader.Comma = ':'
	inputData, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println(err.Error())
	}

	for lineNo, line := range inputData {
		if len(line) != 2 {
			return settings, fmt.Errorf("line %d in settings.txt is too long, must contain ':' exactly once", lineNo+1)
		}
		optionName := strings.Trim(line[0], " ")
		x := strings.Trim(line[1], " ")

		isYes := x == "yes"

		if optionName == "route" {
			settings.route = x
		} else if optionName == "male_crossover" {
			settings.useMaleCrossover = isYes
		} else if optionName == "gaidens" {
			settings.useGaidens = isYes
		} else if optionName == "force_dancer" {
			settings.forceDancer = isYes
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
		} else if optionName == "force_jagen" {
			settings.forceJagen = isYes
		} else {
			return settings, fmt.Errorf("Invalid option name in line %d - check the readme file for valid option names.", lineNo)
		}
	}
	return settings, nil
}
