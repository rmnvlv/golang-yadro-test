package internal

import (
	"container/list"
	"regexp"
	"sort"
)

type StackSubj[T any] struct {
	data []T
}

func (s *StackSubj[T]) Pop() T {
	if len(s.data) == 0 {
		panic("stack is empty")
	}

	last := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]

	return last
}

func (s *StackSubj[T]) Show() T {
	if len(s.data) == 0 {
		panic("stack is empty")
	}

	last := s.data[len(s.data)-1]

	return last
}

func (s *StackSubj[T]) Push(v T) {
	s.data = append(s.data, v)
}

func (s *StackSubj[T]) IsEmpty() bool {
	return len(s.data) == 0
}

func ParseData(data DataIn) (DataOut, error) {
	dataOut := DataOut{}
	var err error
	dataOut.timeStart = data.timeStart
	dataOut.timeEnd = data.timeEnd

	stackInSubjects := StackSubj[Subject]{}

	// Fill stack of our in subjects
	// fmt.Println("Fill the stack -----------------------")
	for i := len(data.subjects) - 1; i >= 0; i-- {
		stackInSubjects.Push(data.subjects[i])
	}

	// fmt.Println(stackInSubjects, "\nMy stack -------------------------")
	// fmt.Println(len(stackInSubjects.data))
	// tables := map[int]bool{}
	tables := []Tables{}
	for i := 0; i < data.tables; i++ {
		tables = append(tables, Tables{
			owner: "none",
			// number: i + 1,
			// timeInUsed: 0,
			// pay:        10,
		})

	}

	clientsQueue := list.New()

	// tables := map[interface{}]interface{}{}
	clintsOnTables := map[string]int{}
	clientsInClub := map[string]interface{}{}

	// clients := map[string]bool{}

	//Hardcore check all subjects
	for {
		if stackInSubjects.IsEmpty() {
			break
		}

		subjectIn := stackInSubjects.Pop()
		subjectOut := Subject{}

		if data.timeEnd.Before(subjectIn.timeOfSubj) {
			break
		}

		dataOut.subjectsOut = append(dataOut.subjectsOut, subjectIn)
		//Generate errors for clients
		errosOfClients := []string{"NotOpenYet", "YouShallNotPass", "ClientUnknown", "ICanWaitNoLonger!", "PleaseIsBusy"}

		//Not open yet error |||| it could be bettr \-_-/
		if subjectIn.timeOfSubj.Before(data.timeStart) || subjectIn.timeOfSubj.After(data.timeEnd) {
			subjectOut.timeOfSubj = subjectIn.timeOfSubj
			subjectOut.id = 13
			subjectOut.name = errosOfClients[0]
			dataOut.subjectsOut = append(dataOut.subjectsOut, subjectOut)
			continue
		}

		tableClear := false
		for i := 0; i < data.tables; i++ {
			if tables[i].owner == "none" {
				tableClear = true
			}
		}

		switch subjectIn.id {
		//ID1: Клиент пришел
		case 1:
			if clientsInClub[subjectIn.name] != nil {
				subjectOut.timeOfSubj = subjectIn.timeOfSubj
				subjectOut.id = 13
				subjectOut.name = errosOfClients[1]
				dataOut.subjectsOut = append(dataOut.subjectsOut, subjectOut)
				continue
			} else {
				clientsInClub[subjectIn.name] = true
				// clintsOnTables[subjectIn.name] = nil
			}
		//ID2: Клиент сел за стол
		case 2:
			//Неизвестный клиент
			if clientsInClub[subjectIn.name] == nil {
				subjectOut.timeOfSubj = subjectIn.timeOfSubj
				subjectOut.id = 13
				subjectOut.name = errosOfClients[2]
				dataOut.subjectsOut = append(dataOut.subjectsOut, subjectOut)
				continue
			}

			//Пытается сесть за занятый стол
			if tables[subjectIn.tableNumb-1].owner != "none" {
				subjectOut.timeOfSubj = subjectIn.timeOfSubj
				subjectOut.id = 13
				subjectOut.name = errosOfClients[4]
				dataOut.subjectsOut = append(dataOut.subjectsOut, subjectOut)
				continue
			} else {
				//Освобождаем стол от клиента если он пересаживается и оплачивание стола
				if clintsOnTables[subjectIn.name] != 0 {
					tables[clintsOnTables[subjectIn.name]-1].owner = "none"
					owned := subjectIn.timeOfSubj.Sub(tables[clintsOnTables[subjectIn.name]-1].startOwned)
					// fmt.Println("в обращении", owned)
					tables[clintsOnTables[subjectIn.name]-1].timeInUsed += owned
					needToPay := 0
					if owned.Minutes() > 0 {
						hours := int(owned.Hours() + 1)
						needToPay = hours * data.costForHour
					} else {
						hours := int(owned.Hours())
						needToPay = hours * data.costForHour
					}
					tables[clintsOnTables[subjectIn.name]-1].pay += needToPay
					tables[clintsOnTables[subjectIn.name]-1].endOwned = subjectIn.timeOfSubj
				}
				// Сажаем за стол
				tables[subjectIn.tableNumb-1].owner = subjectIn.name
				tables[subjectIn.tableNumb-1].startOwned = subjectIn.timeOfSubj
				clintsOnTables[subjectIn.name] = subjectIn.tableNumb
			}
		//ID3: Клиент ожидает
		case 3:
			//Есть свободные столы а клиент ждет
			if tableClear {
				subjectOut.timeOfSubj = subjectIn.timeOfSubj
				subjectOut.id = 13
				subjectOut.name = errosOfClients[3]
				dataOut.subjectsOut = append(dataOut.subjectsOut, subjectOut)
				continue
				//Очередь клиентов больше числа столов
			} else if clientsQueue.Len() > data.tables {
				subjectOut.timeOfSubj = subjectIn.timeOfSubj
				subjectOut.id = 11
				subjectOut.name = subjectIn.name
				dataOut.subjectsOut = append(dataOut.subjectsOut, subjectOut)
				continue
				//Встал в очередь в ожидании компа
			} else if clientsQueue.Len() < data.tables && !tableClear {
				clientsQueue.PushBack(subjectIn.name)
				// fmt.Printf("%s stand in queue in time %s \n", subjectIn.name, subjectIn.timeOfSubj.Format("15:04"))
			}
		//ID4: Клиент ушел
		case 4:
			// Такого клиента нет
			if clientsInClub[subjectIn.name] == nil {
				subjectOut.timeOfSubj = subjectIn.timeOfSubj
				subjectOut.id = 13
				subjectOut.name = errosOfClients[2]
				dataOut.subjectsOut = append(dataOut.subjectsOut, subjectOut)
				continue
			} else {
				clientsInClub[subjectIn.name] = nil

				// Сажаем за пк чела из очереди если она есть
				if clintsOnTables[subjectIn.name] != 0 {
					//Собираем дань!!!
					owned := subjectIn.timeOfSubj.Sub(tables[clintsOnTables[subjectIn.name]-1].startOwned)
					// fmt.Println("Owned:", owned, subjectIn.timeOfSubj.Format("15:04"), tables[clintsOnTables[subjectIn.name]-1].startOwned.Format("15:04"))
					tables[clintsOnTables[subjectIn.name]-1].timeInUsed += owned
					needToPay := 0
					if owned.Minutes() > 0 {
						hours := int(owned.Hours() + 1)
						// fmt.Println(hours, owned)
						needToPay = hours * data.costForHour
					} else {
						hours := int(owned.Hours())
						// fmt.Println(hours)
						needToPay = hours * data.costForHour
					}
					// fmt.Printf("Клиент %s уходит и должен заплатить %v\n", subjectIn.name, needToPay)
					tables[clintsOnTables[subjectIn.name]-1].pay += needToPay
					tables[clintsOnTables[subjectIn.name]-1].endOwned = subjectIn.timeOfSubj
					if clientsQueue.Len() > 0 {
						client := clientsQueue.Front().Value
						clientsQueue.Remove(clientsQueue.Front())
						var name string

						if str, ok := client.(string); ok {
							name = str
						}
						// fmt.Println(name)
						subjectOut.timeOfSubj = subjectIn.timeOfSubj
						subjectOut.id = 12
						subjectOut.name = name
						subjectOut.tableNumb = clintsOnTables[subjectIn.name]
						dataOut.subjectsOut = append(dataOut.subjectsOut, subjectOut)
						clintsOnTables[name] = clintsOnTables[subjectIn.name]
						tables[clintsOnTables[subjectIn.name]-1].owner = name
						tables[clintsOnTables[subjectIn.name]-1].startOwned = subjectIn.timeOfSubj
						continue
					}
					tables[clintsOnTables[subjectIn.name]-1].owner = "none"
					clintsOnTables[subjectIn.name] = 0
					continue
				}
			}
		}
	}

	for i := range tables {
		if tables[i].owner != "none" {
			owned := data.timeEnd.Sub(tables[i].startOwned)
			tables[i].timeInUsed += owned
			needToPay := 0
			if owned.Minutes() > 0 {
				hours := int(owned.Hours() + 1)
				needToPay = hours * data.costForHour
			} else {
				hours := int(owned.Hours())
				needToPay = hours * data.costForHour
			}
			tables[i].pay += needToPay
		}
		re := regexp.MustCompile("[0-9]+")
		mainTime := re.FindAllString(tables[i].timeInUsed.String(), -1)
		// fmt.Println(mainTime)
		if len(mainTime) > 1 {
			if len(mainTime[1]) < 2 {
				mainTime[1] = "0" + mainTime[1]
			}
		}

		workingTime := "00:00"
		if len(mainTime) > 1 {
			workingTime = mainTime[0] + ":" + mainTime[1]
		}

		dataOut.payment = append(dataOut.payment, Pay{
			table:       i + 1,
			money:       tables[i].pay,
			workingTime: workingTime,
		})
	}

	sort.SliceStable(tables, func(i, j int) bool {
		return tables[i].owner < tables[j].owner
	})

	for _, v := range tables {
		if v.owner != "none" {
			dataOut.subjectsOut = append(dataOut.subjectsOut, Subject{
				timeOfSubj: data.timeEnd,
				id:         11,
				name:       v.owner,
			})
		}
	}

	return dataOut, err
}
