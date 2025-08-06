package ui

import (
	"errors"
	"fmt"
	"time"

	"github.com/double-tilde/glo/pkg/config"
)

type DisplayDate struct {
	Date    string
	WeekNum int
	DayNum  int
	Day     time.Weekday
	Commits int
}

func New(date string, weekNum, dayNum int, day time.Weekday, cmits int) DisplayDate {
	return DisplayDate{
		Date:    date,
		WeekNum: weekNum,
		DayNum:  dayNum,
		Day:     day,
		Commits: cmits,
	}
}

func GetWeeksInYear(lastYear int) (int, error) {
	if lastYear < 0 {
		return 0, errors.New("invalid year")
	}

	lastDay := time.Date(lastYear, time.December, 28, 0, 0, 0, 0, time.UTC)
	_, lastYearWeeks := lastDay.ISOWeek()

	return lastYearWeeks, nil
}

func CalculateWeekNumber(day, start time.Time, oneYearAgoWeek, weeksInYear int) int {
	_, dayWeek := day.ISOWeek()

	if dayWeek > oneYearAgoWeek {
		dayWeek = dayWeek - oneYearAgoWeek
	} else if dayWeek == oneYearAgoWeek && day.Year() == start.Year() {
		dayWeek = 0
	} else {
		dayWeek = (weeksInYear - oneYearAgoWeek) + dayWeek
	}

	return dayWeek
}

func CountCmitsForDay(sortedCmits []time.Time, day time.Time, cmitIdx int) (int, int) {
	cmitCount := 0
	for cmitIdx < len(sortedCmits) {
		cmit := sortedCmits[cmitIdx]
		dYear, dMonth, dDay := day.Date()
		cYear, cMonth, cDay := cmit.Date()

		if dYear == cYear && dMonth == cMonth && dDay == cDay {
			cmitCount++
			cmitIdx++
		} else if cmit.After(day) {
			break
		} else {
			cmitIdx++
		}
	}

	return cmitCount, cmitIdx
}

func CollectDates(sortedCmits []time.Time) ([]DisplayDate, error) {
	var collectedDates []DisplayDate
	start := time.Now().AddDate(-1, 0, 0)
	end := time.Now()

	_, oneYearAgoWeek := start.ISOWeek()
	lastYear := start.Year()
	fmt.Println(lastYear)
	weeksInYear, err := GetWeeksInYear(lastYear)
	if err != nil {
		return nil, err
	}

	cmitIdx := 0
	for day := start; !day.After(end); day = day.AddDate(0, 0, 1) {

		dayWeek := CalculateWeekNumber(day, start, oneYearAgoWeek, weeksInYear)
		cmitCount, newCmitIdx := CountCmitsForDay(sortedCmits, day, cmitIdx)
		cmitIdx = newCmitIdx

		dateFormatted := day.Format(config.TimeFormat)
		collectedDate := New(dateFormatted, dayWeek, int(day.Weekday()), day.Weekday(), cmitCount)
		collectedDates = append(collectedDates, collectedDate)
	}

	return collectedDates, nil
}
