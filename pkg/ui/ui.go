package ui

import (
	"fmt"
	"math"

	"github.com/double-tilde/glo/pkg/config"
)

const (
	Gray  = "\033[38;5;236m"
	Reset = "\033[0m"
)

var (
	block            = "◼"
	daysInWeekLabels = []string{
		"  ",
		"m ",
		"  ",
		"w ",
		"  ",
		"f ",
		"  ",
	}
	daysInWeek = 7
)

func setColor(cfg *config.Config) (string, string, string) {
	if cfg.Color == "red" {
		return "\033[38;5;52m", "\033[38;5;88m", "\033[38;5;124m"
	}

	if cfg.Color == "blue" {
		return "\033[38;5;17m", "\033[38;5;19m", "\033[38;5;21m"
	}

	// default: green
	return "\033[38;5;22m", "\033[38;5;34m", "\033[38;5;46m"
}

func SortDates(displayDates []DisplayDate) ([][]DisplayDate, int) {
	var updatedDisplayDatesMatrix [][]DisplayDate
	addedDates := make(map[string]bool)
	mostCommits := 0

	for day := range daysInWeek {
		var updatedDisplayDates []DisplayDate
		for _, date := range displayDates {
			if date.DayNum == day && !addedDates[date.Date] {
				addedDates[date.Date] = true

				updatedDisplayDates = append(updatedDisplayDates, date)

				if date.Commits > mostCommits {
					mostCommits = date.Commits
				}
			}
		}
		updatedDisplayDatesMatrix = append(updatedDisplayDatesMatrix, updatedDisplayDates)
	}

	return updatedDisplayDatesMatrix, mostCommits
}

func Display(cfg *config.Config, displayDates []DisplayDate, monthLabels []string) error {
	updatedDisplayDateMatrix, mostCommits := SortDates(displayDates)
	third := int(math.Floor(float64(mostCommits) / 100 * 33))
	twoThirds := int(math.Floor(float64(mostCommits) / 100 * 66))
	dark, medium, light := setColor(cfg)

	fmt.Print("  ")
	for _, month := range monthLabels {
		fmt.Print(month)
	}
	fmt.Println()

	for weekDayNumber, displayDateDay := range updatedDisplayDateMatrix {
		fmt.Print(daysInWeekLabels[weekDayNumber])
		for _, day := range displayDateDay {
			if day.Commits <= 0 {
				fmt.Print(Gray + block + Reset)
			} else if day.Commits < third {
				fmt.Print(dark + block + Reset)
			} else if day.Commits < twoThirds {
				fmt.Print(medium + block + Reset)
			} else {
				fmt.Print(light + block + Reset)
			}
		}
		fmt.Println()
	}

	return nil
}
