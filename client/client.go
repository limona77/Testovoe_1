package client

import (
	"time"
)

type DBInterfaceClient interface {
	GetTableInWaiting(key string) (int, bool)
	SetTableInWaiting(key string, val int)
	DeleteTableInWaiting(key string)
	GetClient(key int) (Client, bool)
	SetClient(key int, val Client)
	DeleteClient(key int)
	ForEachInClientsNameFromDB() []string
}
type IClient interface {
	CheckValidTime(startTime, endTime time.Time) bool
	GetTableInWaitingFromDB(key string) (int, bool)
	SetInWaitingFromDB(key string, val int)
	DeleteInWaitingFromDB(key string)
	GetClientFromDB(key int) (Client, bool)
	SetClientInDB(key int, val Client)
	DeleteClientInDB(key int)
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

func (c *Client) GetTableInWaitingFromDB(key string) (int, bool) {
	return c.db.GetTableInWaiting(key)
}

func (c *Client) SetInWaitingFromDB(key string, val int) {
	c.db.SetTableInWaiting(key, val)
}

func (c *Client) DeleteInWaitingFromDB(key string) {
	c.db.DeleteTableInWaiting(key)
}

func (c *Client) GetClientFromDB(key int) (Client, bool) {
	return c.db.GetClient(key)
}

func (c *Client) SetClientInDB(key int, val Client) {
	c.db.SetClient(key, val)
}

func (c *Client) DeleteClientInDB(key int) {
	c.db.DeleteClient(key)
}
