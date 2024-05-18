package internal

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func FormatData(data []string) (DataIn, error) {
	formattingData := DataIn{}
	var err error

	//Write tables
	formattingData.tables, err = strconv.Atoi(data[0])
	if err != nil || formattingData.tables < 0 {
		return DataIn{}, fmt.Errorf(fmt.Sprintf("Error with tables in line '%s': %s", data[0], err))
	}

	//Check for two times
	times := strings.Split(data[1], " ")
	if len(times) != 2 {
		return DataIn{}, fmt.Errorf(fmt.Sprintf("Error with start and end times in line '%s': %s", data[1], err))
	}

	//Kostil dlya sravnenia
	var timeStart, timeEnd time.Time

	for i := 0; i < 2; i++ {
		t, err := time.Parse("15:04", times[i])
		if err != nil {
			return DataIn{}, fmt.Errorf(fmt.Sprintf("Error with times in line '%s': %s", data[1], err))
		}
		if i == 0 {
			timeStart = t
		} else if i == 1 {
			timeEnd = t
		}
	}

	//Chek start time before time end
	if timeStart.After(timeEnd) {
		return DataIn{}, fmt.Errorf(fmt.Sprintf("Time start after time end in line %s: %s", data[1], err))
	}

	// Write end and start times
	formattingData.timeStart = timeStart //Format("15:04")
	formattingData.timeEnd = timeEnd     //Format("15:04")

	//Write cost for table
	formattingData.costForHour, err = strconv.Atoi(data[2])
	if err != nil || formattingData.costForHour < 0 {
		return DataIn{}, fmt.Errorf(fmt.Sprintf("Error with cost for table in line '%s': %s", data[2], err))
	}

	r, _ := regexp.Compile("^[a-z0-9-_]+$")

	//Check and write subjects
	for i := 3; i < len(data); i++ {
		subjectString := strings.Split(data[i], " ")

		//Check len of line data
		if len(subjectString) < 3 || len(subjectString) > 4 {
			return DataIn{}, fmt.Errorf(fmt.Sprintf("Too much or not enought subjects in line: %s", data[i]))
		}

		subject := Subject{}

		//Check and Write time of subj
		subject.timeOfSubj, err = time.Parse("15:04", subjectString[0])
		if err != nil || subject.timeOfSubj.After(timeEnd) {
			return DataIn{}, fmt.Errorf(fmt.Sprintf("Error in subject with time in line '%s': %s", data[i], err))
		} else if i != 3 && formattingData.subjects[i-4].timeOfSubj.After(subject.timeOfSubj) {
			return DataIn{}, fmt.Errorf(fmt.Sprintf("Bad time in lines: '%s' after '%s'", data[i-1], data[i]))
		}

		//Check and Write subj id
		subject.id, err = strconv.Atoi(subjectString[1])
		if err != nil {
			return DataIn{}, fmt.Errorf(fmt.Sprintf("Error with subject id in line '%s': %s", data[i], err))
		}

		//Check  0<id<5
		if subject.id > 4 || subject.id < 1 {
			return DataIn{}, fmt.Errorf(fmt.Sprintf("Unknown id (id[1,2,3,4]) in line '%s'", data[i]))
		}

		//Check and Write name of client
		subject.name = subjectString[2]
		matched := r.MatchString(subject.name)
		if !matched {
			return DataIn{}, fmt.Errorf(fmt.Sprintf("Bad name in line '%s'", data[i]))
		}

		//Check table number
		if subject.id == 2 {
			subject.tableNumb, err = strconv.Atoi(subjectString[3])
			if err != nil || subject.tableNumb > formattingData.tables {
				return DataIn{}, fmt.Errorf(fmt.Sprintf("Bad table number in line '%s': %s", data[i], err))
			}
		}

		//Add subject to array
		formattingData.subjects = append(formattingData.subjects, subject)
	}
	return formattingData, err
}
