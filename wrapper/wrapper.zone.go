package wrapper

// Zone zona tempat
type Zone struct {
	ID   int
	Name string
	List ZoneList
}

// ZoneList list dalam zona
type ZoneList struct {
	ID   int `db:"id"`
	Name int `db:"name"`
}
