package main

import (
	"Testovoe/client"
	"Testovoe/custom_errors"
	"Testovoe/data_base"
	"Testovoe/queue"
	"Testovoe/settings_club"
	"Testovoe/table"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/slog"
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
				parsedTime, _ := time.Parse("15:04", line[0])
				parsedActionId, _ := strconv.Atoi(line[1])
				parsedClientName := line[2]
				t := table.NewTable(0, settingsCLub.Price, db)
				c := client.NewClient(parsedTime, parsedActionId, parsedClientName, 0, db)
				err := HandleAction(c, t, settingsCLub)
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

				t := table.NewTable(parsedTableNumber, settingsCLub.Price, db)
				c := client.NewClient(parsedTime, parsedActionId, parsedClientName, parsedTableNumber, db)
				err := HandleAction(c, t, settingsCLub)
				if err != nil {
					fmt.Print(err)
				}
			}
		}
	}
	PrintClientsStayClub(settingsCLub, db)
	PrintSummaryPrices(settingsCLub, db)
}

func PrintClientsStayClub(settingsClub settings_club.SettingsCLub, db *data_base.DB) {
	c := client.NewClient(time.Now(), 0, "", 0, db)
	t := table.NewTable(0, 0, db)
	res := c.ForEachInClientsName()

	sort.Strings(res)
	for _, v := range res {
		tableFromDB, _ := t.GetTable(v)
		c.TableNumber = tableFromDB.Id
		tableFromDB.Price += settingsClub.Price*int(math.Ceil(settingsClub.EndTime.Sub(tableFromDB.StartTime).Hours())) - 10
		tableFromDB.Duration = settingsClub.EndTime.Sub(tableFromDB.StartTime)
		// formattedDuration := formatDurationToHHMM(settingsClub.EndTime.Sub(tableFromDB.StartTime))
		TablesWasBusy[tableFromDB.Id] = tableFromDB
		fmt.Printf("%s %d %s\n", settingsClub.EndTime.Format("15:04"), 11, v)
		fmt.Println()
	}
	fmt.Printf("%s\n", settingsClub.EndTime.Format("15:04"))
}

func PrintSummaryPrices(settingsClub settings_club.SettingsCLub, db *data_base.DB) {
	fmt.Println(TablesWasBusy)
	for _, v := range TablesWasBusy {
		fmt.Println(v.StartTime, v.EndTime)
		if v.EndTime.Hour() != 0 {
			v.Duration = v.EndTime.Sub(v.StartTime)
		}
		formattedDuration := formatDurationToHHMM(v.Duration)
		fmt.Printf("%d %d %s\n", v.Id, v.Price, formattedDuration)
	}
	//		t := table.NewTable(0, 0, db)
	//		res := t.ForEachTables()
	//
	//		sort.Ints(res)
	//		for _, v := range res {
	//			t.
	//		}
	//		fmt.Printf("%s\n", settingsClub.EndTime.Format("15:04"))
}

func formatDurationToHHMM(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	return fmt.Sprintf("%02d:%02d", h, m)
}

// func PrintPrice(settingsClub settings_club.SettingsCLub, db *data_base.DB) {
//	c := client.NewClient(time.Now(), 0, "", 0, db)
//}

func HandleAction(c *client.Client, t *table.Table, settingsClub settings_club.SettingsCLub) error {
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
		_, busy := c.GetClientFromDB(c.TableNumber)
		if busy {
			fmt.Printf("%s %v %s %d\n",
				c.CurrentTime.Format("15:04"),
				c.ActionId, c.ClientName, c.TableNumber)
			ClientsQueue.Enqueue(*c)
			c.SetInWaitingFromDB(c.ClientName, c.TableNumber)
			return fmt.Errorf("%s %v %s\n",
				c.CurrentTime.Format("15:04"),
				custom_errors.ErrCode,
				custom_errors.ErrPlaceIsBusy)
		}

		TablesWasBusy[t.Id] = *t

		t.StartTime = c.CurrentTime
		c.SetClientInDB(c.TableNumber, *c)
		t.SetTable(c.ClientName, *t)
		settingsClub.CountTables--

		fmt.Printf("%s %v %s %d\n",
			c.CurrentTime.Format("15:04"),
			c.ActionId, c.ClientName, c.TableNumber)
		return nil
	case 3:
		fmt.Printf("%s %v %s\n",
			c.CurrentTime.Format("15:04"),
			c.ActionId, c.ClientName)
		if ClientsQueue.Len() > settingsClub.CountTables {
			fmt.Printf("%s %v %s\n",
				c.CurrentTime.Format("15:04"),
				11, c.ClientName)
			return nil
		}
		tableNumber, _ := c.GetTableInWaitingFromDB(c.ClientName)
		_, busy := c.GetClientFromDB(tableNumber)
		if !busy {
			fmt.Printf("%s %v %s\n", c.CurrentTime.Format("15:04"),
				custom_errors.ErrCode,
				custom_errors.ErrICanWaitNoLonger)
		}

		return nil
	case 4:

		tableFromDB, _ := t.GetTable(c.ClientName)
		c.TableNumber = tableFromDB.Id
		tableFromDB.EndTime = c.CurrentTime

		fmt.Printf("%s %v %s\n",
			c.CurrentTime.Format("15:04"),
			c.ActionId, c.ClientName)

		if !ClientsQueue.IsEmpty() {
			cl := ClientsQueue.Dequeue()
			cl.TableNumber, _ = cl.GetTableInWaitingFromDB(cl.ClientName)
			cl.SetClientInDB(c.TableNumber, cl)
			t.SetTable(cl.ClientName, tableFromDB)
			cl.DeleteInWaitingFromDB(cl.ClientName)
			fmt.Printf("%s %v %s %d\n",
				c.CurrentTime.Format("15:04"),
				12, cl.ClientName, c.TableNumber)
		}

		tableFromDB.EndTime = c.CurrentTime
		tableFromDB.Price += settingsClub.Price * int(math.Ceil(tableFromDB.EndTime.Sub(tableFromDB.StartTime).Hours()))
		TablesWasBusy[tableFromDB.Id] = tableFromDB
		// fmt.Println("@dsadsa", tableFromDB.EndTime.Sub(c.CurrentTime).Hours())
		t.DeleteTable(c.ClientName)
		c.DeleteClientInDB(c.TableNumber)

		return nil
	default:
		return fmt.Errorf("%w",
			custom_errors.ErrActionNotExist)
	}
	return nil
}

var (
	ClientsQueue  = queue.Queue{}
	TablesWasBusy = map[int]table.Table{}
)

func main() {
	db := data_base.NewDB()
	ParseFile("test_file.txt", db)
	slog.Info("ParseFile - OK")
}
