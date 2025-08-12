package ui

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/double-tilde/glo/pkg/config"
)

const (
	gray       = "\033[38;5;235m"
	reset      = "\033[0m"
	daysInWeek = 7
	padding    = "  "
)

var daysInWeekLabels = []string{
	padding,
	"m ",
	padding,
	"w ",
	padding,
	"f ",
	padding,
}

// createMonthLabels uses the month of the first day of the week to see if a new
// month has started, if it has, it adds the first letter of that month to the slice
// if not it adds a space
func createMonthLabels(dates []DisplayDate) ([]string, error) {
	var months []string
	var lastMonth string

	for _, date := range dates {
		if date.DayNum == 0 {
			monthNum, err := strconv.Atoi(date.Date[5:7])
			if err != nil {
				return nil, errors.New("cannot get valid month number")
			}

			monthName := time.Month(monthNum).String()
			monthNameLetter := strings.ToLower(string(monthName[0]))

			if monthName != lastMonth {
				months = append(months, monthNameLetter)
				lastMonth = monthName
			} else {
				months = append(months, " ")
			}
		}
	}

	return months, nil
}

func getShape(cfg *config.Config) string {
	defShape := "◼"

	switch cfg.Shape {
	case "circle":
		return "●"
	case "dot":
		return "·"
	case "diamond":
		return "◆"
	case "square":
		return defShape
	default:
		return defShape
	}
}

func getColors(cfg *config.Config) (string, string, string) {
	defClr := []string{"\033[38;5;22m", "\033[38;5;34m", "\033[38;5;46m"}

	switch cfg.Color {
	case "red":
		return "\033[38;5;52m", "\033[38;5;88m", "\033[38;5;124m"
	case "blue":
		return "\033[38;5;17m", "\033[38;5;19m", "\033[38;5;21m"
	case "green":
		return defClr[0], defClr[1], defClr[2]
	default:
		return defClr[0], defClr[1], defClr[2]
	}
}

// createCmitMatrix goes through all of the dates in the last year and sorts them
// into a matrix that starts on sunday and goes through each day of the week
// ending with 7 slices and a int that represents the most commits
func createCmitMatrix(displayDates []DisplayDate) ([][]DisplayDate, int) {
	var cmitMatrix [][]DisplayDate
	var mostCommits int
	addedDates := make(map[string]bool)

	for day := range daysInWeek {
		var updatedCmits []DisplayDate
		for _, date := range displayDates {
			if date.DayNum == day && !addedDates[date.Date] {
				addedDates[date.Date] = true

				updatedCmits = append(updatedCmits, date)

				if date.Commits > mostCommits {
					mostCommits = date.Commits
				}
			}
		}
		cmitMatrix = append(cmitMatrix, updatedCmits)
	}

	return cmitMatrix, mostCommits
}

func DisplayYear(cfg *config.Config, displayDates []DisplayDate) error {
	// create the commit matrix starting with all the sundays, all mondays, etc
	// also get the highest amount of commits for the last year
	cmitMatrix, mostCmits := createCmitMatrix(displayDates)

	// calculate the 1/3rd and 2/3rd amounts based on the most commits
	third := int(math.Floor(float64(mostCmits) / 100 * 33))
	twoThirds := int(math.Floor(float64(mostCmits) / 100 * 66))

	dark, medium, light := getColors(cfg)
	shape := getShape(cfg)

	monthLabels, err := createMonthLabels(displayDates)
	if err != nil {
		return err
	}

	// display the month labels first along the top
	fmt.Print(padding)
	for _, month := range monthLabels {
		fmt.Print(month)
	}
	fmt.Println()

	// cycle through each day starting on sunday
	for dayNum, dates := range cmitMatrix {
		// display the label for the day we are on first
		fmt.Print(daysInWeekLabels[dayNum])

		// then go through all of the dates that match that day and display them
		for _, date := range dates {
			if date.Commits <= 0 {
				fmt.Print(gray + shape + reset)
			} else if date.Commits < third {
				fmt.Print(dark + shape + reset)
			} else if date.Commits < twoThirds {
				fmt.Print(medium + shape + reset)
			} else {
				fmt.Print(light + shape + reset)
			}
		}
		fmt.Println()
	}

	return nil
}
