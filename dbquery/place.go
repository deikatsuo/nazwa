package dbquery

import "github.com/jmoiron/sqlx"

// Place menampung data tempat
type Place struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

// GetProvinces mengambil data provinsi dari database
func GetProvinces(db *sqlx.DB) ([]Place, error) {
	var p []Place
	query := `SELECT id, INITCAP(name) AS name FROM "province"`
	err := db.Select(&p, query)
	if err != nil {
		return []Place{}, err
	}
	return p, nil
}

// GetCities mengambil data kota/kabupaten dari database
func GetCities(db *sqlx.DB, pid int) ([]Place, error) {
	var p []Place
	query := `SELECT id, INITCAP(name) AS name FROM "city" WHERE parent=$1`
	err := db.Select(&p, query, pid)
	if err != nil {
		return []Place{}, err
	}
	return p, nil
}

// GetDistricts mengambil data distrik/kecamatan dari database
func GetDistricts(db *sqlx.DB, pid int) ([]Place, error) {
	var p []Place
	query := `SELECT id, INITCAP(name) AS name FROM "district" WHERE parent=$1`
	err := db.Select(&p, query, pid)
	if err != nil {
		return []Place{}, err
	}
	return p, nil
}

// GetVillages mengambil data kelurahan/desa dari database
func GetVillages(db *sqlx.DB, pid int) ([]Place, error) {
	var p []Place
	query := `SELECT id, INITCAP(name) AS name FROM "village" WHERE parent=$1`
	err := db.Select(&p, query, pid)
	if err != nil {
		return []Place{}, err
	}
	return p, nil
}
