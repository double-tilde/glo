package ui

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
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

func SortDates(displayDates []DisplayDate) ([][]DisplayDate, int, map[string]int) {
	var updatedDisplayDatesMatrix [][]DisplayDate
	addedDates := make(map[string]bool)
	mostCommits := 0
	months := make(map[string]int)

	for day := range daysInWeek {
		var updatedDisplayDates []DisplayDate
		for _, date := range displayDates {
			if date.DayNum == day && !addedDates[date.Date] {
				addedDates[date.Date] = true

				if _, ok := months[date.Date[0:7]]; ok {
					months[date.Date[0:7]] += 1
				} else {
					months[date.Date[0:7]] = 0
				}

				updatedDisplayDates = append(updatedDisplayDates, date)

				if date.Commits > mostCommits {
					mostCommits = date.Commits
				}
			}
		}
		updatedDisplayDatesMatrix = append(updatedDisplayDatesMatrix, updatedDisplayDates)
	}

	return updatedDisplayDatesMatrix, mostCommits, months
}

func Display(displayDates []DisplayDate) error {
	updatedDisplayDateMatrix, mostCommits, months := SortDates(displayDates)

	result := []string{}

	keys := make([]string, 0, len(months))
	for k := range months {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		monthNum, err := strconv.Atoi(k[5:])
		if err != nil {
			return errors.New("cannot get valid month number")
		}

		monthName := time.Month(monthNum).String()
		monthNameLetter := strings.ToLower(string(monthName[0]))
		weeks := math.Round(float64(months[k]) / 7)
		result = append(result, monthNameLetter)

		for i := 1; i < int(weeks); i++ {
			result = append(result, " ")
		}
	}

	third := int(math.Floor(float64(mostCommits) / 100 * 33))
	twoThirds := int(math.Floor(float64(mostCommits) / 100 * 66))

	fmt.Print("  ")
	for _, month := range result {
		fmt.Print(month)
	}
	fmt.Println()

	for weekDayNumber, displayDateDay := range updatedDisplayDateMatrix {
		fmt.Print(daysInWeekLabels[weekDayNumber])
		for _, day := range displayDateDay {
			if day.Commits <= 0 {
				fmt.Print("\033[38;5;236m◼\033[0m")
			} else if day.Commits < third {
				fmt.Print("\033[38;5;22m◼\033[0m")
			} else if day.Commits < twoThirds {
				fmt.Print("\033[38;5;34m◼\033[0m")
			} else {
				fmt.Print("\033[38;5;46m◼\033[0m")
			}
		}
		fmt.Println()
	}

	return nil
}

// type Color int
//
// const (
// 	Red   Color = iota // Red color code
// 	Green              // Green color code
// 	Blue               // Blue color code
// )
//
// func PrintColor(text string, color Color) error {
// 	switch color {
// 	case Red:
// 		fmt.Print("\033[31m" + text + "\033[0m")
// 	case Green:
// 		fmt.Print("\033[32m" + text + "\033[0m")
// 	case Blue:
// 		fmt.Print("\033[34m" + text + "\033[0m")
// 	default:
// 		return errors.New("unsupported color or printing failure")
// 	}
// 	return nil
// }
//
// func Printing() {
// 	// Green backgrounds
// 	fmt.Println("\033[48;5;22m \033[0m")
// 	fmt.Println("\033[48;5;34m \033[0m")
// 	fmt.Println("\033[48;5;46m \033[0m")
//
// 	// Blue backgrounds
// 	fmt.Println("\033[48;5;18m \033[0m")
// 	fmt.Println("\033[48;5;19m \033[0m")
// 	fmt.Println("\033[48;5;21m \033[0m")
//
// 	// Red backgrounds
// 	fmt.Println("\033[48;5;52m \033[0m")
// 	fmt.Println("\033[48;5;88m \033[0m")
// 	fmt.Println("\033[48;5;124m \033[0m")
// }
