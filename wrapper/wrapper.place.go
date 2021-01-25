package wrapper

// Place menampung data tempat
type Place struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

// PlaceNewProvince provinsi baru
type PlaceNewProvince struct {
	Province string `json:"province" binding:"required,min=4,max=50"`
}
