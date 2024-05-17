package handlers

import (
	"Testovoe/client"
	"Testovoe/custom_errors"
	"Testovoe/queue"
	"Testovoe/settings_club"
	"Testovoe/table"
	"fmt"
	"math"
)

var (
	ClientsQueue    = queue.Queue{}
	CodeClientLeave = 11
	CodeClientSeat  = 12
	ErrCode         = 13
	TablesWasBusy   = map[int]table.Table{}
)

func HandleFirstAction(c *client.Client, settingsClub settings_club.SettingsCLub) error {
	fmt.Printf("%s %v %s\n", c.CurrentTime.Format("15:04"), c.ActionId, c.ClientName)
	ok := c.CheckValidTime(settingsClub.StartTime, settingsClub.EndTime)
	if !ok {
		return fmt.Errorf(
			"%s %v %s\n",
			c.CurrentTime.Format("15:04"),
			ErrCode,
			custom_errors.ErrNotOpenYet)
	}
	return nil
}

func HandleSecondAction(c *client.Client, t *table.Table, settingsClub settings_club.SettingsCLub) error {
	_, busy := c.GetClientFromDB(c.TableNumber)
	if busy {
		fmt.Printf("%s %v %s %d\n",
			c.CurrentTime.Format("15:04"),
			c.ActionId, c.ClientName, c.TableNumber)
		ClientsQueue.Enqueue(*c)
		c.SetTableInWaitingFromDB(c.ClientName, c.TableNumber)
		return fmt.Errorf("%s %v %s\n",
			c.CurrentTime.Format("15:04"),
			ErrCode,
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
}

func HandleThirdAction(c *client.Client, settingsClub settings_club.SettingsCLub) error {
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
			13,
			custom_errors.ErrICanWaitNoLonger)
	}
	return nil
}

func HandleFourthAction(c *client.Client, t *table.Table, settingsClub settings_club.SettingsCLub) error {
	tableFromDB, _ := t.GetTable(c.ClientName)
	c.TableNumber = tableFromDB.Id
	tableFromDB.EndTime = c.CurrentTime
	tableFromDB.Duration += tableFromDB.EndTime.Sub(tableFromDB.StartTime)
	fmt.Printf("%s %v %s\n",
		c.CurrentTime.Format("15:04"),
		c.ActionId, c.ClientName)
	if !ClientsQueue.IsEmpty() {
		cl := ClientsQueue.Dequeue()
		cl.TableNumber, _ = cl.GetTableInWaitingFromDB(cl.ClientName)
		cl.TableNumber = c.TableNumber

		tableFromDB.Price += settingsClub.Price * int(math.Ceil(tableFromDB.EndTime.Sub(tableFromDB.StartTime).Hours()))
		tableFromDB.StartTime = tableFromDB.EndTime
		cl.SetClientInDB(cl.TableNumber, cl)

		t.SetTable(cl.ClientName, tableFromDB)
		cl.DeleteInWaitingFromDB(cl.ClientName)
		fmt.Printf("%s %v %s %d\n",
			c.CurrentTime.Format("15:04"),
			12, cl.ClientName, cl.TableNumber)
	} else {
		tableFromDB.Price += settingsClub.Price*int(math.Ceil(tableFromDB.EndTime.Sub(tableFromDB.StartTime).Hours())) - 10
		tableFromDB.StartTime = tableFromDB.EndTime
	}

	TablesWasBusy[tableFromDB.Id] = tableFromDB
	t.DeleteTable(c.ClientName)
	c.DeleteClientInDB(c.TableNumber)
	return nil
}
