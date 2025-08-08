package ui

import (
	"errors"
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

func GetRelWeekNum(curDate, startDate time.Time, startWeek, weeksInStartYear int) (int, error) {
	_, curDateWeek := curDate.ISOWeek()

	if curDateWeek > startWeek {
		curDateWeek = curDateWeek - startWeek
	} else if curDateWeek == startWeek && curDate.Year() == startDate.Year() {
		curDateWeek = 0
	} else {
		curDateWeek = (weeksInStartYear - startWeek) + curDateWeek
	}

	if curDate.Weekday().String() == "Sunday" {
		curDateWeek += 1
	}

	if curDateWeek < 0 {
		return curDateWeek, errors.New("invalid week")
	}

	return curDateWeek, nil
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

func getSundayAYearAgo(start time.Time) time.Time {
	if start.Weekday() == time.Sunday {
		return start
	}

	start = start.AddDate(0, 0, -1)
	start = getSundayAYearAgo(start)

	return start
}

func CollectDates(sortedCmits []time.Time) ([]DisplayDate, error) {
	var collectedDates []DisplayDate
	startDay := time.Now().AddDate(-1, 0, 0)
	start := getSundayAYearAgo(startDay)
	end := time.Now()

	_, oneYearAgoWeek := start.ISOWeek()
	lastYear := start.Year()

	weeksInYear, err := GetWeeksInYear(lastYear)
	if err != nil {
		return nil, err
	}

	cmitIdx := 0
	for day := start; !day.After(end); day = day.AddDate(0, 0, 1) {

		dayWeek, err := GetRelWeekNum(day, start, oneYearAgoWeek, weeksInYear)
		if err != nil {
			return nil, err
		}

		cmitCount, newCmitIdx := CountCmitsForDay(sortedCmits, day, cmitIdx)
		cmitIdx = newCmitIdx

		dateFormatted := day.Format(config.TimeFormat)
		collectedDate := New(dateFormatted, dayWeek, int(day.Weekday()), day.Weekday(), cmitCount)
		collectedDates = append(collectedDates, collectedDate)
	}

	return collectedDates, nil
}
