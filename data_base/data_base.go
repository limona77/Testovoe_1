package data_base

import (
	"Testovoe/client"
	"Testovoe/table"
)

type DB struct {
	clientsInWaiting map[string]int
	clients          map[int]client.Client
	tables           map[string]table.Table
}
type IDB interface {
	GetTableInWaiting(key string) (int, bool)
	SetTableInWaiting(key string, val int)
	DeleteTableInWaiting(key string)
	GetClient(key int) (client.Client, bool)
	SetClient(key int, val client.Client)
	DeleteClient(key int)
	GetTableFromDB(key string) (table.Table, bool)
	SetTableFromDB(key string, val table.Table)
	DeleteTableFromDB(key string)
	ForEachInClientsNameFromDB() []string
	ForEachInTablesNameFromDB() []table.Table
}

func (db *DB) ForEachInClientsNameFromDB() []string {
	res := make([]string, 0)
	for _, v := range db.clients {
		res = append(res, v.ClientName)
	}
	return res
}

func (db *DB) GetTableInWaiting(key string) (int, bool) {
	if t, ok := db.clientsInWaiting[key]; ok {
		return t, true
	}
	return 0, false
}

func (db *DB) SetTableInWaiting(key string, val int) {
	db.clientsInWaiting[key] = val
}

func (db *DB) DeleteTableInWaiting(key string) {
	delete(db.clientsInWaiting, key)
}

func NewDB() *DB {
	return &DB{
		clientsInWaiting: map[string]int{},
		clients:          map[int]client.Client{},
		tables:           map[string]table.Table{},
	}
}

func (db *DB) GetClient(key int) (client.Client, bool) {
	if c, ok := db.clients[key]; ok {
		return c, true
	}
	return client.Client{}, false
}

func (db *DB) SetClient(key int, val client.Client) {
	db.clients[key] = val
}

func (db *DB) DeleteClient(key int) {
	delete(db.clients, key)
}

func (db *DB) GetTableFromDB(key string) (table.Table, bool) {
	if t, ok := db.tables[key]; ok {
		return t, true
	}
	return table.Table{}, false
}

func (db *DB) SetTableFromDB(key string, val table.Table) {
	db.tables[key] = val
}

func (db *DB) DeleteTableFromDB(key string) {
	delete(db.tables, key)
}

func (db *DB) ForEachInTablesNameFromDB() []int {
	res := make([]int, 0)
	for _, v := range db.tables {
		res = append(res, v.Id)
	}
	return res
}
