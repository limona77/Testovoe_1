package main

import (
	"Testovoe/client"
	"Testovoe/custom_errors"
	"Testovoe/data_base"
	"Testovoe/handlers"
	"Testovoe/settings_club"
	"Testovoe/table"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func ParseFile(fileName string, db *data_base.DB) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(file), "\n")

	settingsCLub := settings_club.NewSettingsClub(0, time.Time{}, time.Time{}, 0)
	for i := 0; i < len(lines); i++ {
		if lines[i] != "" {
			newLine := strings.TrimRight(lines[i], "\r")

			line := strings.Split(newLine, " ")
			// 1st three lines
			switch i {
			case 0:
				parsedCountTables, _ := strconv.Atoi(line[0])
				settingsCLub.CountTables = parsedCountTables
			case 1:

				startTime, err := time.Parse("15:04", line[0])
				if err != nil {
					log.Fatal(err)
				}
				endTime, err := time.Parse("15:04", line[1])
				if err != nil {
					log.Fatal(err)
				}

				settingsCLub.StartTime = startTime
				settingsCLub.EndTime = endTime
				settingsCLub.Duration = endTime.Sub(startTime)
				fmt.Println(startTime.Format("15:04"))
			case 2:
				parsedPricePerHour, _ := strconv.Atoi(line[0])
				settingsCLub.Price = parsedPricePerHour
			}
			// 08:48 1 client1 (example)
			if len(line) > 2 && len(line) < 4 {
				parsedTime, err := time.Parse("15:04", line[0])
				if err != nil {
					log.Fatal(fmt.Errorf("%w", errors.New("формат времени должен быть в формате ЧЧ:ММ")))
				}
				parsedActionId, err := strconv.Atoi(line[1])
				if err != nil {
					log.Fatal(fmt.Errorf("%w", errors.New("введите целое число")))
				}
				if parsedActionId < 0 || parsedActionId > 4 {
					log.Fatal(fmt.Errorf("%w", errors.New("вам доступно только 4 действия")))
				}

				parsedClientName := line[2]
				regex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
				if !regex.MatchString(parsedClientName) {
					log.Fatal(fmt.Errorf("%w", errors.New("введите валидное имя")))
				}
				t := table.NewTable(0, settingsCLub.Price, db)
				c := client.NewClient(parsedTime, parsedActionId, parsedClientName, 0, db)
				err = HandleActions(c, t, settingsCLub)
				if err != nil {
					fmt.Print(err)
				}
			}
			// 09:54 2 client1 1 (example)
			if len(line) > 3 {
				parsedTime, err := time.Parse("15:04", line[0])
				if err != nil {
					log.Fatal(fmt.Errorf("%w", errors.New("формат времени должен быть в формате ЧЧ:ММ")))
				}
				parsedActionId, err := strconv.Atoi(line[1])
				if err != nil {
					log.Fatal(fmt.Errorf("%w", errors.New("введите целое число")))
				}
				if parsedActionId < 0 || parsedActionId > 4 {
					log.Fatal(fmt.Errorf("%w", errors.New("вам доступно только 4 действия")))
				}
				parsedClientName := line[2]
				regex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
				if !regex.MatchString(parsedClientName) {
					log.Fatal(fmt.Errorf("%w", errors.New("введите валидное имя")))
				}

				parsedTableNumber, err := strconv.Atoi(line[3])
				if err != nil {
					log.Fatal(fmt.Errorf("%w", errors.New("введите целое число")))
				}

				if parsedTableNumber < 0 || parsedTableNumber > settingsCLub.CountTables {
					log.Fatal(fmt.Errorf("%w",
						errors.New("введите номер стола в диапазоне от 1 до "+
							strconv.Itoa(settingsCLub.CountTables))))
				}
				t := table.NewTable(parsedTableNumber, settingsCLub.Price, db)
				c := client.NewClient(parsedTime, parsedActionId, parsedClientName, parsedTableNumber, db)
				err = HandleActions(c, t, settingsCLub)
				if err != nil {
					fmt.Print(err)
				}
			}
		}
	}
	PrintClientsStayClub(settingsCLub, db)
	PrintSummaryPrices()
}

func PrintClientsStayClub(settingsClub settings_club.SettingsCLub, db *data_base.DB) {
	c := client.NewClient(time.Now(), 0, "", 0, db)
	t := table.NewTable(0, 0, db)
	res := c.ForEachInClientsName()
	sort.Strings(res)
	for _, v := range res {
		tableFromDB, _ := t.GetTable(v)
		c.TableNumber = tableFromDB.Id
		ConvertedTime := int(math.Ceil(settingsClub.EndTime.Sub(tableFromDB.StartTime).Hours()))
		tableFromDB.Price += settingsClub.Price*ConvertedTime - 10
		tableFromDB.Duration = settingsClub.EndTime.Sub(tableFromDB.StartTime)

		handlers.TablesWasBusy[tableFromDB.Id] = tableFromDB
		fmt.Printf("%s %d %s\n", settingsClub.EndTime.Format("15:04"), 11, v)
	}
	fmt.Printf("%s\n", settingsClub.EndTime.Format("15:04"))
}

func PrintSummaryPrices() {
	sliceOfTables := make([]int, 0)
	newTablesWasBusy := make(map[int]table.Table)
	for k, v := range handlers.TablesWasBusy {
		newTablesWasBusy[k] = v
		sliceOfTables = append(sliceOfTables, v.Id)
	}
	sort.Ints(sliceOfTables)
	for _, v := range sliceOfTables {
		formattedDuration := formatDurationToHHMM(newTablesWasBusy[v].Duration)
		fmt.Printf("%d %d %s\n", newTablesWasBusy[v].Id, newTablesWasBusy[v].Price, formattedDuration)
	}
}

func formatDurationToHHMM(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	return fmt.Sprintf("%02d:%02d", h, m)
}

func HandleActions(c *client.Client, t *table.Table, settingsClub settings_club.SettingsCLub) error {
	switch c.ActionId {
	case 1:
		err := handlers.HandleFirstAction(c, settingsClub)
		if err != nil {
			return err
		}
		return nil
	case 2:
		err := handlers.HandleSecondAction(c, t, settingsClub)
		if err != nil {
			return err
		}
		return nil
	case 3:
		err := handlers.HandleThirdAction(c, settingsClub)
		if err != nil {
			return err
		}
		return nil
	case 4:
		err := handlers.HandleFourthAction(c, t, settingsClub)
		if err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("%w",
			custom_errors.ErrActionNotExist)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Введите команду: ./task.exe <input_file>")
		return
	}
	db := data_base.NewDB()
	ParseFile(os.Args[1], db)
}
