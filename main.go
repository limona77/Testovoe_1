package main

import (
	"Testovoe/client"
	"Testovoe/custom_errors"
	"Testovoe/data_base"
	"Testovoe/settings_club"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/slog"
)

func ParseFile(fileName string, db *data_base.DB) {
	file, err := os.ReadFile("./test_file.txt")
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(file), "\n")

	settingsCLub := settings_club.NewSettingsClub(0, time.Time{}, time.Time{}, 0)
	for i := 0; i < 13; i++ {
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
				fmt.Println(startTime.Format("15:04"))
			case 2:
				parsedPricePerHour, _ := strconv.Atoi(line[0])
				settingsCLub.PricePerHour = parsedPricePerHour
			}
			// 08:48 1 client1 (example)
			if len(line) > 2 && len(line) < 4 {
				parsedTime, _ := time.Parse("15:04", line[0])
				parsedActionId, _ := strconv.Atoi(line[1])
				parsedClientName := line[2]

				c := client.NewClient(parsedTime, parsedActionId, parsedClientName, 0, db)
				err := HandleAction(c, settingsCLub)
				if err != nil {
					fmt.Print(err)
				}
			}
			// 09:54 2 client1 1 (example)
			if len(line) > 3 {
				parsedTime, _ := time.Parse("15:04", line[0])
				parsedActionId, _ := strconv.Atoi(line[1])
				parsedClientName := line[2]
				parsedTableNumber, _ := strconv.Atoi(line[3])

				c := client.NewClient(parsedTime, parsedActionId, parsedClientName, parsedTableNumber, db)
				err := HandleAction(c, settingsCLub)
				if err != nil {
					fmt.Print(err)
				}
			}
		}
	}
}

func HandleAction(c *client.Client, settingsClub settings_club.SettingsCLub) error {
	switch c.ActionId {
	case 1:
		fmt.Printf("%s %v %s\n", c.CurrentTime.Format("15:04"), c.ActionId, c.ClientName)
		ok := c.CheckValidTime(settingsClub.StartTime, settingsClub.EndTime)
		if !ok {
			return fmt.Errorf(
				"%s %v %s\n",
				c.CurrentTime.Format("15:04"),
				custom_errors.ErrCode,
				custom_errors.ErrNotOpenYet)
		}
	case 2:
		busy := c.CheckInDB(c.ActionId)
		if busy {
			fmt.Printf("%s %v %s %d\n",
				c.CurrentTime.Format("15:04"),
				c.ActionId, c.ClientName, c.TableNumber)
			return fmt.Errorf("%s %v %s\n",
				c.CurrentTime.Format("15:04"),
				custom_errors.ErrCode,
				custom_errors.ErrPlaceIsBusy)
		}
		c.SetInDB(c.ActionId, c.TableNumber)
		settingsClub.CountTables--
		fmt.Printf("%s %v %s %d\n",
			c.CurrentTime.Format("15:04"),
			c.ActionId, c.ClientName, c.TableNumber)
		return nil
	case 3:
		fmt.Printf("%s %v %s\n",
			c.CurrentTime.Format("15:04"),
			c.ActionId, c.ClientName)
		if settingsClub.CountTables == 0 {
			c.DeleteInDB(c.ActionId)
			settingsClub.CountTables++
		}
		fmt.Printf("%s %v %s\n",
			c.CurrentTime.Format("15:04"),
			custom_errors.ErrCode,
			custom_errors.ErrICanWaitNoLonger)
		return nil
	case 4:
		exist := c.CheckInDB(c.ActionId)
		if !exist {
			return fmt.Errorf("%s %v %s\n",
				c.CurrentTime.Format("15:04"),
				custom_errors.ErrCode,
				custom_errors.ErrClientUnknown)
		}
		return nil
	default:
		return fmt.Errorf("%w",
			custom_errors.ErrActionNotExist)
	}
	return nil
}

func main() {
	db := data_base.NewDB()
	ParseFile("test_file.txt", db)
	slog.Info("ParseFile - OK")
}
