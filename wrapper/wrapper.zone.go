package wrapper

import "database/sql"

// Zone wilayah dalam zona
type Zone struct {
	ID        int
	Collector NameIDCode
	Name      string
	List      []ZoneListSelect
}

// ZoneSelect mengambil data zona wilayah
type ZoneSelect struct {
	ID          int           `db:"id"`
	CollectorID sql.NullInt32 `db:"collector_id"`
	Name        string        `db:"name"`
}

// ZoneListSelect list dalam zona
type ZoneListSelect struct {
	ID         int    `db:"id"`
	DistrictID int    `db:"district_id"`
	Name       string `db:"name"`
}

// ZoneAddListForm tambahkan list wilayah ke zona
type ZoneAddListForm struct {
	Lists []int `json:"lists" binding:"dive,numeric"`
}
