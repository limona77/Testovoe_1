package client

import (
	"Testovoe/data_base"
	"time"
)

type IClient interface {
	CheckValidTime(startTime, endTime time.Time) bool
	GetTableInWaitingFromDB(key string) (int, bool)
	SetInWaitingFromDB(key string, val int)
	DeleteInWaitingFromDB(key string)
	GetClientFromDB(key int) (string, bool)
	SetClientInDB(key int, val string)
	DeleteClientInDB(key int)
	GetTableFromDB(key string) (int, bool)
	SetTableInDB(key string, val int)
	DeleteTableInDB(key string)
}

type Client struct {
	CurrentTime time.Time
	ActionId    int
	ClientName  string
	TableNumber int
	db          *data_base.DB
}

func NewClient(currentTime time.Time, actionId int, clientName string, tableNumber int, db *data_base.DB) *Client {
	return &Client{
		CurrentTime: currentTime,
		ActionId:    actionId,
		ClientName:  clientName,
		TableNumber: tableNumber,
		db:          db,
	}
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

func (c *Client) GetClientFromDB(key int) (string, bool) {
	return c.db.GetClient(key)
}

func (c *Client) SetClientInDB(key int, val string) {
	c.db.SetClient(key, val)
}

func (c *Client) DeleteClientInDB(key int) {
	c.db.DeleteClient(key)
}

func (c *Client) GetTableFromDB(key string) (int, bool) {
	return c.db.GetTable(key)
}

func (c *Client) SetTableInDB(key string, val int) {
	c.db.SetTable(key, val)
}

func (c *Client) DeleteTableInDB(key string) {
	c.db.DeleteTable(key)
}
