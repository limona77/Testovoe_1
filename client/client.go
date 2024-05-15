package client

import (
	"Testovoe/data_base"
	"time"
)

type IClient interface {
	CheckValidTime(startTime, endTime time.Time) bool
	CheckInDB(key int) bool
	SetInDB(key int, val int)
	DeleteInDB(key int)

	//GetCurrentTime() time.Time
	//GetClientName() string
	//GetTableNumber() int
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

func (c *Client) CheckInDB(key int) bool {
	return c.db.Get(key)
}

func (c *Client) SetInDB(key int, val int) {
	c.db.Set(key, val)
}

func (c *Client) DeleteInDB(key int) {
	c.db.Delete(key)
}
