package ui

import "fmt"

func SortDates(displayDates []DisplayDate) [][]DisplayDate {
	var updatedDisplayDatesMatrix [][]DisplayDate
	addedDates := make(map[string]bool)

	for i := range 7 {
		var updatedDisplayDates []DisplayDate
		for _, date := range displayDates {
			if date.DayNum == i && !addedDates[date.Date] {
				updatedDisplayDates = append(updatedDisplayDates, date)
				addedDates[date.Date] = true
			}
		}
		updatedDisplayDatesMatrix = append(updatedDisplayDatesMatrix, updatedDisplayDates)
	}

	return updatedDisplayDatesMatrix
}

func Display(displayDates []DisplayDate) {
	updatedDisplayDateMatrix := SortDates(displayDates)

	for _, sl := range updatedDisplayDateMatrix {
		for _, v := range sl {
			if v.Commits <= 0 {
				fmt.Print("\033[38;5;234m◼\033[0m")
			} else if v.Commits < 3 {
				fmt.Print("\033[38;5;22m◼\033[0m")
			} else if v.Commits < 7 {
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
