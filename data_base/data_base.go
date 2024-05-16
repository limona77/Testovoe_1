package data_base

type DB struct {
	clientsInWaiting map[string]int
	clients          map[int]string
	tables           map[string]int
}
type IDB interface {
	GetTableInWaiting(key string) (int, bool)
	SetTableInWaiting(key string, val int)
	DeleteTableInWaiting(key string)
	GetClient(key int) (string, bool)
	SetClient(key int, val string)
	DeleteClient(key int)
	GetTable(key string) (int, bool)
	SetTable(key string, val int)
	DeleteTable(key string)
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
		clients:          map[int]string{},
		tables:           map[string]int{},
	}
}

func (db *DB) GetClient(key int) (string, bool) {
	if c, ok := db.clients[key]; ok {
		return c, true
	}
	return "", false
}

func (db *DB) SetClient(key int, val string) {
	db.clients[key] = val
}

func (db *DB) DeleteClient(key int) {
	delete(db.clients, key)
}

func (db *DB) GetTable(key string) (int, bool) {
	if t, ok := db.tables[key]; ok {
		return t, true
	}
	return 0, false
}

func (db *DB) SetTable(key string, val int) {
	db.tables[key] = val
}

func (db *DB) DeleteTable(key string) {
	delete(db.tables, key)
}
