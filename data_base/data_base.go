package data_base

type DB struct {
	hashMap map[int]int
}
type IDB interface {
	Get(key int) bool
	Set(key int, val int)
	Delete(key int)
}

func (db *DB) Get(key int) bool {
	if _, ok := db.hashMap[key]; ok {
		return true
	}
	return false
}

func (db *DB) Set(key int, val int) {
	db.hashMap[key] = val
}

func (db *DB) Delete(key int) {
	delete(db.hashMap, key)
}

func NewDB() *DB {
	return &DB{hashMap: map[int]int{}}
}
