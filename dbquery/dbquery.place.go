package dbquery

// Place menampung data tempat
type Place struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

// PlaceGetProvinces mengambil data provinsi dari database
func PlaceGetProvinces() ([]Place, error) {
	db := DB
	var p []Place
	query := `SELECT id, INITCAP(name) AS name FROM "province"`
	err := db.Select(&p, query)
	if err != nil {
		return []Place{}, err
	}
	return p, nil
}

// PlaceGetCities mengambil data kota/kabupaten dari database
func PlaceGetCities(pid int) ([]Place, error) {
	db := DB
	var p []Place
	query := `SELECT id, INITCAP(name) AS name FROM "city" WHERE parent=$1`
	err := db.Select(&p, query, pid)
	if err != nil {
		return []Place{}, err
	}
	return p, nil
}

// PlaceGetDistricts mengambil data distrik/kecamatan dari database
func PlaceGetDistricts(pid int) ([]Place, error) {
	db := DB
	var p []Place
	query := `SELECT id, INITCAP(name) AS name FROM "district" WHERE parent=$1`
	err := db.Select(&p, query, pid)
	if err != nil {
		return []Place{}, err
	}
	return p, nil
}

// PlaceGetVillages mengambil data kelurahan/desa dari database
func PlaceGetVillages(pid int) ([]Place, error) {
	db := DB
	var p []Place
	query := `SELECT id, INITCAP(name) AS name FROM "village" WHERE parent=$1`
	err := db.Select(&p, query, pid)
	if err != nil {
		return []Place{}, err
	}
	return p, nil
}
