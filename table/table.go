package table

import (
	"time"
)

type DBInterfaceTable interface {
	GetTableFromDB(key string) (Table, bool)
	SetTableFromDB(key string, val Table)
	DeleteTableFromDB(key string)
	ForEachInTablesNameFromDB() []int
}
type ITable interface {
	GetTable(key string) (Table, bool)
	SetTable(key string, val Table)
	DeleteTable(key string)
	ForEachTables() []int
}
type Table struct {
	Id        int
	StartTime time.Time
	EndTime   time.Time
	Price     int
	Duration  time.Duration
	db        DBInterfaceTable
}

func NewTable(id int, price int, db DBInterfaceTable) *Table {
	return &Table{
		Id:    id,
		Price: price,
		db:    db,
	}
}

func (t *Table) GetTable(key string) (Table, bool) {
	return t.db.GetTableFromDB(key)
}

func (t *Table) SetTable(key string, val Table) {
	t.db.SetTableFromDB(key, val)
}

func (t *Table) DeleteTable(key string) {
	t.db.DeleteTableFromDB(key)
}

func (t *Table) ForEachTables() []int {
	return t.db.ForEachInTablesNameFromDB()
}
