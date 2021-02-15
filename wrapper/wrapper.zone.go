package wrapper

import "database/sql"

// Zone wilayah dalam zona
type Zone struct {
	ID        int
	Collector NameIDCode
	Name      string
	List      []ZoneListsSelect
	CreatedBy NameIDCode
}

// ZoneSelect mengambil data zona wilayah
type ZoneSelect struct {
	ID          int           `db:"id"`
	CollectorID sql.NullInt32 `db:"collector_id"`
	Name        string        `db:"name"`
	CreatedBy   sql.NullInt32 `db:"created_by"`
}

// ZoneListsSelect list dalam zona
type ZoneListsSelect struct {
	ID         int    `db:"id"`
	DistrictID int    `db:"district_id"`
	Name       string `db:"name"`
}

// ZoneAddListForm list wilayah pada zona
type ZoneAddListForm struct {
	Lists []int `json:"lists" binding:"dive,numeric"`
}

// ZoneNewForm zona baru
type ZoneNewForm struct {
	Zone string `json:"zone" binding:"required,min=4,max=25"`
}
