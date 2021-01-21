package wrapper

// Zone wilayah dalam zona
type Zone struct {
	ID     int
	CityID int
	Name   string
	List   []ZoneListSelect
}

// ZoneSelect mengambil data zona wilayah
type ZoneSelect struct {
	ID     int    `db:"id"`
	CityID int    `db:"city_id"`
	Name   string `db:"name"`
}

// ZoneListSelect list dalam zona
type ZoneListSelect struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}
