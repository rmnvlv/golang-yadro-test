package internal

import (
	"fmt"
	"strconv"
)

func WriteData(data DataOut) {
	fmt.Println(data.timeStart.Format("15:04"))
	for _, s := range data.subjectsOut {
		l := ""
		if s.tableNumb != 0 {
			l = strconv.Itoa(s.tableNumb)
		}
		fmt.Println(s.timeOfSubj.Format("15:04"), s.id, s.name, l)
	}
	fmt.Println(data.timeEnd.Format("15:04"))

	for _, s := range data.payment {
		fmt.Println(s.table, s.money, s.workingTime)
	}
}
