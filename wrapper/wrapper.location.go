package wrapper

import "database/sql"

// LocationZone wilayah dalam zona
type LocationZone struct {
	ID        int
	Collector NameIDCode
	Name      string
	List      []LocationZoneListsSelect
	CreatedBy NameIDCode
}

// LocationZoneSelect mengambil data zona wilayah
type LocationZoneSelect struct {
	ID          int           `db:"id"`
	CollectorID sql.NullInt32 `db:"collector_id"`
	Name        string        `db:"name"`
	CreatedBy   sql.NullInt32 `db:"created_by"`
}

// LocationZoneListsSelect list dalam zona
type LocationZoneListsSelect struct {
	ID         int    `db:"id"`
	ZoneLineID int    `db:"zone_line_id"`
	Name       string `db:"name"`
}

// LocationZoneAddListForm list wilayah pada zona
type LocationZoneAddListForm struct {
	Lists []int `json:"lists" binding:"dive,numeric"`
}

// LocationZoneNewForm zona baru
type LocationZoneNewForm struct {
	Zone string `json:"zone" binding:"required,min=4,max=25"`
}

// LocationLineNewForm line baru
type LocationLineNewForm struct {
	Code string `json:"code" binding:"required,min=1,max=5"`
	Name string `json:"name" binding:"required,min=5,max=25"`
}

// LocationLineSelect wrapper select line
type LocationLineSelect struct {
	ID   int    `db:"id"`
	Code string `db:"code"`
	Name string `db:"name"`
}

// LocationLine struk arah
type LocationLine struct {
	ID    int
	Code  string
	Name  string
	Count int
}
