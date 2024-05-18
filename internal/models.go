package internal

import "time"

type DataIn struct {
	tables      int
	costForHour int
	timeStart   time.Time
	timeEnd     time.Time
	subjects    []Subject
}

type Subject struct {
	timeOfSubj time.Time
	id         int
	name       string
	tableNumb  int
}

type DataOut struct {
	timeStart   time.Time
	timeEnd     time.Time
	payment     []Pay
	subjectsOut []Subject
}

type Pay struct {
	table       int
	money       int
	workingTime string
}

type Tables struct {
	// number int
	pay        int
	owner      string
	startOwned time.Time
	endOwned   time.Time
	timeInUsed time.Duration
}
