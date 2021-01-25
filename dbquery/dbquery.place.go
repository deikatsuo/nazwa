package dbquery

import "nazwa/wrapper"

// PlaceGetProvinces mengambil data provinsi dari database
func PlaceGetProvinces() ([]wrapper.Place, error) {
	db := DB
	var p []wrapper.Place
	query := `SELECT id, INITCAP(name) AS name FROM "province" ORDER BY name`
	err := db.Select(&p, query)
	if err != nil {
		return []wrapper.Place{}, err
	}
	return p, nil
}

// PlaceGetCities mengambil data kota/kabupaten dari database
func PlaceGetCities(pid int) ([]wrapper.Place, error) {
	db := DB
	var p []wrapper.Place
	query := `SELECT id, INITCAP(name) AS name FROM "city" WHERE parent=$1 ORDER BY name`
	err := db.Select(&p, query, pid)
	if err != nil {
		return []wrapper.Place{}, err
	}
	return p, nil
}

// PlaceGetDistricts mengambil data distrik/kecamatan dari database
func PlaceGetDistricts(pid int) ([]wrapper.Place, error) {
	db := DB
	var p []wrapper.Place
	query := `SELECT id, INITCAP(name) AS name FROM "district" WHERE parent=$1 ORDER BY name`
	err := db.Select(&p, query, pid)
	if err != nil {
		return []wrapper.Place{}, err
	}
	return p, nil
}

// PlaceGetVillages mengambil data kelurahan/desa dari database
func PlaceGetVillages(pid int) ([]wrapper.Place, error) {
	db := DB
	var p []wrapper.Place
	query := `SELECT id, INITCAP(name) AS name FROM "village" WHERE parent=$1 ORDER BY name`
	err := db.Select(&p, query, pid)
	if err != nil {
		return []wrapper.Place{}, err
	}
	return p, nil
}

// PlaceNewProvince tambah provinsi baru
func PlaceNewProvince(name string, uid int) error {
	db := DB
	if _, err := db.Exec(`INSERT INTO "province" (parent, name, original, created_by) VALUES ($1, $2, $3, $4)`, 62, name, false, uid); err != nil {
		log.Warn("dbquery.zone.go ZoneNew() Gagal menambahkan provinsi")
		log.Error(err)
		return err
	}
	return nil
}
