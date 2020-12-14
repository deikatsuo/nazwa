package dbquery

import "github.com/jmoiron/sqlx"

// Place menampung data tempat
type Place struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

// PlaceGetProvinces mengambil data provinsi dari database
func PlaceGetProvinces(db *sqlx.DB) ([]Place, error) {
	var p []Place
	query := `SELECT id, INITCAP(name) AS name FROM "province"`
	err := db.Select(&p, query)
	if err != nil {
		return []Place{}, err
	}
	return p, nil
}

// PlaceGetCities mengambil data kota/kabupaten dari database
func PlaceGetCities(db *sqlx.DB, pid int) ([]Place, error) {
	var p []Place
	query := `SELECT id, INITCAP(name) AS name FROM "city" WHERE parent=$1`
	err := db.Select(&p, query, pid)
	if err != nil {
		return []Place{}, err
	}
	return p, nil
}

// PlaceGetDistricts mengambil data distrik/kecamatan dari database
func PlaceGetDistricts(db *sqlx.DB, pid int) ([]Place, error) {
	var p []Place
	query := `SELECT id, INITCAP(name) AS name FROM "district" WHERE parent=$1`
	err := db.Select(&p, query, pid)
	if err != nil {
		return []Place{}, err
	}
	return p, nil
}

// PlaceGetVillages mengambil data kelurahan/desa dari database
func PlaceGetVillages(db *sqlx.DB, pid int) ([]Place, error) {
	var p []Place
	query := `SELECT id, INITCAP(name) AS name FROM "village" WHERE parent=$1`
	err := db.Select(&p, query, pid)
	if err != nil {
		return []Place{}, err
	}
	return p, nil
}
