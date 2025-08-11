package ui

import (
	"fmt"
	"math"
)

const daysInWeek = 7

func SortDates(displayDates []DisplayDate) ([][]DisplayDate, int) {
	var updatedDisplayDatesMatrix [][]DisplayDate
	addedDates := make(map[string]bool)
	mostCommits := 0

	for i := range daysInWeek {
		var updatedDisplayDates []DisplayDate
		for _, date := range displayDates {
			if date.DayNum == i && !addedDates[date.Date] {
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

func Display(displayDates []DisplayDate) {
	updatedDisplayDateMatrix, mostCommits := SortDates(displayDates)

	third := int(math.Floor(float64(mostCommits) / 100 * 33))
	twoThirds := int(math.Floor(float64(mostCommits) / 100 * 66))

	fmt.Println(mostCommits, third, twoThirds)

	for _, displayDateDay := range updatedDisplayDateMatrix {
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
