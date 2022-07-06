package main

import "time"

func initMonth() [][7]time.Time {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, now.Location())

	weekStart := firstOfMonth
	if weekStart.Weekday() != time.Sunday {
		for ; weekStart.Weekday() != time.Sunday; weekStart = weekStart.AddDate(0, 0, -1) {
		}
	}

	lastDayOfMonth := firstOfMonth.AddDate(0, 1, -1)
	lastDayofLastWeek := lastDayOfMonth
	if lastDayofLastWeek.Weekday() != time.Sunday { // until Sunday so, in the next block we can use the condition Before
		for ; lastDayofLastWeek.Weekday() != time.Sunday; lastDayofLastWeek = lastDayofLastWeek.AddDate(0, 0, 1) {
		}
	}

	datesOfMonth := make([][7]time.Time, 5)

	row := 0
	currentDay := weekStart
	for currentDay.Before(lastDayofLastWeek) {

		if currentDay.Weekday() == time.Sunday && !currentDay.Equal(weekStart) {
			row++
		}

		datesOfMonth[row][currentDay.Weekday()] = currentDay

		currentDay = currentDay.AddDate(0, 0, 1)
		continue
	}

	return datesOfMonth
}
