package ui

import (
	"time"
)

type DisplayDate struct {
	// WeekNum, DayNum, Commits int
	// Day                      string
	// Date                     time.Time

	WeekNum, DayNum, Commits int
	Day                      time.Weekday
	Date                     time.Time
}

var WeeksInYear = 52

func New(weekNum, dayNum, commits int, day time.Weekday, date time.Time) DisplayDate {
	return DisplayDate{
		Date:    date,
		WeekNum: weekNum,
		DayNum:  dayNum,
		Day:     day,
		Commits: commits,
	}
}

// TODO: break this up
func CollectDates(sortedCommits []time.Time) []DisplayDate {
	var collectedDates []DisplayDate
	start := time.Now().AddDate(-1, 0, 0)
	end := time.Now()
	_, oneYearAgoWeek := start.ISOWeek()

	commitIndex := 0
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		_, dayWeek := d.ISOWeek()

		if dayWeek > oneYearAgoWeek {
			dayWeek = dayWeek - oneYearAgoWeek
		} else if dayWeek == oneYearAgoWeek && d.Year() == start.Year() {
			dayWeek = 0
		} else {
			dayWeek = (WeeksInYear - oneYearAgoWeek) + dayWeek
		}

		commitCount := 0
		for commitIndex < len(sortedCommits) {
			commit := sortedCommits[commitIndex]
			dYear, dMonth, dDay := d.Date()
			cYear, cMonth, cDay := commit.Date()

			if dYear == cYear && dMonth == cMonth && dDay == cDay {
				commitCount++
				commitIndex++
			} else if commit.After(d) {
				break
			} else {
				commitIndex++
			}
		}

		collectedDate := New(dayWeek, int(d.Weekday()), commitCount, d.Weekday(), d)
		collectedDates = append(collectedDates, collectedDate)
	}

	return collectedDates
}

// TODO: ui stuff will go here, this code is here to learn about go doc
//
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
