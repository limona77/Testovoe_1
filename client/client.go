package client

import "C"
import (
	"Testovoe/table"
	"time"
)

type DBInterfaceClient interface {
	DeleteClientFromDB(key int)
	SetClientFromDB(key int, val Client)
	GetClientFromDB(key int) (Client, bool)
	DeleteTableInWaitingFromDB(key string)
	ForEachInClientsNameFromDB() []string
	SetTableFromDB(key string, val table.Table)
	GetTableInWaiting(key string) (int, bool)
	SetTableInWaiting(key string, val int)
}
type IClient interface {
	CheckValidTime(startTime, endTime time.Time) bool
	GetTableInWaiting(key string) (int, bool)
	SetTableInWaiting(key string, val int)
	DeleteTableInWaiting(key string)
	GetClient(key int) (Client, bool)
	SetClient(key int, val Client)
	DeleteClient(key int)

	ForEachInClientsName() []string
}

type Client struct {
	CurrentTime time.Time
	ActionId    int
	ClientName  string
	TableNumber int
	db          DBInterfaceClient
}

func NewClient(currentTime time.Time, actionId int, clientName string, tableNumber int, db DBInterfaceClient) *Client {
	return &Client{
		CurrentTime: currentTime,
		ActionId:    actionId,
		ClientName:  clientName,
		TableNumber: tableNumber,
		db:          db,
	}
}

func (c *Client) ForEachInClientsName() []string {
	return c.db.ForEachInClientsNameFromDB()
}

func (c *Client) CheckValidTime(startTime, endTime time.Time) bool {
	return c.CurrentTime.Before(endTime) && c.CurrentTime.After(startTime)
}

func (c *Client) GetTableInWaiting(key string) (int, bool) {
	return c.db.GetTableInWaiting(key)
}

func (c *Client) SetTableInWaiting(key string, val int) {
	c.db.SetTableInWaiting(key, val)
}

func (c *Client) DeleteTableInWaiting(key string) {
	c.db.DeleteTableInWaitingFromDB(key)
}

func (c *Client) GetClient(key int) (Client, bool) {
	return c.db.GetClientFromDB(key)
}

func (c *Client) SetClient(key int, val Client) {
	c.db.SetClientFromDB(key, val)
}

func (c *Client) DeleteClient(key int) {
	c.db.DeleteClientFromDB(key)
}
