package wrapper

// Zone wilayah dalam zona
type Zone struct {
	ID     int
	CityID int
	Name   string
	List   []ZoneList
}

// ZoneSelect mengambil data zona wilayah
type ZoneSelect struct {
	ID     int    `db:"id"`
	CityID int    `db:"city_id"`
	Name   string `db:"name"`
}

// ZoneList list dalam zona
type ZoneList struct {
	ID   int `db:"id"`
	Name int `db:"name"`
}
