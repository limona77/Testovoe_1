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
	ClientsInClub   = map[string]struct{}{}
)

func HandleFirstAction(c *client.Client, settingsClub settings_club.SettingsCLub) error {
	fmt.Printf("%s %v %s\n", c.CurrentTime.Format("15:04"), c.ActionId, c.ClientName)
	_, exist := ClientsInClub[c.ClientName]
	if exist {
		return fmt.Errorf("%s %v %s\n",
			c.CurrentTime.Format("15:04"),
			ErrCode,
			custom_errors.ErrYouShallNotPass)
	}

	ok := c.CheckValidTime(settingsClub.StartTime, settingsClub.EndTime)
	if !ok {
		return fmt.Errorf(
			"%s %v %s\n",
			c.CurrentTime.Format("15:04"),
			ErrCode,
			custom_errors.ErrNotOpenYet)
	}
	ClientsInClub[c.ClientName] = struct{}{}
	return nil
}

func HandleSecondAction(c *client.Client, t *table.Table, settingsClub settings_club.SettingsCLub) error {
	_, busy := c.GetClient(c.TableNumber)
	if busy {
		fmt.Printf("%s %v %s %d\n",
			c.CurrentTime.Format("15:04"),
			c.ActionId, c.ClientName, c.TableNumber)
		ClientsQueue.Enqueue(*c)
		c.SetTableInWaiting(c.ClientName, c.TableNumber)
		return fmt.Errorf("%s %v %s\n",
			c.CurrentTime.Format("15:04"),
			ErrCode,
			custom_errors.ErrPlaceIsBusy)
	}
	TablesWasBusy[t.Id] = *t
	t.StartTime = c.CurrentTime
	c.SetClient(c.TableNumber, *c)
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
			CodeClientLeave, c.ClientName)
		return nil
	}

	tableNumber, _ := c.GetTableInWaiting(c.ClientName)
	_, busy := c.GetClient(tableNumber)
	if !busy {
		fmt.Printf("%s %v %s\n", c.CurrentTime.Format("15:04"),
			13,
			custom_errors.ErrICanWaitNoLonger)
	}
	return nil
}

func HandleFourthAction(c *client.Client, t *table.Table, settingsClub settings_club.SettingsCLub) error {
	tableFromDB, exist := t.GetTable(c.ClientName)
	if !exist {
		return fmt.Errorf("%s %v %s\n",
			c.CurrentTime.Format("15:04"),
			ErrCode,
			custom_errors.ErrClientUnknown)
	}
	c.TableNumber = tableFromDB.Id
	tableFromDB.EndTime = c.CurrentTime
	tableFromDB.Duration += tableFromDB.EndTime.Sub(tableFromDB.StartTime)
	fmt.Printf("%s %v %s\n",
		c.CurrentTime.Format("15:04"),
		c.ActionId, c.ClientName)
	if !ClientsQueue.IsEmpty() {
		cl := ClientsQueue.Dequeue()
		cl.TableNumber, _ = cl.GetTableInWaiting(cl.ClientName)
		cl.TableNumber = c.TableNumber

		tableFromDB.Price += settingsClub.Price * int(math.Ceil(tableFromDB.EndTime.Sub(tableFromDB.StartTime).Hours()))
		tableFromDB.StartTime = tableFromDB.EndTime
		cl.SetClient(cl.TableNumber, cl)

		t.SetTable(cl.ClientName, tableFromDB)
		cl.DeleteTableInWaiting(cl.ClientName)
		fmt.Printf("%s %v %s %d\n",
			c.CurrentTime.Format("15:04"),
			CodeClientSeat, cl.ClientName, cl.TableNumber)
	} else {
		tableFromDB.Price += settingsClub.Price*int(math.Ceil(tableFromDB.EndTime.Sub(tableFromDB.StartTime).Hours())) - 10
		tableFromDB.StartTime = tableFromDB.EndTime
	}

	TablesWasBusy[tableFromDB.Id] = tableFromDB
	t.DeleteTable(c.ClientName)
	c.DeleteClient(c.TableNumber)
	delete(ClientsInClub, c.ClientName)
	return nil
}
